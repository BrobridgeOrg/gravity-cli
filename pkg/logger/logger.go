package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCustomEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func GetLogger() *zap.Logger {
	debugLevel := setupLevel()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewCustomEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		debugLevel,
	)

	var logger *zap.Logger
	if os.Getenv("GRAVITY_DEBUG_MODE") == "debug" {
		//		logger.Info(fmt.Sprintf("Debug mode is set to \"%s\"\n", debugLevel.String()))
		logger = zap.New(core, zap.AddCaller(), zap.Development())
	} else {
		logger = zap.New(core)
	}

	zap.ReplaceGlobals(logger)

	//	logger.Info(fmt.Sprintf("Debug level is set to \"%s\"\n", debugLevel.String()))

	return logger
}

func setupLevel() zap.AtomicLevel {

	debugLevel := zap.DebugLevel
	switch os.Getenv("GRAVITY_DEBUG_LEVEL") {
	case zap.InfoLevel.String():
		debugLevel = zap.InfoLevel
	case zap.WarnLevel.String():
		debugLevel = zap.WarnLevel
	case zap.ErrorLevel.String():
		debugLevel = zap.ErrorLevel
	case zap.DPanicLevel.String():
		debugLevel = zap.DPanicLevel
	case zap.PanicLevel.String():
		debugLevel = zap.PanicLevel
	case zap.FatalLevel.String():
		debugLevel = zap.FatalLevel
	}

	return zap.NewAtomicLevelAt(debugLevel)
}
