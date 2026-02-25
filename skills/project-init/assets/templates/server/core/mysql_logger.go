package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"server/config"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLLogger 自定义MySQL日志记录器
type MySQLLogger struct {
	logger                    *zap.Logger
	logLevel                  logger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
	colorful                  bool
	projectRoot               string
}

// NewMySQLLogger 创建MySQL日志记录器
func NewMySQLLogger(cfg *config.Mysql, logCfg *config.Log, sysCfg *config.System) logger.Interface {
	// 如果未启用MySQL日志，返回静默logger
	if !cfg.LogEnabled {
		return logger.Default.LogMode(logger.Silent)
	}

	// 获取日志级别
	logLevel := getGormLogLevel(cfg.LogLevel)

	// 获取慢查询阈值
	slowThreshold := time.Duration(cfg.LogSlowQuery) * time.Millisecond

	// 获取项目根目录
	projectRoot := ""
	if sysCfg != nil {
		projectRoot = sysCfg.ProjectRoot
	}

	// 配置时间编码器
	var timeEncoder zapcore.TimeEncoder
	if logCfg.TimeFormat != "" {
		timeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(logCfg.TimeFormat))
		}
	} else {
		timeEncoder = zapcore.ISO8601TimeEncoder
	}

	// 配置编码器基础配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// 不使用 CallerKey，因为 GORM 通过 utils.FileWithLineNum() 提供了更准确的调用位置
	}

	// 创建 Cores
	var cores []zapcore.Core

	// 根据输出配置创建不同的 Core
	switch logCfg.Output {
	case "stdout":
		// 仅输出到控制台
		cores = append(cores, createMySQLConsoleCore(cfg, encoderConfig))
	case "file":
		// 仅输出到文件
		fileCore := createMySQLFileCore(cfg, logCfg, encoderConfig)
		if fileCore != nil {
			cores = append(cores, fileCore)
		}
	case "both":
		// 同时输出到控制台和文件
		cores = append(cores, createMySQLConsoleCore(cfg, encoderConfig))
		fileCore := createMySQLFileCore(cfg, logCfg, encoderConfig)
		if fileCore != nil {
			cores = append(cores, fileCore)
		}
	default:
		// 默认输出到控制台
		cores = append(cores, createMySQLConsoleCore(cfg, encoderConfig))
	}

	// 如果没有可用的 Core，返回静默 logger
	if len(cores) == 0 {
		return logger.Default.LogMode(logger.Silent)
	}

	// 合并 Core
	core := zapcore.NewTee(cores...)

	// 创建Logger
	zapLogger := zap.New(core)

	return &MySQLLogger{
		logger:                    zapLogger,
		logLevel:                  logLevel,
		slowThreshold:             slowThreshold,
		ignoreRecordNotFoundError: cfg.LogIgnoreTrace,
		colorful:                  cfg.LogColorful,
		projectRoot:               projectRoot,
	}
}

// createMySQLConsoleCore 创建控制台输出 Core
func createMySQLConsoleCore(cfg *config.Mysql, encoderConfig zapcore.EncoderConfig) zapcore.Core {
	// 控制台输出的编码器配置
	consoleConfig := encoderConfig

	// 配置彩色输出
	if cfg.LogColorful {
		consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		consoleConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	// 创建编码器（使用console格式，更易读）
	encoder := zapcore.NewConsoleEncoder(consoleConfig)

	// 创建控制台写入器
	consoleWriter := zapcore.AddSync(os.Stdout)

	return zapcore.NewCore(encoder, consoleWriter, zapcore.DebugLevel)
}

// createMySQLFileCore 创建文件输出 Core
func createMySQLFileCore(cfg *config.Mysql, logCfg *config.Log, encoderConfig zapcore.EncoderConfig) zapcore.Core {
	// 创建日志文件路径
	var logPath string
	if logCfg.SplitByDate {
		// 按日期分文件夹：logs/2025-12-14/mysql.log
		dateStr := time.Now().Format("2006-01-02")
		logPath = filepath.Join(logCfg.LogDir, dateStr, cfg.LogFile+".log")
	} else {
		// 不按日期分文件夹：logs/mysql.log
		logPath = filepath.Join(logCfg.LogDir, cfg.LogFile+".log")
	}

	// 确保日志目录存在
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建MySQL日志目录失败: %v\n", err)
		return nil
	}

	// 配置日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,           // 日志文件路径
		MaxSize:    logCfg.MaxSize,    // 单个文件最大尺寸（MB）
		MaxBackups: logCfg.MaxBackups, // 保留的旧文件数量
		MaxAge:     logCfg.MaxAge,     // 保留的天数
		Compress:   logCfg.Compress,   // 是否压缩
		LocalTime:  true,              // 使用本地时间
	}

	// 文件输出的编码器配置（不使用彩色）
	fileConfig := encoderConfig
	fileConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建编码器（使用console格式，更易读）
	encoder := zapcore.NewConsoleEncoder(fileConfig)

	// 创建文件写入器
	writer := zapcore.AddSync(lumberJackLogger)

	return zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
}

// LogMode 设置日志级别
func (l *MySQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录Info级别日志
func (l *MySQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 记录Warn级别日志
func (l *MySQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 记录Error级别日志
func (l *MySQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 记录SQL执行日志
func (l *MySQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取调用位置并转换为相对路径
	fileWithLineNum := getCallerInfo()
	fileWithLineNum = formatFilePath(fileWithLineNum, l.projectRoot)

	switch {
	case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		// 记录错误 - 格式化为多行
		msg := fmt.Sprintf("%s SQL执行错误 %s\n  %s | %s | %s",
			fileWithLineNum,
			err.Error(),
			formatElapsed(elapsed, l.colorful),
			formatRows(rows, l.colorful),
			sql,
		)
		l.logger.Error(msg)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		// 记录慢查询 - 格式化为多行
		msg := fmt.Sprintf("%s 慢查询 (阈值: %s)\n  %s | %s | %s",
			fileWithLineNum,
			formatDuration(l.slowThreshold, l.colorful),
			formatElapsed(elapsed, l.colorful),
			formatRows(rows, l.colorful),
			sql,
		)
		l.logger.Warn(msg)
	case l.logLevel == logger.Info:
		// 记录普通SQL - 格式化为多行
		msg := fmt.Sprintf("%s SQL执行\n  %s | %s | %s",
			fileWithLineNum,
			formatElapsed(elapsed, l.colorful),
			formatRows(rows, l.colorful),
			sql,
		)
		l.logger.Info(msg)
	}
}

// getCallerInfo 获取真实的调用者信息（跳过 GORM 内部调用栈）
func getCallerInfo() string {
	for i := 2; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			// 跳过 mysql_logger.go 本身
			if strings.Contains(file, "mysql_logger.go") {
				continue
			}

			// 跳过所有 Go modules 中的文件（包括 gorm.io、其他第三方库）
			if strings.Contains(file, "/go/pkg/mod/") {
				continue
			}

			// 返回找到的第一个项目代码位置
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	// 如果没找到，返回默认值
	return "unknown:0"
}

// getGormLogLevel 获取GORM日志级别
func getGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

// formatFilePath 格式化文件路径，转换为编辑器可识别的相对路径
func formatFilePath(fullPathWithLine string, projectRoot string) string {
	// 分离文件路径和行号 (格式: /path/to/file.go:123)
	parts := strings.Split(fullPathWithLine, ":")
	if len(parts) < 2 {
		return fullPathWithLine // 如果没有行号，直接返回
	}

	fullPath := parts[0]
	lineNum := parts[len(parts)-1] // 最后一部分是行号

	var shortPath string

	// 1. 优先使用配置的项目根目录去除前缀
	if projectRoot != "" && strings.HasPrefix(fullPath, projectRoot) {
		// 去除项目根目录前缀，保留相对路径
		shortPath = strings.TrimPrefix(fullPath, projectRoot)
		// 去除开头的斜杠
		shortPath = strings.TrimPrefix(shortPath, "/")
	} else if idx := strings.Index(fullPath, "/server/"); idx != -1 {
		// 2. 后备方案：处理项目内的文件 - 保留 server/ 前缀
		shortPath = fullPath[idx+1:] // 跳过第一个 "/"，保留 "server/"
	} else if strings.Contains(fullPath, "/go/pkg/mod/") {
		// 2. 处理 Go modules 中的文件（如 gorm.io/driver/mysql）
		if idx := strings.Index(fullPath, "/go/pkg/mod/"); idx != -1 {
			modPath := fullPath[idx+12:] // 跳过 "/go/pkg/mod/" (长度为12)
			// 进一步简化：只保留包名和文件名
			// 例如：gorm.io/driver/mysql@v1.6.0/migrator.go:414 -> gorm.io/driver/mysql/migrator.go:414
			atParts := strings.Split(modPath, "@")
			if len(atParts) > 1 {
				// 移除版本号
				pkgName := atParts[0]
				// 查找版本号后的路径
				if versionEnd := strings.Index(atParts[1], "/"); versionEnd != -1 {
					shortPath = pkgName + atParts[1][versionEnd:]
				} else {
					shortPath = modPath
				}
			} else {
				shortPath = modPath
			}
		}
	} else {
		// 3. 处理其他项目相关路径
		markers := []string{"/report/", "/modules/", "/tests/"}
		found := false
		for _, marker := range markers {
			if idx := strings.Index(fullPath, marker); idx != -1 {
				shortPath = fullPath[idx+1:]
				found = true
				break
			}
		}
		// 4. 如果都没匹配，只返回文件名
		if !found {
			shortPath = filepath.Base(fullPath)
		}
	}

	// 重新组合路径和行号
	return shortPath + ":" + lineNum
}

// ANSI 颜色代码
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
)

// formatElapsed 格式化执行时间
func formatElapsed(elapsed time.Duration, colorful bool) string {
	var color string
	if colorful {
		switch {
		case elapsed < 100*time.Millisecond:
			color = colorGreen // 快速查询 - 绿色
		case elapsed < 500*time.Millisecond:
			color = colorYellow // 中等查询 - 黄色
		default:
			color = colorRed // 慢查询 - 红色
		}
	}

	timeStr := formatDurationSimple(elapsed)
	if colorful {
		return fmt.Sprintf("%s%s%s", color, timeStr, colorReset)
	}
	return timeStr
}

// formatRows 格式化影响行数
func formatRows(rows int64, colorful bool) string {
	var color string
	if colorful {
		switch {
		case rows == 0:
			color = colorGray // 无数据 - 灰色
		case rows < 10:
			color = colorCyan // 少量数据 - 青色
		case rows < 100:
			color = colorBlue // 中等数据 - 蓝色
		default:
			color = colorMagenta // 大量数据 - 紫色
		}
	}

	rowStr := fmt.Sprintf("%d rows", rows)
	if colorful {
		return fmt.Sprintf("%s%s%s", color, rowStr, colorReset)
	}
	return rowStr
}

// formatDuration 格式化时间阈值
func formatDuration(d time.Duration, colorful bool) string {
	timeStr := formatDurationSimple(d)
	if colorful {
		return fmt.Sprintf("%s%s%s", colorYellow, timeStr, colorReset)
	}
	return timeStr
}

// formatDurationSimple 简化时间格式
func formatDurationSimple(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%dns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%.2fus", float64(d.Nanoseconds())/1000)
	case d < time.Second:
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}
