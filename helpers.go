package logger

import (
	"fmt"
)

const (
	defaultOutputLevel        = InfoLevel
	defaultStackTraceLevel    = NoneLevel
	defaultOutputPath         = "stdout"
	defaultErrorOutputPath    = "stderr"
	defaultRotationMaxAge     = 30
	defaultRotationMaxSize    = 100 * 1024 * 1024
	defaultRotationMaxBackups = 1000
)

// Level is an enumeration of all supported log levels.
type Level string

const (
	// DebugLevel enables debug level logging
	DebugLevel Level = "debug"
	// InfoLevel enables info level logging
	InfoLevel Level = "info"
	// WarnLevel enables warn level logging
	WarnLevel Level = "warn"
	// ErrorLevel enables error level logging
	ErrorLevel Level = "error"
	// NoneLevel disables logging
	NoneLevel Level = "none"
)

func isValid(level Level) bool {
	switch level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, NoneLevel:
		return true
	default:
		return false
	}
}

// SetOutputLevel sets the minimum log output level.
func (o *Config) SetOutputLevel(level Level) error {
	if ok := isValid(level); !ok {
		return fmt.Errorf("unknown output level: %v", level)
	}
	o.outputLevel = string(level)
	return nil
}

// GetOutputLevel returns the minimum log output level.
func (o *Config) GetOutputLevel() (Level, error) {
	l := Level(o.outputLevel)
	if ok := isValid(l); !ok {
		return "", fmt.Errorf("unknown output level: %v", l)
	}
	return l, nil
}

// SetStackTraceLevel sets the minimum stack trace capture level.
func (o *Config) SetStackTraceLevel(level Level) error {
	if ok := isValid(level); !ok {
		return fmt.Errorf("unknown stack trace level: %v", level)
	}
	o.stackTraceLevel = string(level)
	return nil
}

// GetStackTraceLevel returns the current stack trace level.
func (o *Config) GetStackTraceLevel() (Level, error) {
	l := Level(o.stackTraceLevel)
	if ok := isValid(l); !ok {
		return "", fmt.Errorf("unknown stack trace level: %v", l)
	}
	return l, nil
}
