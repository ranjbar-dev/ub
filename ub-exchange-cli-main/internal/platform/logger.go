package platform

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

const LogFilePath = "/app/var/logs/exchange-go.log"
const TestLogFilePath = "./../var/logs/exchange-go.log"

// Logger provides structured logging via uber/zap with integrated Sentry error
// tracking. Errors logged through Error2 are automatically forwarded to Sentry
// in production unless the error implements NonSentryErr.
type Logger interface {
	// String creates a zap.Field for a string key-value pair, for use as a structured log field.
	String(key string, val string) zap.Field
	// Info logs a message at the Info level with optional structured fields.
	Info(msg string, fields ...zap.Field)
	// Warn logs a message at the Warn level with optional structured fields.
	Warn(msg string, fields ...zap.Field)
	// Panic logs a message at the Panic level and then panics.
	Panic(msg string, fields ...zap.Field)
	// Fatal logs a message at the Fatal level and then calls os.Exit(1).
	Fatal(msg string, fields ...zap.Field)
	// Error2 logs a message at the Error level and reports the error to Sentry
	// (in production) unless the error implements NonSentryErr.
	Error2(msg string, err error, fields ...zap.Field)
}

//var sentry

type logger struct {
	logger    *zap.Logger
	sentryDsn string
	env       string
}

func NewLogger(configs Configs) Logger {
	cfg := zap.NewProductionConfig()
	env := configs.GetEnv()
	logPath := LogFilePath

	if env == EnvTest {
		logPath = TestLogFilePath
	}

	cfg.OutputPaths = []string{
		logPath,
	}

	zapLogger, err := cfg.Build()

	if err != nil {
		panic("can not set logger")
	}

	sentryErr := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: configs.GetSentryDsn(),
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
	})

	if sentryErr != nil {
		panic("can not set sentry")
	}

	//defer logger.Sync()
	//logger.Info("failed to fetch URL",
	//	// Structured context as strongly typed Field values.
	//	zap.String("url", url),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second),
	//)

	sentryDsn := configs.GetSentryDsn()
	return &logger{zapLogger, sentryDsn, env}

}

func (l *logger) String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func (l *logger) error(err error) zap.Field {
	return zap.Error(err)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *logger) Error2(msg string, err error, fields ...zap.Field) {
	allFields := append([]zap.Field{zap.Error(err)}, fields...)
	l.logger.Error(msg, allFields...)

	l.captureError(err)
}

func (l *logger) captureError(err error) {
	if l.env != "prod" {
		return
	}

	_, ok := err.(NonSentryErr)
	if ok {
		return
	}

	//here we set sentry too
	sentryErr := sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: l.sentryDsn,
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: false,
	})

	if sentryErr != nil {
		l.logger.Error("captureError: failed to re-initialize sentry", zap.Error(sentryErr))
		return
	}

	defer sentry.Flush(2 * time.Second)
	sentry.CurrentHub().CaptureException(err)
}
