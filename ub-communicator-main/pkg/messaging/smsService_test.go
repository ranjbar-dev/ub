package messaging

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"testing"
)

// mockHttpClient captures the request and returns configurable responses.
type mockHttpClient struct {
	lastURL     string
	lastPayload string
	lastHeaders map[string]string

	returnBody   []byte
	returnStatus int
	returnErr    error
}

func (m *mockHttpClient) HttpPostForm(url string, body *strings.Reader, headers map[string]string) ([]byte, http.Header, int, error) {
	m.lastURL = url
	// Read the body content
	if body != nil {
		buf := make([]byte, 1024)
		n, _ := body.Read(buf)
		m.lastPayload = string(buf[:n])
	}
	m.lastHeaders = headers
	return m.returnBody, nil, m.returnStatus, m.returnErr
}

func (m *mockHttpClient) HttpGet(url string) ([]byte, http.Header, int, error) {
	return nil, nil, 0, nil
}

func (m *mockHttpClient) HttpPost(url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error) {
	return nil, nil, 0, nil
}

func (m *mockHttpClient) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// mockSmsConfigs returns canned configuration values.
type mockSmsConfigs struct{}

func (m *mockSmsConfigs) GetString(key string) string {
	configs := map[string]string{
		"sms.account_sid": "AC1234567890",
		"sms.auth_token":  "test_token_abc",
		"sms.from":        "+15551234567",
	}
	return configs[key]
}

func (m *mockSmsConfigs) GetInt(name string) int {
	return 0
}

func (m *mockSmsConfigs) GetBool(name string) bool {
	return false
}

func (m *mockSmsConfigs) GetStringSlice(name string) []string {
	return []string{}
}

func (m *mockSmsConfigs) UnmarshalKey(key string, i interface{}) error {
	return nil
}

func (m *mockSmsConfigs) GetAllowedIps() []string {
	return []string{}
}

func (m *mockSmsConfigs) GetEnv() string {
	return "test"
}

func (m *mockSmsConfigs) GetSentryDsn() string {
	return ""
}

func (m *mockSmsConfigs) GetSmsUrl(sId string) string {
	return "https://api.twilio.com/2010-04-01/Accounts/" + sId + "/Messages.json"
}

// Test 1: Auth Header Format — No Double Space
func TestSms_AuthHeaderFormat(t *testing.T) {
	httpClient := &mockHttpClient{returnBody: []byte(`{"sid":"SM123"}`), returnStatus: 201}
	configs := &mockSmsConfigs{}
	svc := NewSmsService(httpClient, configs)

	svc.Send("Test Subject", "+1234567890", "Test message")

	auth, ok := httpClient.lastHeaders["Authorization"]
	if !ok {
		t.Fatal("Authorization header missing")
	}

	// Must be "Basic <base64>" with exactly ONE space
	parts := strings.SplitN(auth, " ", 3)
	if len(parts) != 2 {
		t.Errorf("malformed auth header: %q (expected 'Basic <token>')", auth)
	}
	if parts[0] != "Basic" {
		t.Errorf("auth scheme: got %q, want %q", parts[0], "Basic")
	}

	// Verify base64 decodes to SID:TOKEN
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	expected := "AC1234567890:test_token_abc"
	if string(decoded) != expected {
		t.Errorf("decoded credentials: got %q, want %q", string(decoded), expected)
	}
}

// Test 2: Request URL Contains Account SID
func TestSms_RequestURL(t *testing.T) {
	httpClient := &mockHttpClient{returnBody: []byte(`{"sid":"SM123"}`), returnStatus: 201}
	configs := &mockSmsConfigs{}
	svc := NewSmsService(httpClient, configs)

	svc.Send("Test Subject", "+1234567890", "Test")

	expectedURL := "https://api.twilio.com/2010-04-01/Accounts/AC1234567890/Messages.json"
	if httpClient.lastURL != expectedURL {
		t.Errorf("URL: got %q, want %q", httpClient.lastURL, expectedURL)
	}
}

// Test 3: Request Body Contains To, From, Body
func TestSms_RequestBody(t *testing.T) {
	httpClient := &mockHttpClient{returnBody: []byte(`{"sid":"SM123"}`), returnStatus: 201}
	configs := &mockSmsConfigs{}
	svc := NewSmsService(httpClient, configs)

	svc.Send("Test Subject", "+1234567890", "Hello World")

	body := httpClient.lastPayload
	if !strings.Contains(body, "To=%2B1234567890") && !strings.Contains(body, "To=+1234567890") {
		t.Errorf("body missing To: %q", body)
	}
	if !strings.Contains(body, "From=%2B15551234567") && !strings.Contains(body, "From=+15551234567") {
		t.Errorf("body missing From: %q", body)
	}
	if !strings.Contains(body, "Body=Hello+World") && !strings.Contains(body, "Body=Hello%20World") {
		t.Errorf("body missing Body: %q", body)
	}
}

// Test 4: HTTP Error Returns (false, error)
func TestSms_HttpError(t *testing.T) {
	httpClient := &mockHttpClient{
		returnBody:   nil,
		returnStatus: 0,
		returnErr:    errors.New("connection refused"),
	}
	configs := &mockSmsConfigs{}
	svc := NewSmsService(httpClient, configs)

	ok, err := svc.Send("Test Subject", "+1234567890", "Test")

	if ok {
		t.Error("expected ok=false on HTTP error")
	}
	if err == nil {
		t.Error("expected non-nil error on HTTP error")
	}
}

// Test 5: Non-201 Status Returns (false, error)
func TestSms_Non201Status(t *testing.T) {
	httpClient := &mockHttpClient{
		returnBody:   []byte(`{"message":"Authentication Error","code":20003}`),
		returnStatus: 401,
	}
	configs := &mockSmsConfigs{}
	svc := NewSmsService(httpClient, configs)

	ok, err := svc.Send("Test Subject", "+1234567890", "Test")

	if ok {
		t.Error("expected ok=false on 401 status")
	}
	if err == nil {
		t.Error("expected non-nil error on 401 status")
	}
}
