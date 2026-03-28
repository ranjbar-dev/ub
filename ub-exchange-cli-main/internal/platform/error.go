package platform

// NonSentryErr is a marker interface for errors that should bypass Sentry error
// reporting. Any error type implementing this interface with ShouldSendToSentry
// returning false will be silently logged without being forwarded to Sentry.
type NonSentryErr interface {
	// ShouldSendToSentry returns false to indicate that the error should not be
	// reported to Sentry. Implementations should always return false.
	ShouldSendToSentry() bool
}

type NonSentryError struct {
	Err error
}

func (err NonSentryError) ShouldSendToSentry() bool {
	return false
}

func (err NonSentryError) Error() string {
	return err.Err.Error()
}

type OrderCreateValidationError struct {
	Err error
}

func (err OrderCreateValidationError) Error() string {
	return err.Err.Error()
}

func (OrderCreateValidationError) Is(target error) bool {
	_, ok := target.(OrderCreateValidationError)
	if !ok {
		return false
	}
	return true
}

func (err OrderCreateValidationError) ShouldSendToSentry() bool {
	return false
}
