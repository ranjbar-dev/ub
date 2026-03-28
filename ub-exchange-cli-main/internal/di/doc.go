// Package di provides the dependency injection container for the entire
// application. It registers ~110 services using the sarulabs/di library
// with lazy initialization (di.App scope).
//
// Service name constants are defined as package-level constants. Exported
// constants (e.g., ConfigService, HTTPServer, LoggerService) are used by
// cmd/ entry points; unexported constants are internal to this package.
//
// Registration is organized by domain via addXxx() helper functions called
// from NewContainer(). Registration order matters — services that depend
// on others must be registered after their dependencies.
//
// Usage from entry points:
//
//	container := di.NewContainer()
//	logger := container.Get(di.LoggerService).(platform.Logger)
//	httpServer := container.Get(di.HTTPServer).(*api.HTTPServer)
package di
