package formatter

import (
	"bytes"
	"strconv"

	"github.com/Luvion1/mire/core"
	"github.com/Luvion1/mire/util"
)

// to
var (
	jsonTimestampKey = []byte("\"timestamp\":\"")
	jsonLevelKey     = []byte("\"level_name\":\"")
	jsonMessageKey   = []byte("\"message\":\"")
	jsonPidKey       = []byte(",\"pid\":")
	jsonCallerKey    = []byte(",\"caller\":\"")
	jsonTraceKey     = []byte(",\"trace_id\":\"")
	jsonSpanKey      = []byte(",\"span_id\":\"")
	jsonUserKey      = []byte(",\"user_id\":\"")
	jsonFieldsKey    = []byte(",\"fields\":")
	jsonStackKey     = []byte(",\"stack_trace\":\"")
	jsonQuote        = []byte("\"")
	jsonComma        = []byte(",")
	jsonColon        = []byte(":")
	jsonBraceOpen    = []byte("{")
)

// JSONFormatter formats log entries in JSON format
type JSONFormatter struct {
	PrettyPrint       bool                                     // Enable pretty-printed JSON
	TimestampFormat   string                                   // Custom timestamp format
	ShowCaller        bool                                     // Show caller information
	ShowGoroutine     bool                                     // Show goroutine ID
	ShowPID           bool                                     // Show process ID
	ShowTrace         bool                                     // Show trace information
	IncludeStackTrace bool                                     // Enable stack trace for errors
	ShowDuration      bool                                     // Show operation duration
	FieldKeyMap       map[string]string                        // Map for renaming fields
	DisableHTMLEscape bool                                     // Disable HTML escaping in JSON
	SensitiveFields   []string                                 // List of sensitive field names
	MaskSensitiveData bool                                     // Whether to mask sensitive data
	MaskValue         string                                   // String value to use for masking
	MaskStringBytes   []byte                                   // Byte slice for masking (zero-allocation)
	FieldTransformers map[string]func(interface{}) interface{} // Functions to transform field values
}

// NewJSONFormatter creates a new JSONFormatter
func NewJSON() *JSONFormatter {
	return &JSONFormatter{
		MaskValue:         "[MASKED]",
		MaskStringBytes:   []byte("[MASKED]"),
		FieldKeyMap:       make(map[string]string),
		FieldTransformers: make(map[string]func(interface{}) interface{}),
		SensitiveFields:   make([]string, 0),
	}
}

// Format formats a log entry into JSON
func (f *JSONFormatter) Format(buf *bytes.Buffer, entry *core.LogEntry) error {
	if f.PrettyPrint {
		return f.formatWithStandardEncoder(buf, entry)
	}

	return f.formatManually(buf, entry)
}

// formatManually creates JSON manually
func (f *JSONFormatter) formatManually(buf *bytes.Buffer, entry *core.LogEntry) error {
	buf.Write(jsonBraceOpen)

	// Add timestamp - manually format to avoid allocation
	buf.Write(jsonTimestampKey)
	util.FormatTimestamp(buf, entry.Timestamp, f.TimestampFormat)
	buf.Write(jsonQuote)
	buf.Write(jsonComma)

	// Add level
	buf.Write(jsonLevelKey)
	if len(entry.LevelName) > 0 {
		buf.Write(entry.LevelName)
	} else {
		buf.Write(entry.Level.Bytes())
	}
	buf.Write(jsonQuote)
	buf.Write(jsonComma)

	// Add message
	buf.Write(jsonMessageKey)
	if entry.Message != nil {
		escapeJSON(buf, entry.Message)
	}
	buf.Write(jsonQuote)

	// Add comma after message if more fields will follow
	hasMoreFields := f.ShowPID || (f.ShowCaller && entry.Caller != nil) || len(entry.Fields) > 0 || 
		f.ShowTrace || (f.IncludeStackTrace && len(entry.StackTrace) > 0)
	if hasMoreFields {
		buf.Write(jsonComma)
	}

	// Add PID if needed
	if f.ShowPID {
		buf.Write(jsonPidKey)
		util.WriteInt(buf, int64(entry.PID))
		// Add comma if more fields will follow
		hasMoreAfterPID := (f.ShowCaller && entry.Caller != nil) || len(entry.Fields) > 0 || 
			f.ShowTrace || (f.IncludeStackTrace && len(entry.StackTrace) > 0)
		if hasMoreAfterPID {
			buf.Write(jsonComma)
		}
	}

	// Add caller info if needed
	if f.ShowCaller && entry.Caller != nil {
		buf.Write(jsonCallerKey)
		buf.Write(core.StringToBytes(entry.Caller.File))
		buf.Write(jsonColon)
		util.WriteInt(buf, int64(entry.Caller.Line))
		buf.WriteByte('"')
		// Add comma if more fields will follow
		hasMoreAfterCaller := len(entry.Fields) > 0 || f.ShowTrace || (f.IncludeStackTrace && len(entry.StackTrace) > 0)
		if hasMoreAfterCaller {
			buf.Write(jsonComma)
		}
	}

	// Add fields if present
	if len(entry.Fields) > 0 || len(entry.KeyVals) > 0 {
		buf.Write(jsonFieldsKey)
		f.formatAllFields(buf, entry.Fields, entry.KeyVals)
		// Add comma if trace info or stack trace will follow
		if f.ShowTrace || (f.IncludeStackTrace && len(entry.StackTrace) > 0) {
			buf.Write(jsonComma)
		}
	}

	// Add trace info if needed - organize in a way that reduces branching
	if f.ShowTrace {
		if entry.TraceID != nil {
			buf.Write(jsonTraceKey)
			buf.Write(entry.TraceID)
			buf.Write(jsonQuote)
			// Add comma if more trace fields or stack trace will follow
			hasMoreTrace := entry.SpanID != nil || entry.UserID != nil || (f.IncludeStackTrace && len(entry.StackTrace) > 0)
			if hasMoreTrace {
				buf.Write(jsonComma)
			}
		}
		if entry.SpanID != nil {
			buf.Write(jsonSpanKey)
			buf.Write(entry.SpanID)
			buf.Write(jsonQuote)
			// Add comma if user ID or stack trace will follow
			hasMoreTrace := entry.UserID != nil || (f.IncludeStackTrace && len(entry.StackTrace) > 0)
			if hasMoreTrace {
				buf.Write(jsonComma)
			}
		}
		if entry.UserID != nil {
			buf.Write(jsonUserKey)
			buf.Write(entry.UserID)
			buf.Write(jsonQuote)
			// Add comma if stack trace will follow
			if f.IncludeStackTrace && len(entry.StackTrace) > 0 {
				buf.Write(jsonComma)
			}
		}
	}

	if f.IncludeStackTrace && len(entry.StackTrace) > 0 {
		buf.Write(jsonStackKey)
		escapeJSON(buf, entry.StackTrace)
		buf.WriteByte('"')
	}

	buf.WriteByte('}')
	buf.WriteByte('\n')

	return nil
}

// formatWithStandardEncoder handles pretty printing
func (f *JSONFormatter) formatWithStandardEncoder(buf *bytes.Buffer, entry *core.LogEntry) error {
	if f.PrettyPrint {
		return f.formatManuallyWithIndent(buf, entry)
	}
	return f.formatManually(buf, entry)
}

// formatManuallyWithIndent formats JSON with indentation
func (f *JSONFormatter) formatManuallyWithIndent(buf *bytes.Buffer, entry *core.LogEntry) error {
	indentBuf := util.GetBuf()
	defer util.PutBuf(indentBuf)

	for i := 0; i < 10; i++ {
		indentBuf.WriteString("  ")
	}
	indentLevels := indentBuf.Bytes()

	indent := func(level int) {
		if level <= 0 {
			return
		}
		// Use pre-allocated indentation
		indentSize := level * 2
		if indentSize <= len(indentLevels) {
			buf.Write(indentLevels[:indentSize])
		} else {
			// If we need more indentation than pre-allocated, add more
			for i := 0; i < level; i++ {
				buf.WriteString("  ")
			}
		}
	}

	newline := func(level int) {
		buf.WriteByte('\n')
		indent(level)
	}

	// Start JSON object
	buf.WriteByte('{')

	// Add timestamp
	newline(1)
	buf.WriteString("\"timestamp\": \"")
	util.FormatTimestamp(buf, entry.Timestamp, f.TimestampFormat)
	buf.WriteString("\"")

	// Add level
	buf.WriteString(",\n  ")
	indent(1)
	buf.WriteString("\"level_name\": \"")
	buf.Write(entry.Level.Bytes()) // Using pre-allocated level bytes
	buf.WriteString("\"")

	// Add message
	buf.WriteString(",\n  ")
	indent(1)
	buf.WriteString("\"message\": \"")
	escapeJSON(buf, entry.Message)
	buf.WriteString("\"")

	// Add PID if needed
	if f.ShowPID && entry.PID != 0 {
		buf.WriteString(",\n  ")
		indent(1)
		buf.WriteString("\"pid\": ")
		util.WriteInt(buf, int64(entry.PID))
	}

	// Add caller info if needed
	if f.ShowCaller && entry.Caller != nil {
		buf.WriteString(",\n  ")
		indent(1)
		buf.WriteString("\"caller\": \"")
		buf.Write(core.StringToBytes(entry.Caller.File))
		buf.WriteByte(':')
		util.WriteInt(buf, int64(entry.Caller.Line))
		buf.WriteByte('"')
	}

	// Add fields if present
	if len(entry.Fields) > 0 || len(entry.KeyVals) > 0 {
		buf.WriteString(",\n  ")
		indent(1)
		buf.WriteString("\"fields\": ")
		// For indented fields, we need to format them manually with indentation
		f.formatAllFieldsIndented(buf, entry.Fields, entry.KeyVals, 2)
	}

	// Add trace info if needed
	if f.ShowTrace {
		if entry.TraceID != nil {
			buf.WriteString(",\n  ")
			indent(1)
			buf.WriteString("\"trace_id\": \"")
			buf.Write(entry.TraceID)
			buf.WriteByte('"')
		}
		if entry.SpanID != nil {
			buf.WriteString(",\n  ")
			indent(1)
			buf.WriteString("\"span_id\": \"")
			buf.Write(entry.SpanID)
			buf.WriteByte('"')
		}
		if entry.UserID != nil {
			buf.WriteString(",\n  ")
			indent(1)
			buf.WriteString("\"user_id\": \"")
			buf.Write(entry.UserID)
			buf.WriteByte('"')
		}
	}

	if f.IncludeStackTrace && len(entry.StackTrace) > 0 {
		buf.WriteString(",\n  ")
		indent(1)
		buf.WriteString("\"stack_trace\": \"")
		escapeJSON(buf, entry.StackTrace)
		buf.WriteByte('"')
	}

	newline(0)
	buf.WriteByte('}')
	buf.WriteByte('\n')

	return nil
}

// escapeJSON escapes special characters in JSON strings - optimized for zero allocation
func escapeJSON(buf *bytes.Buffer, data []byte) {
	if len(data) == 0 {
		return
	}

	start := 0
	for i, b := range data {
		var escape []byte
		switch b {
		case '"':
			escape = []byte("\\\"")
		case '\\':
			escape = []byte("\\\\")
		case '\b':
			escape = []byte("\\b")
		case '\f':
			escape = []byte("\\f")
		case '\n':
			escape = []byte("\\n")
		case '\r':
			escape = []byte("\\r")
		case '\t':
			escape = []byte("\\t")
		default:
			if b < 0x20 {
				buf.Write(data[start:i])
				buf.Write([]byte("\\u00"))
				hexChars := "0123456789abcdef"
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0xF])
				start = i + 1
				continue
			}
			continue
		}
		if escape != nil {
			buf.Write(data[start:i])
			buf.Write(escape)
			start = i + 1
		}
	}
	buf.Write(data[start:])
}

// formatJSONValue formats a value for JSON output
func (f *JSONFormatter) formatJSONValue(buf *bytes.Buffer, v interface{}) {
	switch val := v.(type) {
	case string:
		buf.WriteByte('"')
		escapeJSON(buf, core.StringToBytes(val))
		buf.WriteByte('"')
	case []byte:
		buf.WriteByte('"')
		escapeJSON(buf, val)
		buf.WriteByte('"')
	case int:
		tempBuf := util.GetSmallBuf()
		numBytes := strconv.AppendInt(tempBuf[:0], int64(val), 10)
		buf.Write(numBytes)
		util.PutSmallBuf(tempBuf)
	case int64:
		tempBuf := util.GetSmallBuf()
		numBytes := strconv.AppendInt(tempBuf[:0], val, 10)
		buf.Write(numBytes)
		util.PutSmallBuf(tempBuf)
	case float64:
		tempBuf := util.GetSmallBuf()
		numBytes := strconv.AppendFloat(tempBuf[:0], val, 'g', -1, 64)
		buf.Write(numBytes)
		util.PutSmallBuf(tempBuf)
	case bool:
		if val {
			buf.Write([]byte("true"))
		} else {
			buf.Write([]byte("false"))
		}
	case nil:
		buf.Write([]byte("null"))
	default:
		transformed := f.transformValue(val, "<complex-type>")
		buf.WriteByte('"')
		escapeJSON(buf, core.StringToBytes(transformed))
		buf.WriteByte('"')
	}
}

// isSensitiveField checks if a field is in the sensitive fields list
func (f *JSONFormatter) isSensitiveField(field string) bool {
	for _, sensitiveField := range f.SensitiveFields {
		if field == sensitiveField {
			return true
		}
	}
	return false
}

// isSensitive checks if a field name is in the sensitive fields list
func (f *JSONFormatter) isSensitive(field string) bool {
	return f.isSensitiveField(field)
}

// createMapForSensitiveCheck creates a map for O(1) sensitive field lookup when there are many sensitive fields
func (f *JSONFormatter) createSensitiveFieldMap() map[string]bool {
	if len(f.SensitiveFields) == 0 {
		return nil
	}

	// Only create map if there are enough fields to justify it
	if len(f.SensitiveFields) < 5 {
		return nil
	}

	fieldMap := make(map[string]bool, len(f.SensitiveFields))
	for _, field := range f.SensitiveFields {
		fieldMap[field] = true
	}
	return fieldMap
}

// transformValue converts value to string
func (f *JSONFormatter) transformValue(val interface{}, defaultVal string) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		return defaultVal
	}
}

// formatAllFields formats fields map and key-value pairs in JSON format
func (f *JSONFormatter) formatAllFields(buf *bytes.Buffer, fields map[string][]byte, keyvals [][]byte) {
	buf.WriteByte('{')
	first := true

	// Format map fields
	for k, v := range fields {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		buf.WriteByte('"')
		buf.Write(core.StringToBytes(k))
		buf.Write([]byte("\":\""))
		if f.MaskSensitiveData && f.isSensitiveField(k) {
			buf.Write(f.MaskStringBytes)
		} else {
			escapeJSON(buf, v)
		}
		buf.WriteByte('"')
	}

	// Format key-value pairs
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 >= len(keyvals) {
			break
		}
		if !first {
			buf.WriteByte(',')
		}
		first = false

		k := core.BytesToString(keyvals[i])
		v := keyvals[i+1]

		buf.WriteByte('"')
		buf.Write(keyvals[i])
		buf.Write([]byte("\":\""))
		if f.MaskSensitiveData && f.isSensitiveField(k) {
			buf.Write(f.MaskStringBytes)
		} else {
			escapeJSON(buf, v)
		}
		buf.WriteByte('"')
	}

	buf.WriteByte('}')
}

// formatFields formats fields map in JSON format
func (f *JSONFormatter) formatFields(buf *bytes.Buffer, fields map[string][]byte) {
	f.formatAllFields(buf, fields, nil)
}

// formatAllFieldsIndented formats a fields map and key-value pairs in JSON format with indentation
func (f *JSONFormatter) formatAllFieldsIndented(buf *bytes.Buffer, fields map[string][]byte, keyvals [][]byte, indentLevel int) {
	indentBuf := util.GetBuf()
	defer util.PutBuf(indentBuf)

	for i := 0; i < indentLevel; i++ {
		indentBuf.WriteString("  ")
	}
	indentBytes := indentBuf.Bytes()
	originalIndent := make([]byte, len(indentBytes))
	copy(originalIndent, indentBytes)

	buf.WriteByte('{')
	
	innerIndent := append(indentBytes, ' ', ' ')
	first := true

	// Helper to write a field
	writeField := func(k string, kBytes []byte, v []byte) {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.WriteByte('\n')
		buf.Write(innerIndent)
		buf.WriteByte('"')
		buf.Write(kBytes)
		buf.Write([]byte("\": \""))
		if f.MaskSensitiveData && f.isSensitiveField(k) {
			buf.Write(f.MaskStringBytes)
		} else {
			escapeJSON(buf, v)
		}
		buf.WriteByte('"')
	}

	// Map fields
	for k, v := range fields {
		writeField(k, core.StringToBytes(k), v)
	}

	// Key-value pairs
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 >= len(keyvals) {
			break
		}
		kBytes := keyvals[i]
		k := core.BytesToString(kBytes)
		v := keyvals[i+1]
		writeField(k, kBytes, v)
	}

	if !first {
		buf.WriteByte('\n')
		buf.Write(originalIndent)
	}
	buf.WriteByte('}')
}

// formatFieldsIndented formats a fields map in JSON format with indentation
func (f *JSONFormatter) formatFieldsIndented(buf *bytes.Buffer, fields map[string][]byte, indentLevel int) {
	f.formatAllFieldsIndented(buf, fields, nil, indentLevel)
}
