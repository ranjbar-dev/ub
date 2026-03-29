// Package platform provides infrastructure adapters for external services.
// It defines interfaces and implementations for:
//   - Configuration (Viper wrapper with env var support)
//   - HTTP client (GET, POST, POST form, Basic Auth)
//   - Logging (Zap structured logging + Sentry error reporting)
//   - Email providers (SendGrid, Mailjet, Mailgun, SMTP via factory pattern)
//   - MongoDB client initialization
//   - RabbitMQ connection management (mutex-protected, lazy connect)
//
// All types are interface-based for testability via dependency injection.
package platform

import (
	"flag"
	"github.com/spf13/viper"
)

const AllowedIpsConfigKey = "communicator.allowed_ips"
const SentryDsnKey = "sentry.dsn"
const EnvConfigKey = "communicator.environment"
const EnvProd = "prod"
const EnvTest = "test"
const EnvDev = "dev"
const SmsUrlPrefix = "https://api.twilio.com/2010-04-01/Accounts/"
const SmsUrlPostfix = "/Messages.json"

// Configs provides typed access to application configuration values.
// Backed by Viper with env var override support (prefix: UBCOMMUNICATOR_).
type Configs interface {
	GetString(name string) string
	GetInt(name string) int
	GetBool(name string) bool
	GetStringSlice(name string) []string
	UnmarshalKey(key string, i interface{}) error
	GetAllowedIps() []string
	GetEnv() string
	GetSentryDsn() string
	GetSmsUrl(sId string) string
}

type configs struct {
	viper *viper.Viper
}

func (c *configs) GetString(name string) string {
	return c.viper.GetString(name)
}

func (c *configs) GetInt(name string) int {
	return c.viper.GetInt(name)
}

func (c *configs) GetBool(name string) bool {
	return c.viper.GetBool(name)
}

func (c *configs) GetStringSlice(name string) []string {
	return c.viper.GetStringSlice(name)
}

func (c *configs) UnmarshalKey(key string, i interface{}) error {
	err := c.viper.UnmarshalKey(key, i)
	return err
}

func (c *configs) GetAllowedIps() []string {
	return c.GetStringSlice(AllowedIpsConfigKey)
}

func (c *configs) GetEnv() string {
	return c.GetString(EnvConfigKey)

}

func (c *configs) GetSentryDsn() string {
	return c.GetString(SentryDsnKey)
}

func (c *configs) GetSmsUrl(sId string) string {
	url := SmsUrlPrefix + sId + SmsUrlPostfix
	return url
}

// NewConfigs wraps a Viper instance in the Configs interface.
// If running under go test, automatically sets environment to "test".
func NewConfigs(viper *viper.Viper) Configs {

	// WHY: flag.Lookup("test.v") detects when running under `go test`.
	// This automatically sets the environment to "test" to disable
	// production-only features (Sentry reporting, etc.) during test runs.
	if flag.Lookup("test.v") != nil {
		viper.Set(EnvConfigKey, EnvTest)
	}
	return &configs{viper}
}
