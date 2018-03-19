package mocks

import (
	"github.com/NeowayLabs/wabbit"
)

// MockDelivery : mock wabbit delivery
type MockDelivery struct {
	Data []byte
}

// Ack : ack a message
func (m *MockDelivery) Ack(multiple bool) error {
	return nil
}

// Nack : nack a message
func (m *MockDelivery) Nack(multiple, request bool) error {
	return nil
}

// Reject : reject a message
func (m *MockDelivery) Reject(requeue bool) error {
	return nil
}

// Body : message body
func (m *MockDelivery) Body() []byte {
	return m.Data
}

// Headers : message headers
func (m *MockDelivery) Headers() wabbit.Option {
	return wabbit.Option{}
}

// DeliveryTag : delivery tag
func (m *MockDelivery) DeliveryTag() uint64 {
	return 0
}

// ConsumerTag : consumer tag
func (m *MockDelivery) ConsumerTag() string {
	return ""
}

// MessageId : message id
func (m *MockDelivery) MessageId() string {
	return ""
}
