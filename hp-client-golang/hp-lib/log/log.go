package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	output       io.Writer = os.Stdout
	mu           sync.Mutex
	minLevel     Level = LevelInfo
	showTime     bool  = true
	showLevel    bool  = true
	colorEnabled bool  = true
)

const (
	colorReset   = "\x1b[0m"
	colorRed     = "\x1b[31m"
	colorGreen   = "\x1b[32m"
	colorYellow  = "\x1b[33m"
	colorBlue    = "\x1b[34m"
	colorMagenta = "\x1b[35m"
	colorCyan    = "\x1b[36m"
	colorGray    = "\x1b[90m"
)

func levelTag(l Level) (string, string) {
	switch l {
	case LevelDebug:
		return "DEBUG", colorGray
	case LevelInfo:
		return "INFO ", colorGreen
	case LevelWarn:
		return "WARN ", colorYellow
	case LevelError:
		return "ERROR", colorRed
	}
	return "UNKWN", colorReset
}

func formatLine(l Level, msg string) string {
	tag, color := levelTag(l)
	parts := make([]string, 0, 3)

	if showTime {
		ts := time.Now().Format("2006-01-02 15:04:05.000")
		if colorEnabled {
			parts = append(parts, colorGray+ts+colorReset)
		} else {
			parts = append(parts, ts)
		}
	}

	if showLevel {
		if colorEnabled {
			parts = append(parts, color+"["+tag+"]"+colorReset)
		} else {
			parts = append(parts, "["+tag+"]")
		}
	}

	if colorEnabled && l == LevelError {
		parts = append(parts, colorRed+msg+colorReset)
	} else if colorEnabled && l == LevelWarn {
		parts = append(parts, colorYellow+msg+colorReset)
	} else {
		parts = append(parts, msg)
	}

	return joinParts(parts)
}

func joinParts(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += " "
		}
		out += p
	}
	return out
}

func writeLine(l Level, msg string) {
	if l < minLevel {
		return
	}
	line := formatLine(l, msg)
	mu.Lock()
	defer mu.Unlock()
	_, _ = fmt.Fprintln(output, line)
}

func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	output = w
}

func SetLevel(l Level) {
	mu.Lock()
	defer mu.Unlock()
	minLevel = l
}

func DisableColor() {
	mu.Lock()
	defer mu.Unlock()
	colorEnabled = false
}

func EnableColor() {
	mu.Lock()
	defer mu.Unlock()
	colorEnabled = true
}

func Error(v ...interface{}) {
	writeLine(LevelError, fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	writeLine(LevelInfo, fmt.Sprint(v...))
}

func Warn(v ...interface{}) {
	writeLine(LevelWarn, fmt.Sprint(v...))
}

func Debug(v ...interface{}) {
	writeLine(LevelDebug, fmt.Sprint(v...))
}

func Errorf(format string, a ...interface{}) {
	writeLine(LevelError, fmt.Sprintf(format, a...))
}

func Infof(format string, a ...interface{}) {
	writeLine(LevelInfo, fmt.Sprintf(format, a...))
}

func Warnf(format string, a ...interface{}) {
	writeLine(LevelWarn, fmt.Sprintf(format, a...))
}

func Debugf(format string, a ...interface{}) {
	writeLine(LevelDebug, fmt.Sprintf(format, a...))
}

func Fatal(v ...interface{}) {
	writeLine(LevelError, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	writeLine(LevelError, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Printf(format string, a ...interface{}) {
	writeLine(LevelInfo, fmt.Sprintf(format, a...))
}

func Println(v ...interface{}) {
	writeLine(LevelInfo, fmt.Sprint(v...))
}
