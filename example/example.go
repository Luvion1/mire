// Package example demonstrates basic usage of Mire logging library
package example

import (
	"context"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/logger"
	"github.com/Luvion1/mire/util"
)

// BasicExample demonstrates simple logging with Mire
func BasicExample() {
	// Create logger with text formatter
	log := logger.NewDefault()
	defer log.Close()

	// Basic logging
	log.LogZ(context.Background(), core.INFO, []byte("Application started"))

	// With fields
	log.LogZ(context.Background(), core.DEBUG,
		[]byte("Debug message"),
		[]byte("key1"), []byte("value1"),
		[]byte("key2"), []byte("value2"))

	// With context
	ctx := context.Background()
	ctx = util.WithTraceID(ctx, "trace-abc-123")
	log.LogZ(ctx, core.INFO, []byte("Context-aware message"))

	// All log levels
	log.LogZ(context.Background(), core.TRACE, []byte("Trace message"))
	log.LogZ(context.Background(), core.DEBUG, []byte("Debug message"))
	log.LogZ(context.Background(), core.INFO, []byte("Info message"))
	log.LogZ(context.Background(), core.WARN, []byte("Warning message"))
	log.LogZ(context.Background(), core.ERROR, []byte("Error message"))
	// FATAL and PANIC would exit the application
}
