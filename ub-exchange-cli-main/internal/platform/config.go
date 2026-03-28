package platform

import (
	"flag"

	"github.com/spf13/viper"
)

const SentryDsnKey = "sentry.dsn"
const EnvConfigKey = "exchange.environment"
const EnvProd = "prod"
const EnvTest = "test"
const EnvDev = "dev"
const ActiveExternalExchangeConfigKey = "exchange.active_external_exchange"
const DomainConfigKey = "exchange.domain"
const RecaptchaSecretKey = "recaptcha.secretkey"
const RecaptchaSiteKey = "recaptcha.sitekey"
const RecaptchaAndroidSecretKey = "recaptcha.androidsecretkey"
const RecaptchaAndroidSiteKey = "recaptcha.androidsitekey"
const ImagePath = "exchange.imagepath"

// Configs provides centralized, typed access to application configuration backed
// by Viper. It exposes typed getters, environment detection, and convenience
// accessors for exchange-specific settings such as reCAPTCHA keys, Sentry DSN,
// image paths, and the active external exchange.
type Configs interface {
	// GetString returns the configuration value for the given key as a string.
	GetString(name string) string
	// GetInt returns the configuration value for the given key as an int.
	GetInt(name string) int
	// GetBool returns the configuration value for the given key as a bool.
	GetBool(name string) bool
	// GetStringSlice returns the configuration value for the given key as a string slice.
	GetStringSlice(name string) []string
	// UnmarshalKey unmarshals the configuration subtree at key into the provided struct.
	UnmarshalKey(key string, i interface{}) error
	// GetEnv returns the current environment identifier (e.g., "prod", "dev", "test").
	GetEnv() string
	// GetSentryDsn returns the Sentry DSN used for error reporting.
	GetSentryDsn() string
	// Set overrides a configuration value at runtime for the given key.
	Set(key string, value interface{})
	// GetActiveExternalExchange returns the identifier of the currently active external exchange.
	GetActiveExternalExchange() string
	// GetDomain returns the configured exchange domain name.
	GetDomain() string
	// GetRecaptchaSecretKey returns the web reCAPTCHA secret key.
	GetRecaptchaSecretKey() string
	// GetRecaptchaSiteKey returns the web reCAPTCHA site key.
	GetRecaptchaSiteKey() string
	// GetAndroidRecaptchaSecretKey returns the Android reCAPTCHA secret key.
	GetAndroidRecaptchaSecretKey() string
	// GetAndroidRecaptchaSiteKey returns the Android reCAPTCHA site key.
	GetAndroidRecaptchaSiteKey() string
	// GetImagePath returns the configured filesystem path for uploaded images.
	GetImagePath() string
	//WriteTempConfig() error
}

type configs struct {
	viper *viper.Viper
}

func (c *configs) GetActiveExternalExchange() string {
	return c.GetString(ActiveExternalExchangeConfigKey)
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

func NewConfigs(viper *viper.Viper) Configs {

	c := &configs{viper}

	if flag.Lookup("test.v") != nil {
		c.Set(EnvConfigKey, EnvTest)
	}

	return c
}

func (c *configs) GetEnv() string {
	return c.GetString(EnvConfigKey)

}

func (c *configs) GetSentryDsn() string {
	return c.GetString(SentryDsnKey)
}

func (c *configs) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

func (c *configs) GetDomain() string {
	return c.GetString(DomainConfigKey)
}

func (c *configs) GetRecaptchaSecretKey() string {
	return c.GetString(RecaptchaSecretKey)
}

func (c *configs) GetRecaptchaSiteKey() string {
	return c.GetString(RecaptchaSiteKey)
}

func (c *configs) GetAndroidRecaptchaSecretKey() string {
	return c.GetString(RecaptchaAndroidSecretKey)
}

func (c *configs) GetAndroidRecaptchaSiteKey() string {
	return c.GetString(RecaptchaAndroidSiteKey)
}

func (c *configs) GetImagePath() string {
	return c.GetString(ImagePath)
}

//func (c *configs) WriteTempConfig() error {
//	return c.viper.WriteConfig()
//}
