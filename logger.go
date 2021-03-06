package logger

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

var logger = zap.NewNop()

var sugar = logger.Sugar()

type Field struct {
	Key   string
	Value interface{}
}

func parse(fields ...Field) []zapcore.Field {
	pFields := []zapcore.Field{}
	for i := 0; i < len(fields); i++ {
		key := fields[i].Key
		value := fields[i].Value

		switch t := value.(type) {
		case string:
			pFields = append(pFields, zapcore.Field{
				Key:    key,
				String: t,
				Type:   zapcore.StringType,
			})
		case []byte:
			pFields = append(pFields, zapcore.Field{
				Key:    key,
				String: string(t),
				Type:   zapcore.StringType,
			})
		case interface{}:
			b, err := json.Marshal(t)
			if err != nil {
				eStruct := map[string]interface{}{
					"err": err.Error(),
					"msg": "unable to parse",
				}

				eB, _ := json.Marshal(eStruct)

				pFields = append(pFields, zapcore.Field{
					Key:    key,
					String: string(eB),
					Type:   zapcore.StringType,
				})
				continue
			}
			pFields = append(pFields, zapcore.Field{
				Key:    key,
				String: string(b),
				Type:   zapcore.StringType,
			})
		}
	}

	return pFields
}

func formatDate(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	micros := t.Nanosecond() / 1000

	buf := make([]byte, 27)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((micros/100000)%10) + '0'
	buf[21] = byte((micros/10000)%10) + '0'
	buf[22] = byte((micros/1000)%10) + '0'
	buf[23] = byte((micros/100)%10) + '0'
	buf[24] = byte((micros/10)%10) + '0'
	buf[25] = byte((micros)%10) + '0'
	buf[26] = 'Z'

	enc.AppendString(string(buf))
}

// Configure initializes Istio's logging subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the logging system is ready to accept data.
func Configure(config ...Config) error {
	cfg := configDefault(config...)

	outputLevel, err := cfg.GetOutputLevel()
	if err != nil {
		// bad format specified
		return err
	}

	stackTraceLevel, err := cfg.GetStackTraceLevel()
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
		zap.ErrorOutput(zapcore.AddSync(cfg.ErrOutput)),
	}

	if cfg.IncludeCallerSourceLocation {
		opts = append(opts, zap.AddCaller())
	}

	if stackTraceLevel != NoneLevel {
		opts = append(opts, zap.AddStacktrace(levelToZap[stackTraceLevel]))
	}

	l := zap.New(
		zapcore.NewCore(enc, zapcore.AddSync(cfg.Output), zap.NewAtomicLevelAt(levelToZap[outputLevel])),
		opts...,
	)

	logger = l.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(levelToZap[stackTraceLevel]))
	sugar = logger.Sugar()

	// capture global zap logging and force it through our logger
	_ = zap.ReplaceGlobals(l)

	// capture standard golang "log" package output and force it through our logger
	_ = zap.RedirectStdLog(logger)

	// capture gRPC logging
	if cfg.LogGrpc {
		grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.WithOptions(zap.AddCallerSkip(2))))
	}

	return nil
}

// Debug outputs a message at debug level.
// This call is a wrapper around [Logger.Debug](https://godoc.org/go.uber.org/zap#Logger.Debug)
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Error outputs a message at error level.
// This call is a wrapper around [logger.Error](https://godoc.org/go.uber.org/zap#logger.Error)
func Error(msg string, fields ...Field) {
	logger.Error(msg, parse(fields...)...)
}

// Warn outputs a message at warn level.
// This call is a wrapper around [logger.Warn](https://godoc.org/go.uber.org/zap#logger.Warn)
func Warn(msg string, fields ...Field) {
	logger.Warn(msg, parse(fields...)...)
}

// Info outputs a message at information level.
// This call is a wrapper around [logger.Info](https://godoc.org/go.uber.org/zap#logger.Info)
func Info(msg string, fields ...Field) {
	logger.Info(msg, parse(fields...)...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
// This call is a wrapper around [logger.With](https://godoc.org/go.uber.org/zap#logger.With)
func With(fields ...zapcore.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sync flushes any buffered log entries.
// Processes should normally take care to call Sync before exiting.
// This call is a wrapper around [logger.Sync](https://godoc.org/go.uber.org/zap#logger.Sync)
func Sync() error {
	return logger.Sync()
}
