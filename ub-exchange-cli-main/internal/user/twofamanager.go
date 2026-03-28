package user

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const TwoFaIssuer = "unitedbit.com"

// TwoFaManager handles TOTP-based two-factor authentication using Google Authenticator.
type TwoFaManager interface {
	// CheckCode validates a 6-digit TOTP code against the user's stored secret.
	CheckCode(u User, code string) bool
	// GenerateSecretCode creates a new TOTP secret (or reuses an existing one) and
	// returns the secret string and a QR code URL for Google Authenticator setup.
	GenerateSecretCode(u User) (secret string, url string, err error)
}

type twoFaManager struct {
}

func (s *twoFaManager) CheckCode(u User, code string) bool {
	if !u.Google2faSecretCode.Valid {
		return false
	}
	isValid := totp.Validate(code, u.Google2faSecretCode.String)
	return isValid
}

func (s *twoFaManager) GenerateSecretCode(u User) (secret string, qrCodeURL string, err error) {
	var otpKey *otp.Key
	if u.Google2faSecretCode.Valid && u.Google2faSecretCode.String != "" {
		otpKey, err = s.getOtpKeyForUser(u)
		if err != nil {
			return "", "", err
		}
	} else {
		opts := totp.GenerateOpts{
			Issuer:      TwoFaIssuer,
			AccountName: u.Email,
			SecretSize:  10,
		}

		otpKey, err = totp.Generate(opts)
		if err != nil {
			return "", "", err
		}
	}

	qrCodeURL = s.getQRCodeURL(otpKey)
	return otpKey.Secret(), qrCodeURL, nil

}

func (s *twoFaManager) getOtpKeyForUser(u User) (*otp.Key, error) {
	v := url.Values{}
	v.Set("issuer", TwoFaIssuer)
	v.Set("period", strconv.FormatUint(uint64(30), 10))
	v.Set("algorithm", otp.AlgorithmSHA1.String())
	v.Set("digits", otp.DigitsSix.String())
	v.Set("secret", u.Google2faSecretCode.String)

	otpURL := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + TwoFaIssuer + ":" + u.Email,
		RawQuery: v.Encode(),
	}

	return otp.NewKeyFromURL(otpURL.String())
}

func (s *twoFaManager) getQRCodeURL(otpKey *otp.Key) string {
	size := "200x200"

	return fmt.Sprintf(
		"https://chart.googleapis.com/chart?chs=%s&chld=M|0&cht=qr&chl=%s",
		size,
		url.QueryEscape(otpKey.URL()),
	)
}

func NewTwoFaManager() TwoFaManager {
	return &twoFaManager{}
}
