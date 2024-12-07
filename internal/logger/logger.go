package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {

	loggerConfig := zap.NewProductionConfig()
	// 创建一个新的生产环境编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// 配置时间的编码方式为 ISO8601 格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000000")
	// 配置日志级别的编码方式为大写形式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	loggerConfig.EncoderConfig = encoderConfig

	// 创建 logger
	logger, _ = loggerConfig.Build(zap.AddCaller())
}

func Info(msg string, field ...zap.Field) {
	//defer logger.Sync()
	logger.Info(msg, field...)
}

func Error(msg string, field ...zap.Field) {
	//defer logger.Sync()
	logger.Error(msg, field...)
}

func Debug(msg string, field ...zap.Field) {
	//defer logger.Sync()
	logger.Debug(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	//defer logger.Sync()
	logger.Warn(msg, field...)
}
