# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.7] - 2025-12-28

### Added

- **NEW: Microservices Example** - Complete microservice logging patterns
  - Request/response logging with distributed tracing
  - Business events logging
  - Metrics collection
  - Database operation logging
  - Cache operation logging
  - External service call logging
  - HTTP middleware for automatic logging
  - Environment-based configuration

### Changed

- **Examples Updated** - All example files now use \`log.Mire()\` and \`log.LogZ()\` API
  - example.go - Basic usage with zero-allocation
  - efficient_usage.go - Zero-allocation patterns
  - advanced_example.go - Context-aware, hooks, performance
  - external_example.go - External service integration
  - **NEW** microservices.go - Complete microservice patterns

- All examples demonstrate:
  - Global API: \`log.Mire()\`
  - Logger API: \`log.LogZ()\`, \`log.LogZC()\`
  - Zero-allocation with \`[]byte\` parameters

## [v0.0.6] - 2025-12-28

### BREAKING CHANGES

- **Single Public API**: Reduced entire public API surface to one function - `log.Mire()`
  - Removed all `Info()`, `Debug()`, `Warn()`, `Error()` methods
  - Removed all `Infof()`, `Debugf()`, `Warnf()`, `Errorf()` formatted methods
  - Removed all context-aware methods like `InfoC()`, `DebugC()`, etc.
  - Only API is now: `log.Mire(ctx context.Context, level Level, msg []byte, keyvals ...[]byte)`

### Added

- `log.Mire(ctx, level, msg, keyvals...)` - Single public API for all logging
- `log.MireConfig(config)` - Configure default logger
- `log.MireClose()` - Close logger and flush buffers
- `log.TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC` - Level constants
- Strong Mire identity branding

### Changed

- **Zero-Allocation Only**: All public APIs now use `[]byte` only - no `string` or `interface{}`
- Simplified API surface from 20+ functions to 3 functions + level constants
- Updated all documentation and examples to use new `log.Mire()` API
- Updated main.go demonstration to use new single API

### Removed

- `log.Info()`, `log.Debug()`, `log.Warn()`, `log.Error()`, `log.Fatal()`, `log.Panic()`
- `log.Trace()`, `log.Tracef()`, `log.Debugf()`, `log.Infof()`, `log.Warnf()`, `log.Errorf()`, `log.Fatalf()`, `log.Panicf()`
- `log.TraceC()`, `log.DebugC()`, `log.InfoC()`, `log.WarnC()`, `log.ErrorC()`
- `log.TracefC()`, `log.DebugfC()`, `log.InfofC()`, `log.WarnfC()`, `log.ErrorfC()`, `log.FatalfC()`, `log.PanicfC()`
- `log.WithFields()` - Use `log.Mire()` with keyvals instead
- `log.DefaultLogger()` - Internal implementation detail
- `log.SetDefault()` - Internal implementation detail

### Fixed

- Updated all tests to work with new single API
- Updated README with new API documentation

### Migration Guide

If you're upgrading from v0.0.5, you need to update your code:

**Before (v0.0.5):**
```go
log.Info("message")
log.Infof("User %s logged in", "alice")
log.Info("message", "user_id", 12345)
log.InfoC(ctx, "message with context")
```

**After (vv0.0.6):**
```go
log.Mire(ctx, log.INFO, []byte("message"))
log.Mire(ctx, log.INFO, []byte("User alice logged in"))
log.Mire(ctx, log.INFO,
    []byte("message"),
    []byte("user_id"), []byte("12345"))
log.Mire(ctx, log.INFO, []byte("message with context"))
```

## [0.0.5] - 2025-12-28

### Added

- Zero-allocation API with automatic type conversion
- Multiple formatters (Text, JSON, CSV)
- Asynchronous logging
- Log rotation (size and time-based)
- Hook system
- Metrics collection
- Distributed tracing support (trace_id, span_id, user_id, etc.)
- Context-aware logging
- Sensitive data masking

### Performance

- 1M+ logs/second with zero-allocation design
- TextFormatter: 1,390ns/op, 72 B/op, 2 allocs/op
- JSONFormatter: 2,846ns/op, 24 B/op, 1 alloc/op
- CSVFormatter: 2,990ns/op, 48 B/op, 2 allocs/op
