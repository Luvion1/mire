package util

import (
	"context"
	"sync"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string

const (
	TraceIDKey   contextKey = "trace_id"
	SpanIDKey    contextKey = "span_id"
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
	RequestIDKey contextKey = "request_id"
	ClientIPKey  contextKey = "client_ip"
)

// ContextValues holds extracted context values as byte slices for zero-allocation
type ContextValues struct {
	TraceID   []byte
	SpanID    []byte
	UserID    []byte
	SessionID []byte
	RequestID []byte
}

var contextValuesPool = sync.Pool{
	New: func() interface{} {
		return &ContextValues{}
	},
}

// GetContextValues gets a ContextValues from pool
func GetContextValues() *ContextValues {
	return contextValuesPool.Get().(*ContextValues)
}

// PutContextValues returns ContextValues to pool
func PutContextValues(cv *ContextValues) {
	cv.TraceID = nil
	cv.SpanID = nil
	cv.UserID = nil
	cv.SessionID = nil
	cv.RequestID = nil
	contextValuesPool.Put(cv)
}

// ExtractToBytes extracts context values directly as []byte for zero allocation
func ExtractToBytes(ctx context.Context) *ContextValues {
	cv := GetContextValues()

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		cv.TraceID = StringToBytes(traceID)
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		cv.SpanID = StringToBytes(spanID)
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		cv.UserID = StringToBytes(userID)
	}
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok && sessionID != "" {
		cv.SessionID = StringToBytes(sessionID)
	}
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		cv.RequestID = StringToBytes(requestID)
	}

	return cv
}

// WithTraceID adds trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithSpanID adds span ID to context
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, SpanIDKey, spanID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithSessionID adds session ID to context
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// ExtractFromContext extracts all context values - Optimized version
func ExtractFromContext(ctx context.Context) map[string]string {
	result := GetMapStr()

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		result["trace_id"] = traceID
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		result["span_id"] = spanID
	}
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		result["user_id"] = userID
	}
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok && sessionID != "" {
		result["session_id"] = sessionID
	}
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		result["request_id"] = requestID
	}

	return result
}
