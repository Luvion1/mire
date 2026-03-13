package core

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkLogEntryPoolOperations benchmarks the log entry pool operations
func BenchmarkLogEntryPoolOperations(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entry := GetEntry()
		entry.Timestamp = time.Now()
		entry.Level = INFO
		entry.Message = []byte("test message")
		entry.Fields["test"] = []byte(fmt.Sprintf("%d", i))

		PutEntry(entry)
	}
}

// BenchmarkGetEntry benchmarks getting entries from the pool
func BenchmarkGetEntry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		entry := GetEntry()
		// Don't return to pool to measure acquisition performance
		// This tests the hot path of pool acquisition
		PutEntry(entry)
	}
}

// BenchmarkPutEntry benchmarks returning entries to the pool
func BenchmarkPutEntry(b *testing.B) {
	entries := make([]*LogEntry, b.N)
	for i := 0; i < b.N; i++ {
		entries[i] = GetEntry()
		entries[i].Timestamp = time.Now()
		entries[i].Level = DEBUG
		entries[i].Message = []byte("test message")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PutEntry(entries[i])
	}
}

// BenchmarkLogEntryFormatLogToBytes benchmarks the log to bytes formatting
func BenchmarkLogEntryFormatLogToBytes(b *testing.B) {
	entry := GetEntry()
	defer PutEntry(entry)

	entry.Timestamp = time.Now()
	entry.Level = ERROR
	entry.LevelName = []byte("ERROR")
	entry.Message = []byte("test message for benchmark")
	entry.Fields = map[string][]byte{
		"user_id": []byte("123"),
		"action":  []byte("login"),
		"status":  []byte("success"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 1024)
		_ = entry.formatLogToBytes(buf)
	}
}

// BenchmarkLevelBytesConversion benchmarks the conversion of level to bytes
func BenchmarkLevelBytesConversion(b *testing.B) {
	levels := []Level{TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, FATAL, PANIC}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		_ = level.Bytes()
	}
}

// BenchmarkLevelToBytesMethod benchmarks the ToBytes method
func BenchmarkLevelToBytesMethod(b *testing.B) {
	levels := []Level{TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, FATAL, PANIC}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		_ = level.ToBytes()
	}
}

// BenchmarkCallerPoolOperations benchmarks the caller info pool operations
func BenchmarkCallerPoolOperations(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		info := GetCaller()
		info.File = "test.go"
		info.Line = i
		info.Function = "BenchmarkFunction"
		info.Package = "main"

		PutCaller(info)
	}
}

// BenchmarkMapInterfacePoolOperations benchmarks the map interface pool operations
func BenchmarkMapInterfacePoolOperations(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := GetMapByte()
		m["key"] = []byte("value")
		m["num"] = []byte(fmt.Sprintf("%d", i))

		PutMapByte(m)
	}
}

// BenchmarkMapFloatPoolOperations benchmarks the map float pool operations
func BenchmarkMapFloatPoolOperations(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := GetMapFloat()
		m["metric1"] = 1.5
		m["metric2"] = float64(i)

		PutMapFloat(m)
	}
}

// BenchmarkBufferPoolOperations benchmarks the buffer pool operations
func BenchmarkBufferPoolOperations(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := GetBuf()
		*buf = append(*buf, []byte("test data")...)

		PutBuf(buf)
	}
}

// BenchmarkIntToBytesConversion benchmarks converting int to bytes manually
func BenchmarkIntToBytesConversion(b *testing.B) {
	entry := &LogEntry{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 20)
		_ = entry.intToBytes(buf, i)
	}
}

// BenchmarkInt64ToBytesConversion benchmarks converting int64 to bytes manually
func BenchmarkInt64ToBytesConversion(b *testing.B) {
	entry := &LogEntry{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 20)
		_ = entry.int64ToBytes(buf, int64(i))
	}
}

// BenchmarkZeroAllocJSONSerialize benchmarks the zero-allocation JSON serialization
func BenchmarkZeroAllocJSONSerialize(b *testing.B) {
	entry := GetEntry()
	defer PutEntry(entry)

	entry.Timestamp = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry.Level = INFO
	entry.LevelName = []byte("INFO")
	entry.Message = []byte("test message")
	entry.Fields = map[string][]byte{
		"user_id": []byte("123"),
		"action":  []byte("login"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = entry.ZeroAllocJSONSerialize()
	}
}
