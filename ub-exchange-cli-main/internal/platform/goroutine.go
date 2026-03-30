package platform

import (
	"fmt"
	"runtime/debug"

	"go.uber.org/zap"
)

// SafeGo runs fn in a goroutine with panic recovery.
// If fn panics, the panic is logged with a stack trace instead of crashing the process.
func SafeGo(logger Logger, name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error2(
					fmt.Sprintf("goroutine %q panicked", name),
					fmt.Errorf("panic: %v", r),
					zap.String("stack", string(debug.Stack())),
				)
			}
		}()
		fn()
	}()
}
