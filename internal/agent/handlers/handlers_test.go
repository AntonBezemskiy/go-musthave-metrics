package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/compress"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/handlers"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/saver"
	"github.com/AntonBezemskiy/go-musthave-metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	stor := storage.NewDefaultMemStorage()
	saver, err := saver.NewWriter("./TestPush.json")
	require.NoError(t, err)

	type args struct {
		action      string
		typeMetric  string
		nameMetric  string
		valueMetric string
		client      *resty.Client
	}
	tests := []struct {
		name     string
		args     args
		wantStor storage.MemStorage
		wantErr  bool
	}{
		{
			name: "Count #1",
			args: args{
				action:      "update",
				typeMetric:  "counter",
				nameMetric:  "counter1",
				valueMetric: "4",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(nil, map[string]int64{"counter1": 4}),
			wantErr:  false,
		},
		{
			name: "Count error #1",
			args: args{
				action:      "update",
				typeMetric:  "wrangtype",
				nameMetric:  "counter1",
				valueMetric: "4",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(nil, map[string]int64{"counter1": 4}),
			wantErr:  true,
		},
		{
			name: "Gauge #1",
			args: args{
				action:      "update",
				typeMetric:  "gauge",
				nameMetric:  "gauge1",
				valueMetric: "3.14",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(map[string]float64{"gauge1": 3.14}, map[string]int64{"counter1": 4}),
			wantErr:  false,
		},
		{
			name: "Gauge error #1",
			args: args{
				action:      "update",
				typeMetric:  "wrangtype",
				nameMetric:  "gauge1",
				valueMetric: "3.14",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(map[string]float64{"gauge1": 3.14}, map[string]int64{"counter1": 4}),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/update/{metricType}/{metricName}/{metricValue}", func(res http.ResponseWriter, req *http.Request) {
				handlers.UpdateMetrics(res, req, stor, saver)
			})

			// Создаём тестовый сервер
			ts := httptest.NewServer(r)
			defer ts.Close()

			if err := Push(ts.URL, tt.args.action, tt.args.typeMetric, tt.args.nameMetric, tt.args.valueMetric, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("PushJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantStor, *stor)
		})
	}
	// Удаляю тестовый файл
	er := os.Remove("./TestPush.json")
	require.NoError(t, er)
}

func TestPushJSON(t *testing.T) {
	stor := storage.NewDefaultMemStorage()
	saver, err := saver.NewWriter("./TestPushJSON.json")
	require.NoError(t, err)

	type args struct {
		action      string
		typeMetric  string
		nameMetric  string
		valueMetric string
		client      *resty.Client
	}
	tests := []struct {
		name     string
		args     args
		wantStor storage.MemStorage
		wantErr  bool
	}{
		{
			name: "Count #1",
			args: args{
				action:      "update",
				typeMetric:  "counter",
				nameMetric:  "counter1",
				valueMetric: "4",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(nil, map[string]int64{"counter1": 4}),
			wantErr:  false,
		},
		{
			name: "Count error #1",
			args: args{
				action:      "update",
				typeMetric:  "wrangtype",
				nameMetric:  "counter1",
				valueMetric: "4",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(nil, map[string]int64{"counter1": 4}),
			wantErr:  true,
		},
		{
			name: "Gauge #1",
			args: args{
				action:      "update",
				typeMetric:  "gauge",
				nameMetric:  "gauge1",
				valueMetric: "3.14",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(map[string]float64{"gauge1": 3.14}, map[string]int64{"counter1": 4}),
			wantErr:  false,
		},
		{
			name: "Gauge error #1",
			args: args{
				action:      "update",
				typeMetric:  "wrangtype",
				nameMetric:  "gauge1",
				valueMetric: "3.14",
				client:      resty.New(),
			},
			wantStor: *storage.NewMemStorage(map[string]float64{"gauge1": 3.14}, map[string]int64{"counter1": 4}),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Post("/update", compress.GzipMiddleware(handlers.UpdateMetricsJSONHandler(stor, saver)))

			// Создаём тестовый сервер
			ts := httptest.NewServer(r)
			defer ts.Close()

			if err := PushJSON(ts.URL, tt.args.action, tt.args.typeMetric, tt.args.nameMetric, tt.args.valueMetric, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("PushJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantStor, *stor)
		})
	}
	// Удаляю тестовый файл
	er := os.Remove("./TestPushJSON.json")
	require.NoError(t, er)
}
