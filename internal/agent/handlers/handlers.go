package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agent/compress"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agent/logger"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/agent/storage"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/repositories"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	pollInterval   time.Duration = 2
	reportInterval time.Duration = 10
)

func SetPollInterval(interval time.Duration) {
	pollInterval = interval
}

func GetPollInterval() time.Duration {
	return pollInterval
}

func SetReportInterval(interval time.Duration) {
	reportInterval = interval
}

func GetReportInterval() time.Duration {
	return reportInterval
}

// CollectMetrics собирает метрики
func SyncCollectMetrics(metrics *storage.MetricsStats) {
	metrics.Lock()
	defer metrics.Unlock()
	metrics.CollectMetrics()
}

// CollectMetricsTimer запускает сбор метрик с интервалом
func CollectMetricsTimer(metrics *storage.MetricsStats) {
	sleepInterval := GetPollInterval() * time.Second
	for {
		SyncCollectMetrics(metrics)
		time.Sleep(sleepInterval)
	}
}

func BuildMetric(typeMetric, nameMetric, valueMetric string) (metric repositories.Metric, err error) {
	// Строю структуру метрики для сериализации из принятых параметров
	metric.ID = nameMetric
	metric.MType = typeMetric

	switch typeMetric {
	case "counter":
		val, errParse := strconv.ParseInt(valueMetric, 10, 64)
		if errParse != nil {
			logger.AgentLog.Error("Convert string to int64 error: ", zap.String("error: ", error.Error(err)))
			err = errParse
			return
		}
		metric.Delta = &val
	case "gauge":
		val, errParse := strconv.ParseFloat(valueMetric, 64)
		if errParse != nil {
			logger.AgentLog.Error("Convert string to float64 error: ", zap.String("error: ", error.Error(err)))
			err = errParse
			return
		}
		metric.Value = &val
	default:
		logger.AgentLog.Error("Invalid type of metric", zap.String("type: ", metric.MType)) //---------------------------------------------
		err = fmt.Errorf("get invalid type of metric: %s", typeMetric)
		return
	}
	logger.AgentLog.Debug(fmt.Sprintf("Success build metric structure for JSON: name: %s, type: %s, delta: %d, value: %d", metric.ID, metric.MType, metric.Delta, metric.Value))
	return
}

// Push отправляет метрику на сервер в JSON формате и возвращает ошибку при неудаче
func PushJSON(address, action, typeMetric, nameMetric, valueMetric string, client *resty.Client) error {
	metric, err := BuildMetric(typeMetric, nameMetric, valueMetric)
	if err != nil {
		logger.AgentLog.Error("Build metric error", zap.String("error", error.Error(err)))
		return err
	}

	// сериализую полученную струтктуру с метриками в json-представление  в виде слайса байт
	var bufEncode bytes.Buffer
	enc := json.NewEncoder(&bufEncode)
	if err := enc.Encode(metric); err != nil {
		logger.AgentLog.Error("Encode message error", zap.String("error", error.Error(err)))
		return err
	}

	// Сжатие данных для передачи
	compressBody, err := compress.Compress(bufEncode.Bytes())
	if err != nil {
		logger.AgentLog.Error("Fail to comperess push data ", zap.String("error", error.Error(err)))
		return err
	}

	url := fmt.Sprintf("%s/%s", address, action)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressBody).
		Post(url)

	if err != nil {
		logger.AgentLog.Error("Push json metric to server error ", zap.String("error", error.Error(err)))
		return err
	}

	logger.AgentLog.Debug("Get answer from server", zap.String("Content-Encoding", resp.Header().Get("Content-Encoding")),
		zap.String("statusCode", fmt.Sprintf("%d", resp.StatusCode())),
		zap.String("Content-Type", resp.Header().Get("Content-Type")),
		zap.String("Content-Encoding", fmt.Sprint(resp.Header().Values("Content-Encoding"))))

	if resp.StatusCode() != http.StatusOK {
		logger.AgentLog.Error("Geting status is not 200 ", zap.String("statusCode", fmt.Sprintf("%d", resp.StatusCode())))
		return fmt.Errorf("status code is: %d", resp.StatusCode())
	}

	contentEncoding := resp.Header().Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		logger.AgentLog.Debug("Get compress answer data in PushJSON function", zap.String("Content-Encoding", contentEncoding))
	} else {
		logger.AgentLog.Debug("Get uncompress answer data in PushJSON function", zap.String("Content-Encoding", contentEncoding))
	}

	responceMetric := resp.Body()
	if !bytes.Equal(bufEncode.Bytes(), responceMetric) {
		return fmt.Errorf("answer metric from server not equal pushing metric: get %d, want %d", responceMetric, bufEncode.Bytes())
	}

	// Десериализую данные полученные от сервера, в основном для дебага
	var resJSON repositories.Metric
	buRes := bytes.NewBuffer(responceMetric)
	dec := json.NewDecoder(buRes)
	if err := dec.Decode(&resJSON); err != nil {
		logger.AgentLog.Error("decode decompress data from server error ", zap.String("error", error.Error(err)))
		return err
	}
	logger.AgentLog.Debug(fmt.Sprintf("decode metric from server %s", resJSON.String()))

	logger.AgentLog.Debug(fmt.Sprintf("Success push metric in JSON format: typeMetric - %s, nameMetric - %s, valueMetric - %s", typeMetric, nameMetric, valueMetric))
	return nil
}

// Push отправляет метрику на сервер и возвращает ошибку при неудаче
func Push(address, action, typemetric, namemetric, valuemetric string, client *resty.Client) error {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", address, action, typemetric, namemetric, valuemetric)
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		return fmt.Errorf("error with post: %s, %w", url, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d for url: %s", resp.StatusCode(), url)
	}
	return nil
}

// PushMetrics отправляет все метрики
func PushMetrics(address, action string, metrics *storage.MetricsStats, client *resty.Client) {
	metrics.Lock()
	defer metrics.Unlock()

	for _, metricName := range storage.AllMetrics {
		typeMetric, value, err := metrics.GetMetricString(metricName)
		if err != nil {
			logger.AgentLog.Error(fmt.Sprintf("Failed to get metric %s: %v\n", typeMetric, err), zap.String("action", "push metrics"))
			continue
		}
		er := PushJSON(address, action, typeMetric, metricName, value, client)
		if er != nil {
			logger.AgentLog.Error(fmt.Sprintf("Failed to push metric %s: %v\n", typeMetric, er), zap.String("action", "push metrics"))
		}
	}
}

// Push отправляет метрику на сервер в JSON формате и возвращает ошибку при неудаче
func PushBatch(address, action string, metricsSlice []repositories.Metric, client *resty.Client) error {

	// сериализую полученную слайс с метриками в json-представление  в виде слайса байт
	var bufEncode bytes.Buffer
	enc := json.NewEncoder(&bufEncode)
	if err := enc.Encode(metricsSlice); err != nil {
		logger.AgentLog.Error("Encode message error", zap.String("error", error.Error(err)))
		return err
	}

	// Сжатие данных для передачи
	compressBody, err := compress.Compress(bufEncode.Bytes())
	if err != nil {
		logger.AgentLog.Error("Fail to comperess push data ", zap.String("error", error.Error(err)))
		return err
	}

	url := fmt.Sprintf("%s/%s", address, action)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(compressBody).
		Post(url)

	if err != nil {
		logger.AgentLog.Error("Push batch json metrics to server error ", zap.String("error", error.Error(err)))
		return err
	}

	logger.AgentLog.Debug("Get answer from server", zap.String("Content-Encoding", resp.Header().Get("Content-Encoding")),
		zap.String("statusCode", fmt.Sprintf("%d", resp.StatusCode())),
		zap.String("Content-Type", resp.Header().Get("Content-Type")),
		zap.String("Content-Encoding", fmt.Sprint(resp.Header().Values("Content-Encoding"))))

	if resp.StatusCode() != http.StatusOK {
		logger.AgentLog.Error("Geting status is not 200 ", zap.String("statusCode", fmt.Sprintf("%d", resp.StatusCode())))
		return fmt.Errorf("status code is: %d", resp.StatusCode())
	}

	contentEncoding := resp.Header().Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		logger.AgentLog.Debug("Get compress answer data in PushBatch function", zap.String("Content-Encoding", contentEncoding))
	} else {
		logger.AgentLog.Debug("Get uncompress answer data in PushBatch function", zap.String("Content-Encoding", contentEncoding))
	}

	responceMetrics := resp.Body()
	if !bytes.Equal(bufEncode.Bytes(), responceMetrics) {
		return fmt.Errorf("answer metric from server not equal pushing metric: get %d, want %d", responceMetrics, bufEncode.Bytes())
	}

	// Десериализую данные полученные от сервера, в основном для дебага
	var resJSON []repositories.Metric
	buRes := bytes.NewBuffer(responceMetrics)
	dec := json.NewDecoder(buRes)
	if err := dec.Decode(&resJSON); err != nil {
		logger.AgentLog.Error("decode decompress data from server error ", zap.String("error", error.Error(err)))
		return err
	}

	logger.AgentLog.Debug("Success push batch metrics in JSON format")
	return nil
}

// PushMetrics отправляет все метрики батчем
func PushMetricsBatch(address, action string, metrics *storage.MetricsStats, client *resty.Client) {
	metrics.Lock()
	defer metrics.Unlock()
	metricsSlice := make([]repositories.Metric, 0)

	// создаю слайс с метриками для отправки батчем
	for _, metricName := range storage.AllMetrics {
		typeMetric, value, err := metrics.GetMetricString(metricName)
		if err != nil {
			logger.AgentLog.Error(fmt.Sprintf("Failed to get metric %s: %v\n", typeMetric, err), zap.String("action", "push metrics"))
			continue
		}
		metric, err := BuildMetric(typeMetric, metricName, value)
		if err != nil {
			logger.AgentLog.Error(fmt.Sprintf("Failed to build metric structer %s: %v\n", typeMetric, err), zap.String("action", "push metrics"))
			continue
		}
		metricsSlice = append(metricsSlice, metric)
	}
	er := PushBatch(address, action, metricsSlice, client)
	if er != nil {
		logger.AgentLog.Error("Failed to push batch metrics", zap.String("action", "push metrics"), zap.String("error", error.Error(er)))
	}
}

// PushMetricsTimer запускает отправку метрик с интервалом
func PushMetricsTimer(address, action string, metrics *storage.MetricsStats) {
	sleepInterval := GetReportInterval() * time.Second
	for {
		client := resty.New()
		PushMetricsBatch(address, action, metrics, client)
		logger.AgentLog.Debug("Running agent", zap.String("action", "push metrics"))
		time.Sleep(sleepInterval)
	}
}
