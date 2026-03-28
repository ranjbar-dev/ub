package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient provides outbound HTTP capabilities for calling external APIs.
// It supports GET, POST (JSON body), and form-encoded POST requests with custom headers.
type HTTPClient interface {
	// HTTPGet performs an HTTP GET request to the given URL with optional headers.
	// It returns the response body, headers, status code, and any error encountered.
	HTTPGet(ctx context.Context, url string, headers map[string]string) ([]byte, http.Header, int, error)
	// HTTPPost performs an HTTP POST request with a JSON-serialized body and optional headers.
	// If body is already a []byte it is sent as-is; otherwise it is marshaled to JSON.
	// It returns the response body, headers, status code, and any error encountered.
	HTTPPost(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error)
	// PostForm performs an HTTP POST request with URL-encoded form data.
	PostForm(ctx context.Context, url string, data url.Values) (resp *http.Response, err error)
}

type httpClient struct {
}

func (httpClient *httpClient) PostForm(ctx context.Context, url string, data url.Values) (resp *http.Response, err error) {
	return http.PostForm(url, data)
}

func (httpClient *httpClient) HTTPGet(ctx context.Context, url string, headers map[string]string) ([]byte, http.Header, int, error) {
	respBody := []byte("")
	header := http.Header{}
	statusCode := http.StatusBadRequest
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return respBody, header, statusCode, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	http.DefaultClient.Timeout = 10 * time.Second

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return respBody, header, statusCode, err
	}

	defer resp.Body.Close()

	header = resp.Header
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return respBody, header, statusCode, err
	}
	statusCode = resp.StatusCode

	return respBody, header, statusCode, nil

}

func (httpClient *httpClient) HTTPPost(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error) {
	finalRes := []byte("")
	statusCode := http.StatusBadRequest
	header := http.Header{}

	requestBody := []byte("")
	if b, ok := body.([]byte); ok {
		requestBody = b
	} else {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return finalRes, header, statusCode, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return finalRes, header, statusCode, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return finalRes, header, statusCode, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return finalRes, header, statusCode, err
	}
	finalRes = respBody
	header = resp.Header
	statusCode = resp.StatusCode
	return finalRes, header, statusCode, err
}

func NewHTTPClient() HTTPClient {
	return &httpClient{}
}
