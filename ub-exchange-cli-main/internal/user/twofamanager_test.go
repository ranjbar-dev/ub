// Package user_test tests the two-factor authentication (TOTP) manager. Covers:
//   - CheckCode: validating a TOTP code against a user's stored secret
//   - GenerateSecretCode: generating a new TOTP secret and provisioning URL
//
// Test data: dynamically generated TOTP keys via pquerna/otp/totp with
// no external mocks; uses real cryptographic secret generation.
package user_test

import (
	"database/sql"
	"exchange-go/internal/user"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTwoFaManager_CheckCode(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "somesite",
		AccountName: "testenv",
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	secret := key.Secret()
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	twoFaManager := user.NewTwoFaManager()
	u := user.User{
		Google2faSecretCode: sql.NullString{String: secret, Valid: true},
	}
	isValid := twoFaManager.CheckCode(u, code)
	assert.Equal(t, true, isValid)
}

func TestTwoFaManager_GenerateSecretCode(t *testing.T) {
	u := user.User{
		Email: "test@test.com",
	}
	twoFaManager := user.NewTwoFaManager()
	secret, url, err := twoFaManager.GenerateSecretCode(u)
	assert.Nil(t, err)
	assert.Equal(t, 16, len(secret))
	assert.NotEqual(t, "", url)
}
