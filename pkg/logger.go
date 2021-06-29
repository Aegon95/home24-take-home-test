package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// InitLogger creates a new zap logger instance and configures it
func InitLogger() *zap.SugaredLogger {

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.ErrorLevel),
	)
	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	return sugarLogger

}
