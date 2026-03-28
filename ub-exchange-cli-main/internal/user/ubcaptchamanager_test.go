// Package user_test tests the UB captcha manager. Covers:
//   - CheckUbCaptcha: validating a correctly signed captcha token (valid and invalid)
//   - Decrypt: decrypting a PEM-encoded RSA-encrypted message
//   - Encrypt: encrypting a message and round-tripping through Decrypt
//   - GetKey: retrieving the current RSA key pair
//   - NewKey: generating a fresh RSA key pair
//
// Test data: testify mock for Logger; hard-coded Base64 captcha tokens
// and PEM-encoded encrypted messages for deterministic decryption tests.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUbCaptchaManager_CheckUbCaptcha(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	//valid ub-captcha
	ubCaptchaString := "ub-captcha_Qaac7K4w8HXcAqa6fC2NR6VeXeKQzn9oUUMjqFZHM51k21KtaZ+AdwWnTmGeo89Y1rP/zK35TXe5jWjc/ZiHViMeITryVE62TPpx2GjW4SSakhzX6cO1ji6k/wCa87l73a19sSFFZhkQbcx0VqQOH47mncMla/D/v80+7/JFX9DbZRN633HA1ksvtWfdtJVZQ/5NZE1i1zIinnaCTmb2wWR7DHju6It863wTpXhCbyGgs16BTUO6UpoEN8jH3roE6TW7d+kc2nCjxvtV6J5lvJDe+ET602MDvhic97OGw3rY1XZjFJnfCA82w+B/yDdD98pXumikyA5IpQgxgG5ZpQ=="

	result, err := ubCaptchaManager.CheckUbCaptcha(ubCaptchaString)
	assert.Nil(t, err)
	assert.True(t, result)
	logger.AssertExpectations(t)
}

func TestUbCaptchaManager_CheckUbCaptcha2(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	//invalid ub-captcha
	ubCaptchaString := "ub-captcha_l/hi26loujWcrNEZrUtDG5a4CTb0KhK52WOJfYTJ9Bgsj/uwoy90w02YAArwS1SnccnaHQfU2OYsDwMCRLMUuin+koBhZzzXhx16K2kUBMmMrmt0eFp0vvhaAmFWpBW5WOAaug1dxmtwrW0hBPDxSc3ixeYuFJZffbIfb/TyXijzP4AzpjgOU6cNI44/OnOhYNIJXlggMRgoAdeSBJejhcL2O3oLY0ejJ2nWOpJ0+6Ip5L/t75UpeRnksbK6tHvGMt4eK3tpAVTTKnBiOzLT3gdYNej8yS3Q3ctmnC8e64n2Bv1nRUdfvJ6cs/IrVDEZGPqnptdsF5FXOPw6OWQV7w=="

	result, err := ubCaptchaManager.CheckUbCaptcha(ubCaptchaString)
	assert.Nil(t, err)
	assert.False(t, result)
	logger.AssertExpectations(t)
}

func TestUbCaptchaManager_Decrypt(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	encryptedMsg := `-----BEGIN MESSAGE-----
eOKnzwV63Wv88swnI3khNUCTlH6+HP5I0ivBoGnsATRotCdhSh8P2tKR5wj9WCHq
ZVu1mBVLlgfnh5Cze0HV6zG82XX1BEx3/Dt1cr3V/JKIKT2ltq+zrQSguMuqcviB
u9aSzZ6DT6HIlFHU+COMepulA48gc5S6y4GnrRSOV+sIfqCCyEZM69p2H6HeGAWk
50Lf70wKzAa7yaXcim3OaGWCzqI5Dih9rqw6xo0Z6iDfIJ4pFHJ0N9+v1mHPz0YV
tlwAruh3QOD9xdubJeLAHLt/o157KdUzmazn19sj8d1Tv7hq9P8KQ9Fsd2wu0Ie+
BX3usZFsocdGhrnmGZBT2g==
-----END MESSAGE-----`

	msg := "hello, world"

	decryptedMsg, err := ubCaptchaManager.Decrypt(encryptedMsg)
	assert.Nil(t, err)
	assert.Equal(t, msg, decryptedMsg)
	logger.AssertExpectations(t)
}

func TestUbCaptchaManager_Encrypt(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	msg := "hello, world"

	encryptedMsg, err := ubCaptchaManager.Encrypt(msg)
	assert.Nil(t, err)
	decryptedMsg, err := ubCaptchaManager.Decrypt(encryptedMsg)
	assert.Nil(t, err)
	assert.Equal(t, msg, decryptedMsg)
	logger.AssertExpectations(t)
}

func TestUbCaptchaManager_GetKey(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	key, err := ubCaptchaManager.GetKey()

	assert.Nil(t, err)
	assert.IsType(t, user.UbCaptchaKey{}, key)
	logger.AssertExpectations(t)
}

func TestUbCaptchaManager_NewKey(t *testing.T) {
	logger := new(mocks.Logger)

	ubCaptchaManager := user.NewUbCaptchaManager(logger)

	key, err := ubCaptchaManager.NewKey()

	assert.Nil(t, err)
	assert.IsType(t, user.UbCaptchaKey{}, key)
	logger.AssertExpectations(t)
}
