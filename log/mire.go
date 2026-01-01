package log

import (
	"context"

	"github.com/Lunar-Chipter/mire/core"
	"github.com/Lunar-Chipter/mire/logger"
)

// defaultLogger is the default logger instance
var defaultLogger *logger.Logger

// Mire is the single public API for logging.
// Zero-allocation logging using []byte parameters only.
//
// Usage:
//
//	Mire(context.Background(), INFO, []byte("message"))
//	Mire(ctx, INFO, []byte("message"), []byte("key"), []byte("value"))
//	Mire(ctx, ERROR, []byte("error"), []byte("code"), []byte("500"))
//
// Parameters:
//   - ctx: context.Context for distributed tracing and context extraction
//   - level: log level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)
//   - msg: log message as []byte (zero-allocation)
//   - keyvals: optional key-value pairs as []byte (zero-allocation, must be even number)
func Mire(ctx context.Context, level Level, msg []byte, keyvals ...[]byte) {
	if defaultLogger == nil {
		defaultLogger = logger.NewDefaultLogger()
	}
	defaultLogger.LogZ(ctx, core.Level(level), msg, keyvals...)
}

// MireConfig allows customization of the default logger.
// Must be called before any logging operations.
func MireConfig(config logger.LoggerConfig) {
	defaultLogger = logger.NewLogger(config)
}

// MireClose closes the default logger and flushes all buffers.
// Call this before application exit to ensure all logs are written.
func MireClose() {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
}

// Level represents the severity level of a log entry.
type Level int

// Log level constants
const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

// String returns the string representation of the log level.
func (l Level) String() string {
	return core.Level(l).String()
}
