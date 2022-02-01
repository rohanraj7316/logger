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

// Options defines the set of options supported by Istio's component logging package.
type Options struct {

	// JSONEncoding controls whether the log is formatted as JSON.
	JSONEncoding bool

	// IncludeCallerSourceLocation determines whether log messages include the source location of the caller.
	IncludeCallerSourceLocation bool

	// LogGrpc indicates that Grpc logs should be captured. The default is true.
	// This is not exposed through the command-line flags, as this flag is mainly useful for testing: Grpc
	// stack will hold on to the logger even though it gets closed. This causes data races.
	LogGrpc bool

	stackTraceLevel string
	outputLevel     string
}

// NewOptions returns a new set of options, initialized to the defaults
func NewOptions() *Options {
	return &Options{
		outputLevel:     string(defaultOutputLevel),
		stackTraceLevel: string(defaultStackTraceLevel),
		LogGrpc:         true,
	}
}

func isValid(level Level) bool {
	switch level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, NoneLevel:
		return true
	default:
		return false
	}
}

// SetOutputLevel sets the minimum log output level.
func (o *Options) SetOutputLevel(level Level) error {
	if ok := isValid(level); !ok {
		return fmt.Errorf("unknown output level: %v", level)
	}
	o.outputLevel = string(level)
	return nil
}

// GetOutputLevel returns the minimum log output level.
func (o *Options) GetOutputLevel() (Level, error) {
	l := Level(o.outputLevel)
	if ok := isValid(l); !ok {
		return "", fmt.Errorf("unknown output level: %v", l)
	}
	return l, nil
}

// SetStackTraceLevel sets the minimum stack trace capture level.
func (o *Options) SetStackTraceLevel(level Level) error {
	if ok := isValid(level); !ok {
		return fmt.Errorf("unknown stack trace level: %v", level)
	}
	o.stackTraceLevel = string(level)
	return nil
}

// GetStackTraceLevel returns the current stack trace level.
func (o *Options) GetStackTraceLevel() (Level, error) {
	l := Level(o.stackTraceLevel)
	if ok := isValid(l); !ok {
		return "", fmt.Errorf("unknown stack trace level: %v", l)
	}
	return l, nil
}
