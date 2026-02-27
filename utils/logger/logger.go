package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	debugEnabled = func() bool {
		// support multiple env vars for flexibility (amk it doesnt even work)
		if v := strings.ToLower(os.Getenv("LOG_DEBUG")); v == "1" || v == "true" {
			return true
		}
		if v := strings.ToLower(os.Getenv("DEBUG")); v == "1" || v == "true" {
			return true
		}
		if v := strings.ToLower(os.Getenv("LOG_LEVEL")); strings.Contains(v, "debug") {
			return true
		}
		return false
	}()
)

type logger struct {
	std *log.Logger
}

var Logger = &logger{std: log.New(os.Stdout, "", 0)}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func color(code string, s string) string {
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", code, s)
}

func (l *logger) Info(format string, args ...interface{}) {
	prefix := color("36", "[INFO]")
	msg := fmt.Sprintf(format, args...)
	l.std.Printf("%s %s %s\n", color("90", timestamp()), prefix, msg)
}

func (l *logger) Warn(format string, args ...interface{}) {
	prefix := color("33", "[WARN]")
	msg := fmt.Sprintf(format, args...)
	l.std.Printf("%s %s %s\n", color("90", timestamp()), prefix, msg)
}

func (l *logger) Error(format string, args ...interface{}) {
	prefix := color("31", "[ERROR]")
	msg := fmt.Sprintf(format, args...)
	l.std.Printf("%s %s %s\n", color("90", timestamp()), prefix, msg)
}

func (l *logger) Debug(format string, args ...interface{}) {
	if !debugEnabled {
		return
	}
	prefix := color("35", "[DEBUG]")
	msg := fmt.Sprintf(format, args...)
	l.std.Printf("%s %s %s\n", color("90", timestamp()), prefix, msg)
}

func SetDebug(enabled bool) { debugEnabled = enabled }

func Info(format string, args ...interface{})  { Logger.Info(format, args...) }
func Warn(format string, args ...interface{})  { Logger.Warn(format, args...) }
func Error(format string, args ...interface{}) { Logger.Error(format, args...) }
func Debug(format string, args ...interface{}) { Logger.Debug(format, args...) }
