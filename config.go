package logger

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
)

// none is used to disable logging output as well as to disable stack tracing.
const none zapcore.Level = 100

var levelToZap = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	NoneLevel:  none,
}

type EncoderConfig struct {
	// Path is array consists of all the json path.
	Path []string

	// EncoderFunc function which helps in encoding the string.
	EncoderFunc func(string) string
}

// Config defines the set of cfg supported by Istio's component logging package.
type Config struct {

	// JSONEncoding controls whether the log is formatted as JSON.
	JSONEncoding bool

	// IncludeCallerSourceLocation determines whether log messages include the source location of the caller.
	IncludeCallerSourceLocation bool

	// LogGrpc indicates that Grpc logs should be captured. The default is true.
	// This is not exposed through the command-line flags, as this flag is mainly useful for testing: Grpc
	// stack will hold on to the logger even though it gets closed. This causes data races.
	LogGrpc bool

	// Output is a writer where logs are written
	//
	// Default: os.Stdout
	Output io.Writer

	// ErrOutput is a writer where logs are written
	//
	// Default: os.Stderr
	ErrOutput io.Writer

	// IsEncoding enables encoding in logging
	//
	// Default: false
	IsEncoding bool

	// EncoderConfig has encoder config which helps in encoding
	EncoderConfig EncoderConfig

	stackTraceLevel string
	outputLevel     string
}

// ConfigDefault it is the set of default configs
var ConfigDefault = Config{
	outputLevel:     string(defaultOutputLevel),
	stackTraceLevel: string(defaultStackTraceLevel),
	Output:          os.Stdout,
	ErrOutput:       os.Stderr,
	LogGrpc:         true,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
