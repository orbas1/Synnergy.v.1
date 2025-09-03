package log

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// entry represents a single log line.
type entry map[string]interface{}

var mu sync.Mutex

// log writes a structured JSON entry with level and message.
func log(level, msg string, kv []interface{}) {
	e := entry{
		"level": level,
		"msg":   msg,
		"ts":    time.Now().UTC().Format(time.RFC3339Nano),
	}
	for i := 0; i+1 < len(kv); i += 2 {
		key := fmt.Sprint(kv[i])
		e[key] = kv[i+1]
	}
	b, err := json.Marshal(e)
	if err != nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	os.Stdout.Write(append(b, '\n'))
}

// Info logs at info level.
func Info(msg string, kv ...interface{}) { log("info", msg, kv) }

// Error logs at error level.
func Error(msg string, kv ...interface{}) { log("error", msg, kv) }
