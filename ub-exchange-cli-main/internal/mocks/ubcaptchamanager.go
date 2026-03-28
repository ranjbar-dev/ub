package mocks

import (
	"exchange-go/internal/user"
	"github.com/stretchr/testify/mock"
)

type UbCaptchaManager struct {
	mock.Mock
}

func (_m *UbCaptchaManager) CheckUbCaptcha(ubCaptchaStr string) (bool, error) {
	args := _m.Called(ubCaptchaStr)
	return args.Bool(0), args.Error(1)
}

func (_m *UbCaptchaManager) Encrypt(plainText string) (string, error) {
	args := _m.Called(plainText)
	return args.String(0), args.Error(1)
}

func (_m *UbCaptchaManager) Decrypt(encryptedText string) (string, error) {
	args := _m.Called(encryptedText)
	return args.String(0), args.Error(1)
}

func (_m *UbCaptchaManager) NewKey() (user.UbCaptchaKey, error) {
	args := _m.Called()
	return args.Get(0).(user.UbCaptchaKey), args.Error(1)
}

func (_m *UbCaptchaManager) GetKey() (user.UbCaptchaKey, error) {
	args := _m.Called()
	return args.Get(0).(user.UbCaptchaKey), args.Error(1)
}

func (_m *UbCaptchaManager) SaveKeyToPemFile(key user.UbCaptchaKey) error {
	args := _m.Called(key)
	return args.Error(0)
}
