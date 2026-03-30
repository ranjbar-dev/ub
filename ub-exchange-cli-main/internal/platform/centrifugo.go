package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// CentrifugoClient provides real-time message publishing to connected clients
// via the Centrifugo server HTTP API. Messages are published to channels that
// clients subscribe to for live tickers, order updates, and trade notifications.
type CentrifugoClient interface {
	// Publish sends data to a single Centrifugo channel.
	Publish(channel string, data interface{}) error
	// Broadcast sends data to multiple Centrifugo channels at once.
	Broadcast(channels []string, data interface{}) error
}

type centrifugoClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
	logger  Logger
	env     string
}

func NewCentrifugoClient(configs Configs, logger Logger) CentrifugoClient {
	// Expect centrifugo.api_url to include /api (e.g. http://centrifugo:8000/api),
	// matching the phpcent convention used by ub-server-main.
	baseURL := strings.TrimRight(configs.GetString("centrifugo.api_url"), "/")
	return &centrifugoClient{
		baseURL: baseURL,
		apiKey:  configs.GetString("centrifugo.api_key"),
		client:  &http.Client{Timeout: 5 * time.Second},
		logger:  logger,
		env:     configs.GetEnv(),
	}
}

func (c *centrifugoClient) Publish(channel string, data interface{}) error {
	if c.env == EnvTest {
		return nil
	}
	payload := map[string]interface{}{
		"channel": channel,
		"data":    data,
	}
	return c.doRequest("/publish", payload)
}

func (c *centrifugoClient) Broadcast(channels []string, data interface{}) error {
	if c.env == EnvTest {
		return nil
	}
	payload := map[string]interface{}{
		"channels": channels,
		"data":     data,
	}
	return c.doRequest("/broadcast", payload)
}

func (c *centrifugoClient) doRequest(path string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("centrifugo marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("centrifugo request error: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Warn("centrifugo publish error",
			zap.Error(err),
			zap.String("service", "centrifugoClient"),
			zap.String("method", "doRequest"),
		)
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("centrifugo API error: status %d", resp.StatusCode)
	}
	return nil
}
