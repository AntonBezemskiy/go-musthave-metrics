// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AntonBezemskiy/go-musthave-metrics/internal/repositories (interfaces: ServerRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	repositories "github.com/AntonBezemskiy/go-musthave-metrics/internal/repositories"
	gomock "github.com/golang/mock/gomock"
)

// MockServerRepo is a mock of ServerRepo interface.
type MockServerRepo struct {
	ctrl     *gomock.Controller
	recorder *MockServerRepoMockRecorder
}

// MockServerRepoMockRecorder is the mock recorder for MockServerRepo.
type MockServerRepoMockRecorder struct {
	mock *MockServerRepo
}

// NewMockServerRepo creates a new mock instance.
func NewMockServerRepo(ctrl *gomock.Controller) *MockServerRepo {
	mock := &MockServerRepo{ctrl: ctrl}
	mock.recorder = &MockServerRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerRepo) EXPECT() *MockServerRepoMockRecorder {
	return m.recorder
}

// AddCounter mocks base method.
func (m *MockServerRepo) AddCounter(arg0 context.Context, arg1 string, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCounter", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCounter indicates an expected call of AddCounter.
func (mr *MockServerRepoMockRecorder) AddCounter(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCounter", reflect.TypeOf((*MockServerRepo)(nil).AddCounter), arg0, arg1, arg2)
}

// AddGauge mocks base method.
func (m *MockServerRepo) AddGauge(arg0 context.Context, arg1 string, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddGauge", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddGauge indicates an expected call of AddGauge.
func (mr *MockServerRepoMockRecorder) AddGauge(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGauge", reflect.TypeOf((*MockServerRepo)(nil).AddGauge), arg0, arg1, arg2)
}

// AddMetricsFromSlice mocks base method.
func (m *MockServerRepo) AddMetricsFromSlice(arg0 context.Context, arg1 []repositories.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetricsFromSlice", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMetricsFromSlice indicates an expected call of AddMetricsFromSlice.
func (mr *MockServerRepoMockRecorder) AddMetricsFromSlice(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetricsFromSlice", reflect.TypeOf((*MockServerRepo)(nil).AddMetricsFromSlice), arg0, arg1)
}

// Bootstrap mocks base method.
func (m *MockServerRepo) Bootstrap(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrap", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Bootstrap indicates an expected call of Bootstrap.
func (mr *MockServerRepoMockRecorder) Bootstrap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrap", reflect.TypeOf((*MockServerRepo)(nil).Bootstrap), arg0)
}

// GetAllMetrics mocks base method.
func (m *MockServerRepo) GetAllMetrics(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockServerRepoMockRecorder) GetAllMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockServerRepo)(nil).GetAllMetrics), arg0)
}

// GetAllMetricsSlice mocks base method.
func (m *MockServerRepo) GetAllMetricsSlice(arg0 context.Context) ([]repositories.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetricsSlice", arg0)
	ret0, _ := ret[0].([]repositories.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetricsSlice indicates an expected call of GetAllMetricsSlice.
func (mr *MockServerRepoMockRecorder) GetAllMetricsSlice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetricsSlice", reflect.TypeOf((*MockServerRepo)(nil).GetAllMetricsSlice), arg0)
}

// GetMetric mocks base method.
func (m *MockServerRepo) GetMetric(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockServerRepoMockRecorder) GetMetric(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockServerRepo)(nil).GetMetric), arg0, arg1, arg2)
}