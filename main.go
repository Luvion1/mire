package main

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/formatter"
	"github.com/Luvion1/mire/log"
	"github.com/Luvion1/mire/logger"
	"github.com/Luvion1/mire/util"
)

// wrappedError wraps an error with a message
type wrappedError struct {
	msg   string
	cause error
}

func (e *wrappedError) Error() string {
	if e.cause != nil {
		return e.msg + ": " + e.cause.Error()
	}
	return e.msg
}

func (e *wrappedError) Unwrap() error {
	return e.cause
}

// printLine is a helper function to print lines without fmt
func printLine(s string) {
	_, _ = os.Stdout.Write([]byte(s))
	_, _ = os.Stdout.Write([]byte("\n"))
}

func main() {
	printLine("===================================================")
	printLine("  MIRE LOGGING LIBRARY DEMONSTRATION")
	printLine("===================================================")

	// Example 1: Default Logger using Mire API
	printLine("### 1. Default Logger with Mire API ###")
	ctx := context.Background()

	// Simple message
	log.Mire(ctx, log.INFO, []byte("Simple INFO message"))

	// With fields - zero-allocation
	log.Mire(ctx, log.WARN,
		[]byte("There are 2 warnings in the system."),
		[]byte("count"), []byte("2"))

	// Debug and Trace won't appear due to WARN level
	log.Mire(ctx, log.DEBUG, []byte("This debug message will NOT appear"))
	log.Mire(ctx, log.TRACE, []byte("This trace message will NOT appear"))

	// Error logging
	log.Mire(ctx, log.ERROR,
		[]byte("A simple error occurred."),
		[]byte("severity"), []byte("high"))
	printLine("---------------------------------------------------")
	time.Sleep(10 * time.Millisecond)

	// --- Example 2: Logger with Fields and Context ---
	printLine("### 2. Logger with Fields & Context ###")
	// Adding TraceID, SpanID, UserID to context.
	// These will be automatically extracted by the logger.
	ctx = util.WithTraceID(ctx, "trace-xyz-987")
	ctx = util.WithSpanID(ctx, "span-123")
	ctx = util.WithUserID(ctx, "user-alice")

	// Mire API with fields
	log.Mire(ctx, log.INFO,
		[]byte("User successfully logged in."),
		[]byte("service"), []byte("auth-service"),
		[]byte("version"), []byte("1.0.0"),
		[]byte("username"), []byte("alice"),
		[]byte("ip_address"), []byte("192.168.1.100"),
	)

	// Context-aware logging
	log.Mire(ctx, log.INFO,
		[]byte("Processing authorization request for token-ABC."))
	printLine("---------------------------------------------------")
	time.Sleep(10 * time.Millisecond)

	// --- Example 3: Error Logging with Stack Trace ---
	printLine("### 3. Error Logging with Stack Trace ###")
	errSample := errors.New("failed to read database configuration")
	log.Mire(ctx, log.ERROR,
		[]byte("Error during initialization: "+errSample.Error()),
		[]byte("error_code"), []byte("500"),
		[]byte("component"), []byte("database-connector"),
	)
	printLine("---------------------------------------------------")
	time.Sleep(10 * time.Millisecond)

	// --- Example 4: JSON Logger to File (app.log) using custom logger ---
	printLine("### 4. Custom JSON Logger to File (app.log) ###")
	printLine("JSON logs will be written to 'app.log'. Check its contents after program completes.")

	// Configure custom logger for file output
	jsonConfig := logger.LoggerConfig{
		Level:       core.DEBUG,
		Output:      mustOpenFile("app.log"),
		ErrorOutput: io.Discard,
		BufferSize:  1024,
		Formatter: &formatter.JSONFormatter{
			PrettyPrint:       true,
			ShowCaller:        true,
			IncludeStackTrace: true,
			TimestampFormat:   logger.DEFAULT_TIMESTAMP_FORMAT,
		},
	}
	jsonLogger := logger.NewLogger(jsonConfig)
	defer jsonLogger.Close()

	// Use custom logger directly with LogZ
	jsonLogger.LogZ(ctx, core.DEBUG, []byte("Debug message for JSON file logger."))
	jsonLogger.LogZ(ctx, core.INFO,
		[]byte("Transaction processed successfully."),
		[]byte("trans_id"), []byte("TXN-001"),
		[]byte("amount"), []byte("123.45"),
		[]byte("currency"), []byte("IDR"),
	)
	jsonLogger.LogZ(ctx, core.ERROR,
		[]byte("Failed to save user data to cache."),
		[]byte("user_id"), []byte("user-bob"),
		[]byte("cache_key"), []byte("user:bob"),
	)
	printLine("---------------------------------------------------")
	time.Sleep(10 * time.Millisecond)

	// --- Example 5: Custom Text Logger (Without Timestamp & Caller) ---
	printLine("### 5. Custom Text Logger ###")
	customTextLogger := logger.NewLogger(logger.LoggerConfig{
		Level:       core.TRACE,
		ErrorOutput: io.Discard,
		Formatter: &formatter.TextFormatter{
			EnableColors:  true,
			ShowTimestamp: false,
			ShowCaller:    false,
			ShowPID:       true,
			ShowGoroutine: true,
		},
	})
	defer customTextLogger.Close()

	customTextLogger.LogZ(ctx, core.Level(log.TRACE),
		[]byte("This is a 'TRACE' message from custom logger (without timestamp/caller)."))
	customTextLogger.LogZ(ctx, core.Level(log.INFO),
		[]byte("INFO level message."),
		[]byte("level_name"), []byte(log.INFO.String()))
	printLine("---------------------------------------------------")
	time.Sleep(10 * time.Millisecond)

	// --- Example 6: All Log Levels ---
	printLine("### 6. All Log Levels Demo ###")
	testCtx := context.Background()
	log.Mire(testCtx, log.TRACE, []byte("TRACE level message"))
	log.Mire(testCtx, log.DEBUG, []byte("DEBUG level message"))
	log.Mire(testCtx, log.INFO, []byte("INFO level message"))
	log.Mire(testCtx, log.WARN, []byte("WARN level message"))
	log.Mire(testCtx, log.ERROR, []byte("ERROR level message"))
	// Note: FATAL and PANIC would exit the application
	printLine("---------------------------------------------------")

	// Ensure all buffers are flushed before program ends.
	printLine("===================================================")
	printLine("  DEMONSTRATION COMPLETED                          ")
	printLine("===================================================")
	printLine("Check 'app.log' and 'errors.log' files to see the log output.")
}

func mustOpenFile(filePath string) io.Writer {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
