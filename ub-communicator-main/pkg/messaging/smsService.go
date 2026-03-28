package messaging

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"ub-communicator/pkg/platform"
)

// SmsService sends SMS messages via the Twilio HTTP API.
type SmsService interface {
	Send(subject string, receiver string, content string) (bool, error)
}

type smsService struct {
	httpClient platform.HttpClient
	sId        string
	authToken  string
	smsUrl     string
	from       string
}

func (s *smsService) Send(subject string, receiver string, content string) (bool, error) {
	headers := s.getRequestHeaders()
	body := s.getRequestBody(receiver, content)
	resp, _, statusCode, err := s.httpClient.HttpPostForm(s.smsUrl, &body, headers)
	if err != nil {
		return false, err
	}

	if statusCode >= 200 && statusCode < 300 {
		var data map[string]interface{}
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return false, err
		}
		return true, nil

	} else {
		//bodyString := string(resp)
		return false, fmt.Errorf("sms send failed with status code %d", statusCode)
	}

}

func (s *smsService) getRequestHeaders() map[string]string {
	basicAuthToken := s.httpClient.BasicAuth(s.sId, s.authToken)
	rh := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": "Basic " + basicAuthToken,
		"Accept":        "application/json",
	}
	return rh
}

func (s *smsService) getRequestBody(receiver string, content string) strings.Reader {
	msgData := url.Values{}
	msgData.Set("To", receiver)
	msgData.Set("From", s.from)
	msgData.Set("Body", content)
	msgDataReader := *strings.NewReader(msgData.Encode())
	return msgDataReader
}

// NewSmsService creates an SmsService configured with Twilio credentials from config.
func NewSmsService(httpClient platform.HttpClient, configs platform.Configs) SmsService {
	sId := configs.GetString("sms.account_sid")
	authToken := configs.GetString("sms.auth_token")
	from := configs.GetString("sms.from")
	smsUrl := configs.GetSmsUrl(sId)
	return &smsService{httpClient, sId, authToken, smsUrl, from}
}
