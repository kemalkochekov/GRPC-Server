// Code generated by MockGen. DO NOT EDIT.
// Source: ./kafkaInterface.go
//
// Generated by this command:
//
//	mockgen -source ./kafkaInterface.go -destination=./mocks/kafkaInterface.go -package=mock_kafka
//
// Package mock_kafka is a generated GoMock package.
package mock_kafka

import (
	kafkaEntities "GRPC_Server/internal/kafka/kafkaEntities"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockProducerInterface is a mock of ProducerInterface interface.
type MockProducerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockProducerInterfaceMockRecorder
}

// MockProducerInterfaceMockRecorder is the mock recorder for MockProducerInterface.
type MockProducerInterfaceMockRecorder struct {
	mock *MockProducerInterface
}

// NewMockProducerInterface creates a new mock instance.
func NewMockProducerInterface(ctrl *gomock.Controller) *MockProducerInterface {
	mock := &MockProducerInterface{ctrl: ctrl}
	mock.recorder = &MockProducerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducerInterface) EXPECT() *MockProducerInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockProducerInterface) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockProducerInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockProducerInterface)(nil).Close))
}

// SendMessage mocks base method.
func (m *MockProducerInterface) SendMessage(topic string, kakfaRequestMessage kafkaEntities.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", topic, kakfaRequestMessage)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockProducerInterfaceMockRecorder) SendMessage(topic, kakfaRequestMessage any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockProducerInterface)(nil).SendMessage), topic, kakfaRequestMessage)
}

// Topic mocks base method.
func (m *MockProducerInterface) Topic() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Topic")
	ret0, _ := ret[0].(string)
	return ret0
}

// Topic indicates an expected call of Topic.
func (mr *MockProducerInterfaceMockRecorder) Topic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Topic", reflect.TypeOf((*MockProducerInterface)(nil).Topic))
}

// MockConsumerInterface is a mock of ConsumerInterface interface.
type MockConsumerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerInterfaceMockRecorder
}

// MockConsumerInterfaceMockRecorder is the mock recorder for MockConsumerInterface.
type MockConsumerInterfaceMockRecorder struct {
	mock *MockConsumerInterface
}

// NewMockConsumerInterface creates a new mock instance.
func NewMockConsumerInterface(ctrl *gomock.Controller) *MockConsumerInterface {
	mock := &MockConsumerInterface{ctrl: ctrl}
	mock.recorder = &MockConsumerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerInterface) EXPECT() *MockConsumerInterfaceMockRecorder {
	return m.recorder
}

// ConsumerClose mocks base method.
func (m *MockConsumerInterface) ConsumerClose() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumerClose")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumerClose indicates an expected call of ConsumerClose.
func (mr *MockConsumerInterfaceMockRecorder) ConsumerClose() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumerClose", reflect.TypeOf((*MockConsumerInterface)(nil).ConsumerClose))
}

// ReadMessages mocks base method.
func (m *MockConsumerInterface) ReadMessages(ctx context.Context, topic string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessages", ctx, topic)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadMessages indicates an expected call of ReadMessages.
func (mr *MockConsumerInterfaceMockRecorder) ReadMessages(ctx, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessages", reflect.TypeOf((*MockConsumerInterface)(nil).ReadMessages), ctx, topic)
}
