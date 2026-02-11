package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"dove/pkg/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// MultiWriter 多写入器，支持同时写入多个目标
type MultiWriter struct {
	writers []io.Writer
}

// NewMultiWriter 创建多写入器
func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

// Write 实现 io.Writer 接口
func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return len(p), nil
}

var Logger *slog.Logger

// InitLogger 初始化日志
func InitLogger() {
	var level slog.Level
	switch config.GlobalConfig.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 根据配置选择输出目标
	writer := createLogWriter(config.GlobalConfig.Log)

	var handler slog.Handler
	if config.GlobalConfig.Log.Format == "json" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// createLogWriter 创建日志写入器
func createLogWriter(logConfig config.LogConfig) io.Writer {
	var writers []io.Writer

	switch logConfig.Output {
	case "stdout":
		writers = append(writers, os.Stdout)
	case "stderr":
		writers = append(writers, os.Stderr)
	case "file":
		// 确保日志目录存在
		if err := ensureLogDir(logConfig.File.Path); err != nil {
			// 如果无法创建目录，回退到 stdout
			slog.Error("Failed to create log directory, falling back to stdout", "error", err)
			writers = append(writers, os.Stdout)
		} else {
			// 创建带轮转的日志写入器
			writers = append(writers, &lumberjack.Logger{
				Filename:   logConfig.File.Path,
				MaxSize:    logConfig.File.MaxSize, // MB
				MaxAge:     logConfig.File.MaxAge,  // days
				MaxBackups: logConfig.File.MaxBackups,
				Compress:   true, // 压缩旧日志文件
			})
		}
	case "both":
		// 同时输出到控制台和文件
		writers = append(writers, os.Stdout)

		// 确保日志目录存在
		if err := ensureLogDir(logConfig.File.Path); err != nil {
			// 如果无法创建目录，只输出到控制台
			slog.Error("Failed to create log directory, outputting to console only", "error", err)
		} else {
			// 创建带轮转的日志写入器
			writers = append(writers, &lumberjack.Logger{
				Filename:   logConfig.File.Path,
				MaxSize:    logConfig.File.MaxSize, // MB
				MaxAge:     logConfig.File.MaxAge,  // days
				MaxBackups: logConfig.File.MaxBackups,
				Compress:   true, // 压缩旧日志文件
			})
		}
	case "all":
		// 输出到所有目标：控制台、错误输出和文件
		writers = append(writers, os.Stdout, os.Stderr)

		// 确保日志目录存在
		if err := ensureLogDir(logConfig.File.Path); err != nil {
			// 如果无法创建目录，只输出到控制台
			slog.Error("Failed to create log directory, outputting to console only", "error", err)
		} else {
			// 创建带轮转的日志写入器
			writers = append(writers, &lumberjack.Logger{
				Filename:   logConfig.File.Path,
				MaxSize:    logConfig.File.MaxSize, // MB
				MaxAge:     logConfig.File.MaxAge,  // days
				MaxBackups: logConfig.File.MaxBackups,
				Compress:   true, // 压缩旧日志文件
			})
		}
	default:
		// 默认输出到 stdout
		writers = append(writers, os.Stdout)
	}

	// 如果只有一个写入器，直接返回
	if len(writers) == 1 {
		return writers[0]
	}

	// 多个写入器，使用 MultiWriter
	return NewMultiWriter(writers...)
}

// ensureLogDir 确保日志目录存在
func ensureLogDir(logPath string) error {
	dir := filepath.Dir(logPath)
	if dir == "." {
		// 如果路径中没有目录，直接返回
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// getCallerInfo 获取调用者信息
func getCallerInfo(skip int) (file string, line int, function string) {
	// 获取调用栈信息，跳过指定数量的函数
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}

	// 获取函数名
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return filepath.Base(file), line, "unknown"
	}

	// 获取相对路径的文件名
	fileName := filepath.Base(file)

	// 获取函数名（去掉包路径）
	funcName := fn.Name()
	if idx := filepath.Ext(funcName); idx != "" {
		funcName = funcName[:len(funcName)-len(idx)]
	}

	return fileName, line, funcName
}

// Debug 记录调试日志
func Debug(msg string, args ...any) {
	file, line, function := getCallerInfo(2)
	Logger.With("file", file, "line", line, "function", function).Debug(msg, args...)
}

// Info 记录信息日志
func Info(msg string, args ...any) {
	file, line, function := getCallerInfo(2)
	Logger.With("file", file, "line", line, "function", function).Info(msg, args...)
}

// Warn 记录警告日志
func Warn(msg string, args ...any) {
	file, line, function := getCallerInfo(2)
	Logger.With("file", file, "line", line, "function", function).Warn(msg, args...)
}

// Error 记录错误日志
func Error(msg string, args ...any) {
	file, line, function := getCallerInfo(2)
	Logger.With("file", file, "line", line, "function", function).Error(msg, args...)
}

// WithGroup 创建带分组的日志器
func WithGroup(name string) *slog.Logger {
	return Logger.WithGroup(name)
}

// With 创建带字段的日志器
func With(args ...any) *slog.Logger {
	return Logger.With(args...)
}

// GetLogStatus 获取日志状态信息
func GetLogStatus() map[string]interface{} {
	status := map[string]interface{}{
		"level":  config.GlobalConfig.Log.Level,
		"format": config.GlobalConfig.Log.Format,
		"output": config.GlobalConfig.Log.Output,
	}

	// 解析输出目标
	outputs := parseOutputTargets(config.GlobalConfig.Log.Output)
	status["targets"] = outputs

	// 如果是文件输出或包含文件输出，添加文件信息
	if config.GlobalConfig.Log.Output == "file" || config.GlobalConfig.Log.Output == "both" || config.GlobalConfig.Log.Output == "all" {
		fileInfo, err := os.Stat(config.GlobalConfig.Log.File.Path)
		if err == nil {
			status["file"] = map[string]interface{}{
				"path":     config.GlobalConfig.Log.File.Path,
				"size":     fileInfo.Size(),
				"mod_time": fileInfo.ModTime(),
			}
		} else {
			status["file"] = map[string]interface{}{
				"path":  config.GlobalConfig.Log.File.Path,
				"error": "file not found or not accessible",
			}
		}

		// 添加轮转配置
		status["rotation"] = map[string]interface{}{
			"max_size":    config.GlobalConfig.Log.File.MaxSize,
			"max_age":     config.GlobalConfig.Log.File.MaxAge,
			"max_backups": config.GlobalConfig.Log.File.MaxBackups,
		}
	}

	return status
}

// parseOutputTargets 解析输出目标
func parseOutputTargets(output string) []string {
	switch output {
	case "stdout":
		return []string{"stdout"}
	case "stderr":
		return []string{"stderr"}
	case "file":
		return []string{"file"}
	case "both":
		return []string{"stdout", "file"}
	case "all":
		return []string{"stdout", "stderr", "file"}
	default:
		return []string{"stdout"}
	}
}

// ValidateLogConfig 验证日志配置
func ValidateLogConfig() error {
	logConfig := config.GlobalConfig.Log

	// 验证日志级别
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[logConfig.Level] {
		return fmt.Errorf("invalid log level: %s", logConfig.Level)
	}

	// 验证日志格式
	validFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validFormats[logConfig.Format] {
		return fmt.Errorf("invalid log format: %s", logConfig.Format)
	}

	// 验证输出目标
	validOutputs := map[string]bool{
		"stdout": true, "stderr": true, "file": true, "both": true, "all": true,
	}
	if !validOutputs[logConfig.Output] {
		return fmt.Errorf("invalid log output: %s, valid options: stdout, stderr, file, both, all", logConfig.Output)
	}

	// 如果是文件输出或包含文件输出，验证文件配置
	if logConfig.Output == "file" || logConfig.Output == "both" || logConfig.Output == "all" {
		if logConfig.File.Path == "" {
			return fmt.Errorf("log file path is required when output is 'file', 'both', or 'all'")
		}
		if logConfig.File.MaxSize <= 0 {
			return fmt.Errorf("log file max_size must be greater than 0")
		}
		if logConfig.File.MaxAge < 0 {
			return fmt.Errorf("log file max_age must be non-negative")
		}
		if logConfig.File.MaxBackups < 0 {
			return fmt.Errorf("log file max_backups must be non-negative")
		}
	}

	return nil
}
