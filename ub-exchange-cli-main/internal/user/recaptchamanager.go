package user

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

const GoogleRecaptchaURL = "https://www.google.com/recaptcha/api/siteverify"

// RecaptchaManager validates CAPTCHA responses for anti-bot protection.
// It supports both Google reCAPTCHA and the custom UbCaptcha mechanism.
type RecaptchaManager interface {
	// CheckRecaptcha validates a CAPTCHA response token. It auto-detects whether the
	// token is a Google reCAPTCHA or UbCaptcha response and delegates accordingly.
	// Verification is skipped in non-production environments.
	CheckRecaptcha(recaptchaResponse string, device string, ip string) (bool, error)
}

type recaptchaManager struct {
	httpClient       platform.HTTPClient
	configs          platform.Configs
	logger           platform.Logger
	ubCaptchaManager UbCaptchaManager
}

type recaptchaVerifyResponseBody struct {
	Success     bool    `json:"success"`
	Score       float64 `json:"score"`
	Action      string  `json:"action"`
	ChallengeTs string  `json:"challenge_ts"`
}

func (rm recaptchaManager) CheckRecaptcha(recaptchaResponse string, device string, ip string) (bool, error) {
	if rm.configs.GetEnv() != platform.EnvProd {
		return true, nil
	}

	recaptchaResponse = strings.Trim(recaptchaResponse, "")
	if recaptchaResponse == "" {
		return false, nil
	}

	//check if it is ubCaptcha
	if strings.Contains(recaptchaResponse, UbCaptchaPrefix) {
		return rm.checkUbCaptcha(recaptchaResponse)
	}

	recaptchaSecretKey := rm.configs.GetRecaptchaSecretKey()
	ctx := context.Background()
	formData := url.Values{}
	formData.Set("secret", recaptchaSecretKey)
	formData.Set("response", recaptchaResponse)
	formData.Set("remoteip", ip)
	resp, err := rm.httpClient.PostForm(ctx, GoogleRecaptchaURL, formData)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		rb := recaptchaVerifyResponseBody{}
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			rm.logger.Error2("can not read resp body", err,
				zap.String("service", "recaptchaManager"),
				zap.String("method", "CheckRecaptcha"),
			)
			return false, err
		}

		err = json.Unmarshal(respBody, &rb)
		if err != nil {
			rm.logger.Error2("can not unmarshal recaptcha response", err,
				zap.String("service", "recaptchaManager"),
				zap.String("method", "CheckRecaptcha"),
				zap.String("body", string(respBody)),
			)
			return false, nil
		}

		if rb.Success {
			// todo this score if should be uncommented if we want to have restricted check
			//if rb.Score > 0.3 {
			//	return true,nil
			//}
			//return false, nil

			return true, nil
		}
		return false, nil
	}

	respBody, _ := io.ReadAll(resp.Body)
	err = fmt.Errorf("status code is not 200 it is %d with the body %s", resp.StatusCode, string(respBody))
	rm.logger.Error2("recaptcha response from google is not 200", err,
		zap.String("service", "recaptchaManager"),
		zap.String("method", "CheckRecaptcha"),
	)
	return false, nil
}

func (rm recaptchaManager) checkUbCaptcha(ubCaptchaStr string) (bool, error) {
	return rm.ubCaptchaManager.CheckUbCaptcha(ubCaptchaStr)
}

func NewRecaptchaManager(httpClient platform.HTTPClient, configs platform.Configs, logger platform.Logger, ubCaptchaManager UbCaptchaManager) RecaptchaManager {
	return &recaptchaManager{
		httpClient:       httpClient,
		configs:          configs,
		logger:           logger,
		ubCaptchaManager: ubCaptchaManager,
	}
}
