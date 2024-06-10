package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/**
 * @author: gagral.x@gmail.com
 * @time: 2024/6/8 1:21
 * @file: logger_file.go
 * @description: log writer file
 */
const filename string = "log.LOG"

// getFileLogWriter returns the WriteSyncer for logging to a file.
func getFileLogWriter(config *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s", config.LogPath, filename),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
