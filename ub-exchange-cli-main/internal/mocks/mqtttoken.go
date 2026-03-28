package mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MqttToken struct {
	mock.Mock
}

func (m *MqttToken) Wait() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MqttToken) WaitTimeout(d time.Duration) bool {
	args := m.Called(d)
	return args.Bool(0)
}

func (m *MqttToken) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(<-chan struct{})

}

func (m *MqttToken) Error() error {
	args := m.Called()
	return args.Error(0)
}
