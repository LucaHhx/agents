package core

import (
	"fmt"
	"os"
	"path/filepath"
	"server/config"
	"server/global"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志系统
func InitLogger(cfg *config.Log) error {
	// 配置日志级别
	level := getLogLevel(cfg.Level)

	// 配置时间编码器
	timeEncoder := getTimeEncoder(cfg.TimeFormat)

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	// 创建 Core
	var cores []zapcore.Core

	// 根据输出配置创建不同的 Core
	switch cfg.Output {
	case "stdout":
		// 仅输出到控制台
		cores = append(cores, createConsoleCore(cfg, encoderConfig, level))
	case "file":
		// 仅输出到文件
		fileCores := createFileCores(cfg, encoderConfig, level)
		cores = append(cores, fileCores...)
	case "both":
		// 同时输出到控制台和文件
		cores = append(cores, createConsoleCore(cfg, encoderConfig, level))
		fileCores := createFileCores(cfg, encoderConfig, level)
		cores = append(cores, fileCores...)
	default:
		// 默认输出到控制台
		cores = append(cores, createConsoleCore(cfg, encoderConfig, level))
	}

	// 合并 Core
	core := zapcore.NewTee(cores...)

	// 创建 HZ_LOG 选项
	options := []zap.Option{}
	if cfg.ShowCaller {
		options = append(options, zap.AddCaller())
		options = append(options, zap.AddCallerSkip(0))
	}
	if cfg.ShowStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 创建 HZ_LOG
	logger := zap.New(core, options...)

	// 将日志实例赋值给全局变量
	global.HZ_LOG = logger

	global.HZ_LOG.Debug("✓ HZ_LOG 初始化成功")
	return nil
}

// createConsoleCore 创建控制台输出 Core（支持彩色）
func createConsoleCore(cfg *config.Log, encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
	// 控制台输出的编码器配置
	consoleConfig := encoderConfig

	// 配置彩色输出
	if cfg.EnableColor {
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色大写级别
	} else {
		consoleConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写级别
	}

	// 选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(consoleConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(consoleConfig)
	}

	// 创建控制台写入器
	consoleWriter := zapcore.AddSync(os.Stdout)

	return zapcore.NewCore(encoder, consoleWriter, level)
}

// createFileCores 创建文件输出 Cores（按级别分文件）
func createFileCores(cfg *config.Log, encoderConfig zapcore.EncoderConfig, level zapcore.Level) []zapcore.Core {
	var cores []zapcore.Core

	// 文件输出的编码器配置（不使用彩色）
	fileConfig := encoderConfig
	fileConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写级别

	// 选择编码器
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(fileConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(fileConfig)
	}

	if cfg.SplitByLevel {
		// 按级别分文件：debug.log, info.log, warn.log, error.log
		levels := []zapcore.Level{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
		}

		levelNames := map[zapcore.Level]string{
			zapcore.DebugLevel: "debug",
			zapcore.InfoLevel:  "info",
			zapcore.WarnLevel:  "warn",
			zapcore.ErrorLevel: "error",
		}

		for _, lvl := range levels {
			// 只创建当前配置级别及以上的日志文件
			if lvl < level {
				continue
			}

			filename := getLogFilename(cfg, levelNames[lvl])
			writer := getFileWriter(cfg, filename)

			// 创建只记录特定级别的 Core
			core := zapcore.NewCore(
				encoder,
				writer,
				zap.LevelEnablerFunc(func(l zapcore.Level) bool {
					return l == lvl
				}),
			)
			cores = append(cores, core)
		}
	} else {
		// 不分级别，所有日志写入一个文件
		filename := getLogFilename(cfg, "app")
		writer := getFileWriter(cfg, filename)
		core := zapcore.NewCore(encoder, writer, level)
		cores = append(cores, core)
	}

	return cores
}

// getLogFilename 获取日志文件名
func getLogFilename(cfg *config.Log, levelName string) string {
	var logPath string

	if cfg.SplitByDate {
		// 按日期分文件夹：logs/2024-12-14/info.log
		dateStr := time.Now().Format("2006-01-02")
		logPath = filepath.Join(cfg.LogDir, dateStr, levelName+".log")
	} else {
		// 不按日期分文件夹：logs/info.log
		logPath = filepath.Join(cfg.LogDir, levelName+".log")
	}

	return logPath
}

// getFileWriter 获取文件写入器（使用 lumberjack 实现日志轮转）
func getFileWriter(cfg *config.Log, filename string) zapcore.WriteSyncer {
	// 确保日志目录存在
	logDir := filepath.Dir(filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		return nil
	}

	// 配置日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,       // 日志文件路径
		MaxSize:    cfg.MaxSize,    // 单个文件最大尺寸（MB）
		MaxBackups: cfg.MaxBackups, // 保留的旧文件数量
		MaxAge:     cfg.MaxAge,     // 保留的天数
		Compress:   cfg.Compress,   // 是否压缩
		LocalTime:  true,           // 使用本地时间
	}

	return zapcore.AddSync(lumberJackLogger)
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
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

// getTimeEncoder 获取时间编码器
func getTimeEncoder(format string) zapcore.TimeEncoder {
	if format == "" {
		return zapcore.ISO8601TimeEncoder
	}
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

// customCallerEncoder 自定义调用者编码器，显示完整路径（包含 server/）
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}

	// 格式化路径：保留 server/ 前缀
	fullPath := caller.File
	shortPath := formatLogFilePath(fullPath)

	// 格式：server/path/file.go:123
	enc.AppendString(fmt.Sprintf("%s:%d", shortPath, caller.Line))
}

// formatLogFilePath 格式化日志文件路径
func formatLogFilePath(fullPath string) string {
	// 查找 /server/ 标记
	if idx := strings.Index(fullPath, "/server/"); idx != -1 {
		// 保留 server/ 前缀
		return fullPath[idx+1:] // 跳过第一个 "/"
	}
	// 如果没有 /server/，返回文件名
	return filepath.Base(fullPath)
}
