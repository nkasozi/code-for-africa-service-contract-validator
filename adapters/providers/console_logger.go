package providers

import (
	"fmt"
	"time"
)

type ConsoleLogger struct {
	debug_enabled bool
}

func NewConsoleLogger(debug_enabled bool) *ConsoleLogger {
	return &ConsoleLogger{
		debug_enabled: debug_enabled,
	}
}

func (l *ConsoleLogger) LogInfo(message string) {
	if !l.debug_enabled {
		return
	}
	l.logWithLevel("INFO", message)
}

func (l *ConsoleLogger) LogWarning(message string) {
	if !l.debug_enabled {
		return
	}
	l.logWithLevel("WARN", message)
}

func (l *ConsoleLogger) LogError(message string) {
	l.logWithLevel("ERROR", message)
}

func (l *ConsoleLogger) logWithLevel(level string, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", timestamp, level, message)
}
