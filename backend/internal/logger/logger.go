package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(level string, logFile string) {
	// Создаём директорию логов если её нет
	if logFile != "" {
		dir := filepath.Dir(logFile) // FIX: корректно получить каталог
		if dir != "" && dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(encCfg)

	// Определяем уровень логирования
	zapLevel := parseLevel(level)

	var cores []zapcore.Core

	// Console output
	cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLevel))

	// File output
	if logFile != "" {
		if f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zapLevel))
		}
	}

	core := zapcore.NewTee(cores...)
	z := zap.New(core, zap.AddCaller())
	Log = z.Sugar()
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
