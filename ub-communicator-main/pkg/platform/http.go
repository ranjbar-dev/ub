package platform

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

// HttpClient provides HTTP request methods with Basic Auth support.
// All methods return (body, headers, statusCode, error).
type HttpClient interface {
	HttpGet(url string) ([]byte, http.Header, int, error)
	HttpPost(url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error)
	HttpPostForm(url string, body *strings.Reader, headers map[string]string) ([]byte, http.Header, int, error)
	BasicAuth(username, password string) string
}

type httpClient struct {
	client *http.Client
}

func (c *httpClient) HttpGet(url string) ([]byte, http.Header, int, error) {
	respBody := []byte("")
	header := http.Header{}
	statusCode := 0

	resp, err := c.client.Get(url)

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

func (c *httpClient) HttpPost(url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error) {
	finalRes := []byte("")
	statusCode := 0
	header := http.Header{}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return finalRes, header, statusCode, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	if err != nil {
		return finalRes, header, statusCode, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)

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

func (c *httpClient) HttpPostForm(url string, body *strings.Reader, headers map[string]string) ([]byte, http.Header, int, error) {
	finalRes := []byte("")
	statusCode := 0
	header := http.Header{}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return finalRes, header, statusCode, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.client.Do(req)

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

func (c *httpClient) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// NewHttpClient creates an HttpClient with a shared connection pool and 10s timeout.
func NewHttpClient() HttpClient {
	return &httpClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
