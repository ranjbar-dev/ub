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

type HTTPClient interface {
	HTTPGet(ctx context.Context, url string, headers map[string]string) ([]byte, http.Header, int, error)
	HTTPPost(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error)
	PostForm(ctx context.Context, url string, data url.Values) (resp *http.Response, err error)
}

type httpClient struct {
	client *http.Client
}

func (h *httpClient) PostForm(ctx context.Context, url string, data url.Values) (resp *http.Response, err error) {
	return h.client.PostForm(url, data)
}

func (h *httpClient) HTTPGet(ctx context.Context, url string, headers map[string]string) ([]byte, http.Header, int, error) {
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

	resp, err := h.client.Do(req)

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

func (h *httpClient) HTTPPost(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error) {
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

	resp, err := h.client.Do(req)

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
	return &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
