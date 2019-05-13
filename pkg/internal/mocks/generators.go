// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/solo-io/service-mesh-hub/pkg/kustomize/plugins (interfaces: NamedGenerator,NamedTransformer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	resmap "sigs.k8s.io/kustomize/pkg/resmap"
)

// MockNamedGenerator is a mock of NamedGenerator interface
type MockNamedGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockNamedGeneratorMockRecorder
}

// MockNamedGeneratorMockRecorder is the mock recorder for MockNamedGenerator
type MockNamedGeneratorMockRecorder struct {
	mock *MockNamedGenerator
}

// NewMockNamedGenerator creates a new mock instance
func NewMockNamedGenerator(ctrl *gomock.Controller) *MockNamedGenerator {
	mock := &MockNamedGenerator{ctrl: ctrl}
	mock.recorder = &MockNamedGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNamedGenerator) EXPECT() *MockNamedGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockNamedGenerator) Generate() (resmap.ResMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(resmap.ResMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (mr *MockNamedGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockNamedGenerator)(nil).Generate))
}

// Name mocks base method
func (m *MockNamedGenerator) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockNamedGeneratorMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockNamedGenerator)(nil).Name))
}

// MockNamedTransformer is a mock of NamedTransformer interface
type MockNamedTransformer struct {
	ctrl     *gomock.Controller
	recorder *MockNamedTransformerMockRecorder
}

// MockNamedTransformerMockRecorder is the mock recorder for MockNamedTransformer
type MockNamedTransformerMockRecorder struct {
	mock *MockNamedTransformer
}

// NewMockNamedTransformer creates a new mock instance
func NewMockNamedTransformer(ctrl *gomock.Controller) *MockNamedTransformer {
	mock := &MockNamedTransformer{ctrl: ctrl}
	mock.recorder = &MockNamedTransformerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNamedTransformer) EXPECT() *MockNamedTransformerMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockNamedTransformer) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockNamedTransformerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockNamedTransformer)(nil).Name))
}

// Transform mocks base method
func (m *MockNamedTransformer) Transform(arg0 resmap.ResMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transform", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transform indicates an expected call of Transform
func (mr *MockNamedTransformerMockRecorder) Transform(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transform", reflect.TypeOf((*MockNamedTransformer)(nil).Transform), arg0)
}
