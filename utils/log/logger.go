package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func EnsureBaseLogDirectoryExist(cfg *LoggerConfig) error {
	if _, err := os.Stat(cfg.BASE_LOG_DIR); os.IsNotExist(err) {
		return os.Mkdir(cfg.BASE_LOG_DIR, 0755)
	}
	return nil
}
func GetNextLogFileName(cfg *LoggerConfig) (string, error) {
	dateStr := time.Now().Format(cfg.FILE_PATTERN)

	baseName := fmt.Sprintf("%s/%s_%s", cfg.BASE_LOG_DIR, "log", dateStr)

	if err := EnsureBaseLogDirectoryExist(cfg); err != nil {
		return "", err
	}
	var index int
	var logFile string

	for {
		logFile = fmt.Sprintf("%s_%d.log", baseName, index+1)
		logPath := filepath.Join(cfg.BASE_LOG_DIR, logFile)
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			return logFile, nil
		}

		fileInfo, err := os.Stat(logFile)
		if err != nil {
			return "", err
		}

		if fileInfo.Size() < int64(cfg.MAX_SIZE_MB*1024*1024) {
			return logPath, nil
		}
		index++
	}
}

func ParseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

var LogInstance *zap.Logger

func InitializeNewLogInstance(cfg *LoggerConfig) (*zap.Logger, error) {
	logFilename, err := GetNextLogFileName(cfg)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writer := zapcore.AddSync(file)

	logLevel := ParseLogLevel(cfg.LOG_LEVEL)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		writer,
		logLevel,
	)

	return zap.New(core, zap.AddCaller()), nil
}

func InitializeLogger() {
	logCfg, err := LoadLoggerConfig()
	if err != nil {
		panic("Error reading logger config " + err.Error())
	}
	LogInstance, err = InitializeNewLogInstance(logCfg)
	if err != nil {
		panic("Error initializing logger: " + err.Error())
	}
}

func CleanupQueuedLogs() {
	if LogInstance != nil {
		_ = LogInstance.Sync()
	}
}
