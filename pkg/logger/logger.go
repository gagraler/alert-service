package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	once        sync.Once
)

// LogLevel defines the severity of a logger internal.
type LogLevel int8

const (
	// DebugLevel logs debug messages.
	DebugLevel LogLevel = iota - 1
	// InfoLevel logs informational messages.
	InfoLevel
	// WarnLevel logs warning messages.
	WarnLevel
	// ErrorLevel logs error messages.
	ErrorLevel
	FatalLevel
)

// Config holds logger configuration options.
type Config struct {
	LogPath      string
	MaxSize      int
	MaxBackups   int
	MaxAge       int
	LogLevel     LogLevel
	Output       string
	Mode         string
	KafkaBrokers string
	KafkaTopic   string
}

// init initializes the logger when the package is imported.
func init() {
	once.Do(func() {
		var err error
		sugarLogger, logger, err = NewLogger()
		if err != nil {
			fmt.Printf("Could not initialize logger: %v\n", err)
			os.Exit(1)
		}
	})
}

// NewLogger initializes the logger and returns a sugared logger.
func NewLogger() (*zap.SugaredLogger, *zap.Logger, error) {
	var (
		writeSyncer zapcore.WriteSyncer
		encoder     zapcore.Encoder
		core        zapcore.Core
		logConfig   *Config
	)

	// Load configuration from environment variables or use default values
	logConfig = loggerConfigParse()

	encoder = getEncoder(logConfig.Mode)

	switch logConfig.Output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "file":
		writeSyncer = getFileLogWriter(logConfig)
	case "kafka":
		kafkaSyncer, err := getKafkaLogWriter(logConfig)
		if err != nil {
			return nil, nil, err
		}
		writeSyncer = kafkaSyncer
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	switch logConfig.LogLevel {
	case DebugLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	case InfoLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	case WarnLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.WarnLevel)
	case ErrorLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	case FatalLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.FatalLevel)
	default:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	}

	logger = zap.New(core, zap.AddCaller())
	sugaredLogger := logger.Sugar()
	sugaredLogger.Infof("log output is %s", logConfig.Output)
	sugaredLogger.Infof("log mode is %s", strings.ToTitle(logConfig.Mode))
	return sugaredLogger, logger, nil
}

// Logger returns the initialized zap logger.
func Logger() *zap.Logger {
	return logger
}

// SugaredLogger returns the initialized sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return sugarLogger
}

// getEncoder returns the appropriate encoder based on the mode.
func getEncoder(mode string) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig

	if mode == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = customTimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写编码器
	encoderConfig.EncodeTime = customTimeEncoder            // ISO8601 UTC 时间格式
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 相对路径编码器
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// customTimeEncoder formats the time as 2024-06-08 00:51:55.
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// loggerConfigParse loads configuration from environment variables or uses default values.
func loggerConfigParse() *Config {
	logPath := getEnv("LOG_PATH", "./logs")
	logMaxSizeStr := getEnv("LOG_MAX_SIZE", "100")
	logMaxBackupsStr := getEnv("LOG_MAX_BACKUPS", "30")
	logMaxAgeStr := getEnv("LOG_MAX_AGE", "7")
	// debug, info, warn, error, fatal
	logLevelStr := getEnv("LOG_LEVEL", "info")
	// output, file, kafka
	output := getEnv("LOG_OUTPUT", "stdout")
	// dev, prod
	mode := getEnv("LOG_MODE", "prod")
	kafkaBrokers := getEnv("KAFKA_BROKERS", "")
	kafkaTopic := getEnv("KAFKA_TOPIC", "")

	logMaxSize, err := strconv.Atoi(logMaxSizeStr)
	if err != nil {
		logMaxSize = 100
	}
	logMaxBackups, err := strconv.Atoi(logMaxBackupsStr)
	if err != nil {
		logMaxBackups = 30
	}
	logMaxAge, err := strconv.Atoi(logMaxAgeStr)
	if err != nil {
		logMaxAge = 7
	}

	var logLevel LogLevel
	switch strings.ToLower(logLevelStr) {
	case "info":
		logLevel = InfoLevel
	case "debug":
		logLevel = DebugLevel
	case "error":
		logLevel = ErrorLevel
	case "warn":
		logLevel = WarnLevel
	default:
		logLevel = InfoLevel
	}

	return &Config{
		LogPath:      logPath,
		MaxSize:      logMaxSize,
		MaxBackups:   logMaxBackups,
		MaxAge:       logMaxAge,
		LogLevel:     logLevel,
		Output:       output,
		Mode:         mode,
		KafkaBrokers: kafkaBrokers,
		KafkaTopic:   kafkaTopic,
	}
}

// getEnv retrieves the value of the environment variable named by the key or returns the default value if the variable is not present.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
