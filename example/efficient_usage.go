package example

import (
	"context"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/logger"
)

// Example of unified API usage with zero-allocation
func EfficientLoggingExample() {
	// Create logger with default configuration
	log := logger.NewDefault()
	defer log.Close()

	// Basic logging with LogZC
	log.LogZC(context.Background(), core.INFO, []byte("This is unified API"))

	// With context
	ctx := context.Background()
	log.LogZ(ctx, core.INFO, []byte("Context-aware logging"), nil)

	// With fields using LogZ (zero-allocation)
	log.LogZ(context.Background(), core.INFO,
		[]byte("User login successful"),
		[]byte("user_id"), []byte("12345"),
		[]byte("action"), []byte("login"),
		[]byte("success"), []byte("true"),
		[]byte("score"), []byte("98.5"))

	// Complete API: context + fields
	log.LogZ(ctx, core.INFO,
		[]byte("Complete logging"),
		[]byte("user_id"), []byte("12345"),
		[]byte("action"), []byte("login"))
}
