// Пакет для сохранения метрик в файл
package saver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/repositories"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/logger"
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
	WriteMetrics(repositories.ServerRepo) error
}

type ReadInterface interface {
	ReadMetrics() ([]repositories.Metric, error)
}

// SaverWriter --------------------------------------------------------------------------------------------------
type Writer struct {
	file     *os.File
	writer   *bufio.Writer
	filename string
}

func NewWriter(filename string) (*Writer, error) {
	// При создании файла удаляю предыдущее содержимое
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:     file,
		writer:   bufio.NewWriter(file),
		filename: filename,
	}, nil
}

func (storage *Writer) Close() error {
	if err := storage.writer.Flush(); err != nil {
		return err
	}
	return storage.file.Close()
}

// Сохраняю метрики из сервера в файл, причем предыдущее содержимое файла удаляю
func (storage *Writer) WriteMetrics(metrics repositories.ServerRepo) error {
	metricsSlice, err := metrics.GetAllMetricsSlice(context.Background())
	if err != nil {
		return err
	}
	if len(metricsSlice) == 0 {
		return nil
	}

	var metricsJSON bytes.Buffer
	enc := json.NewEncoder(&metricsJSON)
	if err := enc.Encode(metricsSlice); err != nil {
		return err
	}

	// Закрываем текущий writer и файл
	if err := storage.Close(); err != nil {
		return err
	}

	// Переоткрываем файл с флагами O_WRONLY и O_TRUNC для очистки
	file, err := os.OpenFile(storage.file.Name(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	storage.file = file
	storage.writer = bufio.NewWriter(file)

	n, err := storage.writer.Write(metricsJSON.Bytes())
	if err != nil {
		return err
	}
	if n != len(metricsJSON.Bytes()) {
		return fmt.Errorf("write metrics to file error: want write %d bytes, actual write %d bytes", len(metricsJSON.Bytes()), n)
	}
	if err := storage.writer.Flush(); err != nil {
		return err
	}

	logger.ServerLog.Info("write metrics to file")
	return nil
}

// Reader --------------------------------------------------------------------------------------------------
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

// Метод для чтения метрик из файла и записи их в слайс
func (saver *Reader) ReadMetrics() ([]repositories.Metric, error) {
	var bufRead bytes.Buffer

	_, err := bufRead.ReadFrom(saver.reader)
	if err != nil {
		return nil, err
	}

	bytesForRead := bufRead.Bytes()
	if len(bytesForRead) == 0 {
		return nil, nil
	}

	// преобразуем данные из JSON-представления в структуру
	var metrics = make([]repositories.Metric, 0)

	dec := json.NewDecoder(&bufRead)
	er := dec.Decode(&metrics)
	if er != nil {
		return nil, err
	}

	return metrics, nil
}

// Функция для загрузки метрик из файла в сервер
func AddMetricsFromFile(stor repositories.ServerRepo, reader ReadInterface) error {
	if GetRestore() {
		metrics, err := reader.ReadMetrics()
		if err != nil {
			return err
		}
		if err := stor.AddMetricsFromSlice(context.Background(), metrics); err != nil {
			return err
		}
	}
	return nil
}
