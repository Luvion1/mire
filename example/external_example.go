package example

import (
	"context"
	"os"

	"github.com/Lunar-Chipter/mire/core"
	"github.com/Lunar-Chipter/mire/formatter"
	"github.com/Lunar-Chipter/mire/logger"
)

type key string

const traceIDKey key = "trace_id"

// ExternalServiceExample demonstrates how to use external hook
func ExternalServiceExample() {
	// Create a logger with JSON formatter for production
	log := logger.NewLogger(logger.LoggerConfig{
		Level:  core.INFO,
		Output: os.Stdout,
		Formatter: &formatter.JSONFormatter{
			TimestampFormat: logger.DEFAULT_TIMESTAMP_FORMAT,
			ShowCaller:    true,
		},
	})
	defer log.Close()

	// Log application lifecycle
	log.LogZ(context.Background(), core.INFO, []byte("Application started"))

	// Simulate some processing with context
	ctx := context.Background()
	ctx = context.WithValue(ctx, traceIDKey, "trace-12345")

	log.LogZ(ctx, core.INFO, []byte("Connecting to external service"))

	// Log successful connection
	log.LogZ(ctx, core.INFO, []byte("External service connected successfully"))

	// Log some operations
	log.LogZ(ctx, core.DEBUG,
		[]byte("Processing request"),
		[]byte("request_id"), []byte("req-67890"))

	log.LogZ(ctx, core.DEBUG,
		[]byte("Request completed successfully"),
		[]byte("request_id"), []byte("req-67890"),
		[]byte("duration_ms"), []byte("150"))

	// Simulate a warning
	log.LogZ(context.Background(), core.WARN,
		[]byte("This is a warning that might indicate a problem"))

	// Simulate an error
	log.LogZ(context.Background(), core.ERROR,
		[]byte("This is an error that occurred during processing"))

	// Cleanup
	log.LogZ(context.Background(), core.INFO, []byte("Application shutting down"))
}
