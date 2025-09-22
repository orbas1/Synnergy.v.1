package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Level represents a logging severity level.
type Level int

const (
	// DebugLevel emits verbose diagnostic information.
	DebugLevel Level = iota
	// InfoLevel captures standard operational events.
	InfoLevel
	// WarnLevel highlights anomalies that should be triaged.
	WarnLevel
	// ErrorLevel records failures that require immediate attention.
	ErrorLevel
)

// Settings configure a logger instance.
type Settings struct {
	Level         Level
	Format        string
	IncludeCaller bool
	Writers       []io.Writer
	StaticFields  map[string]any
}

// Logger emits structured log lines that align with the enterprise CLI and web dashboards.
type Logger struct {
	mu            sync.Mutex
	level         Level
	format        string
	includeCaller bool
	writers       []io.Writer
	staticFields  map[string]any
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// NewLogger creates a Logger from the supplied settings.
func NewLogger(settings Settings) *Logger {
	writers := settings.Writers
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	format := settings.Format
	if format == "" {
		format = "json"
	}
	return &Logger{
		level:         settings.Level,
		format:        strings.ToLower(format),
		includeCaller: settings.IncludeCaller,
		writers:       append([]io.Writer(nil), writers...),
		staticFields:  copyMap(settings.StaticFields),
	}
}

// Configure installs the global logger using the provided settings.
func Configure(settings Settings) {
	once.Do(func() {})
	defaultLogger = NewLogger(settings)
}

// Default returns the process wide logger, initialising it if required.
func Default() *Logger {
	once.Do(func() {
		defaultLogger = NewLogger(Settings{Level: InfoLevel, Format: "json", IncludeCaller: true})
	})
	return defaultLogger
}

// WithFields returns a derived logger with the provided fields attached to every event.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	clone := *l
	clone.staticFields = mergeMaps(l.staticFields, fields)
	return &clone
}

// WithContext returns a derived logger enriched with trace identifiers from the context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return l
	}
	sc := span.SpanContext()
	return l.WithFields(map[string]any{
		"trace_id": sc.TraceID().String(),
		"span_id":  sc.SpanID().String(),
	})
}

// SetLevel updates the minimum severity level for the logger.
func (l *Logger) SetLevel(level Level) { l.level = level }

// Debug logs a debug level message.
func (l *Logger) Debug(msg string, kv ...any) { l.log(DebugLevel, msg, kv...) }

// Info logs an info level message.
func (l *Logger) Info(msg string, kv ...any) { l.log(InfoLevel, msg, kv...) }

// Warn logs a warning level message.
func (l *Logger) Warn(msg string, kv ...any) { l.log(WarnLevel, msg, kv...) }

// Error logs an error level message.
func (l *Logger) Error(msg string, kv ...any) { l.log(ErrorLevel, msg, kv...) }

func (l *Logger) log(level Level, msg string, kv ...any) {
	if l == nil {
		return
	}
	if level < l.level {
		return
	}
	entry := map[string]any{
		"level": levelString(level),
		"msg":   msg,
		"ts":    time.Now().UTC().Format(time.RFC3339Nano),
	}
	for k, v := range l.staticFields {
		entry[k] = v
	}
	if l.includeCaller {
		if pc, file, line, ok := runtime.Caller(2); ok {
			fn := runtime.FuncForPC(pc)
			entry["caller"] = fmt.Sprintf("%s:%d", trimPath(file), line)
			if fn != nil {
				entry["function"] = fn.Name()
			}
		}
	}
	fields := kvToMap(kv)
	for k, v := range fields {
		entry[k] = v
	}
	l.write(entry)
}

func (l *Logger) write(entry map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var payload []byte
	switch l.format {
	case "text":
		payload = []byte(formatText(entry))
	default:
		data, err := json.Marshal(entry)
		if err != nil {
			return
		}
		payload = append(data, '\n')
	}

	for _, w := range l.writers {
		_, _ = w.Write(payload)
	}
}

func levelString(level Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	default:
		return "error"
	}
}

func kvToMap(kv []any) map[string]any {
	out := map[string]any{}
	for i := 0; i < len(kv); i += 2 {
		key := fmt.Sprint(kv[i])
		var value any = "<missing>"
		if i+1 < len(kv) {
			value = kv[i+1]
		}
		out[key] = value
	}
	return out
}

func copyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	cp := make(map[string]any, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

func mergeMaps(base, overlay map[string]any) map[string]any {
	merged := copyMap(base)
	if merged == nil {
		merged = map[string]any{}
	}
	for k, v := range overlay {
		merged[k] = v
	}
	return merged
}

func formatText(entry map[string]any) string {
	timestamp := entry["ts"].(string)
	level := entry["level"].(string)
	msg := entry["msg"].(string)
	delete(entry, "ts")
	delete(entry, "level")
	delete(entry, "msg")
	builder := strings.Builder{}
	builder.WriteString(timestamp)
	builder.WriteString(" ")
	builder.WriteString(level)
	builder.WriteString(" ")
	builder.WriteString(msg)
	for k, v := range entry {
		builder.WriteString(" ")
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprint(v))
	}
	builder.WriteString("\n")
	return builder.String()
}

func trimPath(path string) string {
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}
	return path
}

// Convenience helpers using the default logger -------------------------------------------------

// Debug logs using the default logger.
func Debug(msg string, kv ...any) { Default().Debug(msg, kv...) }

// Info logs at info level using the default logger.
func Info(msg string, kv ...any) { Default().Info(msg, kv...) }

// Warn logs at warn level using the default logger.
func Warn(msg string, kv ...any) { Default().Warn(msg, kv...) }

// Error logs at error level using the default logger.
func Error(msg string, kv ...any) { Default().Error(msg, kv...) }
