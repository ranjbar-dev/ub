package platform_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ub-communicator/pkg/platform"
)

func TestHttpGet_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := platform.NewHttpClient()
	body, _, statusCode, err := client.HttpGet(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != 200 {
		t.Errorf("expected status 200, got %d", statusCode)
	}
	if string(body) != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestHttpGet_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer server.Close()

	client := platform.NewHttpClient()
	_, _, statusCode, err := client.HttpGet(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != 500 {
		t.Errorf("expected status 500, got %d", statusCode)
	}
}

func TestHttpGet_InvalidURL(t *testing.T) {
	client := platform.NewHttpClient()
	_, _, statusCode, err := client.HttpGet("http://invalid.localhost.test:99999")

	if err == nil {
		t.Error("expected error for invalid URL")
	}
	if statusCode != 0 {
		t.Errorf("expected default status 0 on error, got %d", statusCode)
	}
}

func TestHttpPost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":1}`))
	}))
	defer server.Close()

	client := platform.NewHttpClient()
	headers := map[string]string{"Content-Type": "application/json"}
	body := map[string]string{"key": "value"}

	resp, _, statusCode, err := client.HttpPost(server.URL, body, headers)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != 201 {
		t.Errorf("expected status 201, got %d", statusCode)
	}
	if string(resp) != `{"id":1}` {
		t.Errorf("unexpected response: %s", string(resp))
	}
}

func TestHttpPostForm_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"sid":"SM123"}`))
	}))
	defer server.Close()

	client := platform.NewHttpClient()
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	body := strings.NewReader("To=%2B1234567890&From=%2B9876543210&Body=Hello")

	resp, _, statusCode, err := client.HttpPostForm(server.URL, body, headers)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != 200 {
		t.Errorf("expected status 200, got %d", statusCode)
	}
	if string(resp) != `{"sid":"SM123"}` {
		t.Errorf("unexpected response: %s", string(resp))
	}
}

func TestBasicAuth(t *testing.T) {
	client := platform.NewHttpClient()
	result := client.BasicAuth("user", "pass")
	expected := base64.StdEncoding.EncodeToString([]byte("user:pass"))
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestBasicAuth_EmptyCredentials(t *testing.T) {
	client := platform.NewHttpClient()
	result := client.BasicAuth("", "")
	expected := base64.StdEncoding.EncodeToString([]byte(":"))
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
