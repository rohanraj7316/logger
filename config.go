package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
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

// Configure initializes Istio's logging subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the logging system is ready to accept data.
func Configure(options *Options) error {
	outputLevel, err := options.GetOutputLevel()
	if err != nil {
		// bad format specified
		return err
	}

	stackTraceLevel, err := options.GetStackTraceLevel()
	if err != nil {
		// bad format specified
		return err
	}

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeTime:     formatDate,
	}

	enc := zapcore.NewJSONEncoder(encCfg)

	opts := []zap.Option{
		zap.ErrorOutput(zapcore.AddSync(options.ErrOutput)),
	}

	if options.IncludeCallerSourceLocation {
		opts = append(opts, zap.AddCaller())
	}

	if stackTraceLevel != NoneLevel {
		opts = append(opts, zap.AddStacktrace(levelToZap[stackTraceLevel]))
	}

	l := zap.New(
		zapcore.NewCore(enc, zapcore.AddSync(options.Output), zap.NewAtomicLevelAt(levelToZap[outputLevel])),
		opts...,
	)

	logger = l.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(levelToZap[stackTraceLevel]))
	sugar = logger.Sugar()

	// capture global zap logging and force it through our logger
	_ = zap.ReplaceGlobals(l)

	// capture standard golang "log" package output and force it through our logger
	_ = zap.RedirectStdLog(logger)

	// capture gRPC logging
	if options.LogGrpc {
		grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.WithOptions(zap.AddCallerSkip(2))))
	}

	return nil
}
