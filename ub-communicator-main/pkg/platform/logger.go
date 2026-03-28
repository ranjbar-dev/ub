package platform

import (
	"os"
	"path/filepath"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Logger provides structured logging with Sentry integration.
// All methods accept optional zap.Field arguments for structured context.
// Use zap.Error(err) to attach an error to any log message.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	logger    *zap.Logger
	sentryDsn string
	env       string
}

// NewLogger creates a Logger backed by zap with Sentry error reporting.
// Falls back to stdout if the configured log file path is inaccessible.
func NewLogger(configs Configs) Logger {
	logPath := configs.GetString("logging.file_path")

	cfg := zap.NewProductionConfig()

	if logPath == "" || logPath == "stdout" {
		cfg.OutputPaths = []string{"stdout"}
	} else {
		// Ensure log directory exists
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// Fall back to stdout if directory creation fails
			cfg.OutputPaths = []string{"stdout"}
		} else {
			cfg.OutputPaths = []string{"stdout", logPath}
		}
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		// Fall back to basic logger if production config fails
		fallback, _ := zap.NewDevelopment()
		if fallback == nil {
			// Last resort: nop logger (should never happen)
			fallback = zap.NewNop()
		}
		zapLogger = fallback
	}

	sentryDsn := configs.GetSentryDsn()
	env := configs.GetEnv()

	if sentryDsn != "" {
		sentryErr := sentry.Init(sentry.ClientOptions{
			Dsn:   sentryDsn,
			Debug: false,
		})
		if sentryErr != nil {
			zapLogger.Warn("failed to initialize sentry, error reporting disabled",
				zap.Error(sentryErr))
		}
	}

	return &logger{zapLogger, sentryDsn, env}
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)

	// Extract error from fields for Sentry reporting
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			if err, ok := f.Interface.(error); ok {
				l.captureError(err)
			}
		}
	}
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *logger) captureError(err error) {
	// WHY: Sentry error reporting is disabled outside production to avoid
	// polluting the error dashboard with development/test noise.
	if l.env != "prod" {
		return
	}
	sentry.CaptureException(err)
	sentry.Flush(2 * time.Second)
}
