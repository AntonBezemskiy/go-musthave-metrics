package saver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/repositories"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/logger"
	"go.uber.org/zap"
)

// Global variable -------------------------------------------------
var (
	storeInterval   time.Duration
	fileStoragePath string
	restore         bool
)

func SetStoreInterval(interval time.Duration) {
	storeInterval = interval
}

func GetStoreInterval() time.Duration {
	return storeInterval
}

func SetFilestoragePath(path string) {
	fileStoragePath = path
}

func GetFilestoragePath() string {
	return fileStoragePath
}

func SetRestore(r bool) {
	restore = r
}

func GetRestore() bool {
	return restore
}

// end Global variable -------------------------------------------------

type WriterInterface interface {
	WriteMetrics(metric repositories.Metrics)
	FlushMetrics() error
}

type ReadInterface interface {
	ReadMetrics() ([]repositories.Metrics, error)
}

// SaverWriter --------------------------------------------------------------------------------------------------
type Writer struct {
	sync.Mutex
	file   *os.File
	writer *bufio.Writer
	buf    []repositories.Metrics
}

func NewWriter(filename string) (*Writer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:   file,
		writer: bufio.NewWriter(file),
		buf:    make([]repositories.Metrics, 0),
	}, nil
}

// Сохраняю метрики в буфер для последующей записи в файл
func (storage *Writer) WriteMetrics(metric repositories.Metrics) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	storage.buf = append(storage.buf, metric)

	logger.ServerLog.Info("write metrics to buffer for future flushing to file")
}

// Записываю накопленные метрики в файл и обнуляю буфер
func (storage *Writer) FlushMetrics() error {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()

	if len(storage.buf) == 0 {
		return nil
	}

	var metricsJSON bytes.Buffer
	enc := json.NewEncoder(&metricsJSON)
	if err := enc.Encode(storage.buf); err != nil {
		logger.ServerLog.Error("parse metric to json error", zap.String("error", error.Error(err)))
		return err
	}

	n, err := storage.file.Write(metricsJSON.Bytes())
	if err != nil {
		logger.ServerLog.Error("write metrics to file error", zap.String("error", error.Error(err)))
		return err
	}
	if n != len(metricsJSON.Bytes()) {
		logger.ServerLog.Error("write metrics to file error", zap.String("want write byte", strconv.Itoa(len(metricsJSON.Bytes()))),
			zap.String("actual write byte", strconv.Itoa(n)))
		return fmt.Errorf("write metrics to file error: want write %d bytes, actual write %d bytes", len(metricsJSON.Bytes()), n)
	}
	if err := storage.writer.Flush(); err != nil {
		logger.ServerLog.Error("flash buffer to the file error", zap.String("error", error.Error(err)))
		return err
	}

	// Обнуляю накопленные метрики, которые уже сохранены в файл
	storage.buf = make([]repositories.Metrics, 0)

	logger.ServerLog.Info("flush metrcis to file")
	return nil
}

// end SaverWriter --------------------------------------------------------------------------------------------------

// SaverReader --------------------------------------------------------------------------------------------------
type Reader struct {
	file   *os.File
	reader *bufio.Reader
}

func NewReader(filename string) (*Reader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file: file,
		// создаём новый Reader
		reader: bufio.NewReader(file),
	}, nil
}

func (saver *Reader) ReadMetrics() ([]repositories.Metrics, error) {
	var bufRead bytes.Buffer

	_, err := bufRead.ReadFrom(saver.reader)
	if err != nil {
		logger.ServerLog.Error("read metrics from file error", zap.String("error", error.Error(err)))
		return nil, err
	}

	bytesForRead := bufRead.Bytes()
	if len(bytesForRead) == 0 {
		return nil, nil
	}

	// преобразуем данные из JSON-представления в структуру
	var metrics = make([]repositories.Metrics, 0)

	dec := json.NewDecoder(&bufRead)
	er := dec.Decode(&metrics)
	if er != nil {
		logger.ServerLog.Error("decode metrics from file error", zap.String("error", error.Error(err)))
		return nil, err
	}

	return metrics, nil
}

func AddMetricsFromFile(stor repositories.ServerRepo, reader ReadInterface) {
	if GetRestore() {
		metrics, err := reader.ReadMetrics()
		if err != nil {
			log.Fatalf("read metrics from file error, file: %s. Error is: %s\n", GetFilestoragePath(), error.Error(err))
			//return
		}
		if err := stor.AddMetricsFromSlice(metrics); err != nil {
			log.Fatalf("add metrics from file: %s into server error. Error is: %s\n", GetFilestoragePath(), error.Error(err))
		}
	}
}