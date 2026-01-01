// Package example demonstrates advanced usage patterns of Mire
package example

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/formatter"
	"github.com/Luvion1/mire/hook"
	"github.com/Luvion1/mire/logger"
	"github.com/Luvion1/mire/util"
)

// AdvancedExample demonstrates more complex usage patterns of Mire
func AdvancedExample() {
	fmt.Println("=== Advanced Mire Usage Examples ===")

	contextAwareExample()
	customFormatterExample()
	hookExample()
	performanceExample()

	fmt.Println("Advanced examples completed")
}

func contextAwareExample() {
	fmt.Println("\n--- Context-Aware Logging Example ---")

	ctx := context.Background()
	ctx = util.WithTraceID(ctx, "trace-abc-123")
	ctx = util.WithSpanID(ctx, "span-xyz-789")
	ctx = util.WithUserID(ctx, "user-john-doe")
	ctx = util.WithRequestID(ctx, "req-456-def")

	log := logger.NewLogger(logger.LoggerConfig{
		Level: core.DEBUG,
		Output: os.Stdout,
		Formatter: &formatter.JSONFormatter{
			TimestampFormat: logger.DEFAULT_TIMESTAMP_FORMAT,
			ShowTrace:    true,
			IncludeStackTrace: true,
		},
	})
	defer log.Close()

	log.LogZC(ctx, core.INFO, []byte("Processing user request"))

	log.LogZ(ctx, core.INFO,
		[]byte("User action completed"),
		[]byte("service"), []byte("auth-service"),
		[]byte("action"), []byte("login"),
		[]byte("success"), []byte("true"))
}

func customFormatterExample() {
	fmt.Println("\n--- Custom Formatter with Field Transformers Example ---")

	jsonFormatter := &formatter.JSONFormatter{
		TimestampFormat: logger.DEFAULT_TIMESTAMP_FORMAT,
		ShowCaller:      true,
		SensitiveFields: []string{"password", "token", "secret"},
		MaskSensitiveData: true,
	}

	log := logger.NewLogger(logger.LoggerConfig{
		Level:  core.INFO,
		Output: os.Stdout,
		Formatter: jsonFormatter,
	})
	defer log.Close()

	log.LogZ(context.Background(), core.INFO,
		[]byte("Secure transaction processed"),
		[]byte("user"), []byte("jane.doe@example.com"),
		[]byte("credit_card"), []byte("1234567890123456"),
		[]byte("password"), []byte("supersecret123"),
		[]byte("action"), []byte("purchase"),
		[]byte("amount"), []byte("199.99"))
}

func hookExample() {
	fmt.Println("\n--- Hook Implementation Example ---")

	customHook := &CustomHTTPHook{
		endpoint: "https://logs.example.com/api/logs",
	}

	log := logger.NewLogger(logger.LoggerConfig{
		Level: core.WARN,
		Output: os.Stdout,
		Formatter: &formatter.JSONFormatter{
			TimestampFormat: logger.DEFAULT_TIMESTAMP_FORMAT,
		},
		Hooks: []hook.Hook{customHook},
	})
	defer log.Close()

	log.LogZ(context.Background(), core.WARN,
		[]byte("Database connection issue, retrying..."),
		[]byte("error_code"), []byte("E500"),
		[]byte("component"), []byte("database"))

	_ = customHook.Close()
}

func performanceExample() {
	fmt.Println("\n--- Performance Optimization Example ---")

	perfLog := logger.NewLogger(logger.LoggerConfig{
		Level:          core.INFO,
		Output:         os.Stdout,
		AsyncMode:      true,
		WorkerCount:    8,
		ChannelSize:    5000,
		ProcessTimeout: 5 * time.Second,
		NoTimeout:      true,
		BufferSize:     4096,
		FlushInterval:  50 * time.Millisecond,
		NoLocking:      true,
		IncludeStackTrace: false,
		Formatter: &formatter.CSVFormatter{
			IncludeHeader:   false,
			SensitiveFields: []string{"password", "token"},
			MaskSensitiveData: true,
		},
	})
	defer perfLog.Close()

	start := time.Now()
	for i := 0; i < 1000; i++ {
		perfLog.LogZ(context.Background(), core.INFO,
			[]byte("Performance test message"),
			[]byte("iteration"), []byte(fmt.Sprintf("%d", i)),
			[]byte("worker"), []byte("perf-worker"))
	}
	duration := time.Since(start)

	fmt.Printf("Logged 1000 messages in %v using async logging\n", duration)

	perfLog.Close()
}

// CustomHTTPHook is a custom hook implementation
type CustomHTTPHook struct {
	endpoint string
}

func (h *CustomHTTPHook) Fire(entry *core.LogEntry) error {
	_ = entry
	fmt.Printf("Simulating sending log to %s\n", h.endpoint)
	return nil
}

func (h *CustomHTTPHook) Close() error {
	fmt.Println("CustomHTTPHook closed")
	return nil
}
