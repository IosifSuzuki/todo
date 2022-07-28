package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var zapLog *zap.Logger

func init() {
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	core := zapcore.NewTee(
		zapcore.NewCore(getLogFileEncoder(), getLogFileWriter(), zapcore.DebugLevel),
		zapcore.NewCore(getConsoleEncoder(), getConsoleWriter(), zapcore.DebugLevel),
	)
	zapLog = zap.New(core)
}

func getLogFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogFileWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./logging.log")
	return zapcore.AddSync(file)
}

func getConsoleWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}
