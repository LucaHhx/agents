package config

// Log 日志配置
type Log struct {
	Level          string `yaml:"level"`          // 日志级别: debug, info, warn, error
	Format         string `yaml:"format"`         // 日志格式: json, console
	Output         string `yaml:"output"`         // 输出位置: stdout, file, both
	LogDir         string `yaml:"logDir"`         // 日志文件根目录
	EnableColor    bool   `yaml:"enableColor"`    // 是否启用彩色输出
	ShowCaller     bool   `yaml:"showCaller"`     // 是否显示调用者信息
	ShowStacktrace bool   `yaml:"showStacktrace"` // 是否显示堆栈跟踪
	TimeFormat     string `yaml:"timeFormat"`     // 时间格式
	SplitByLevel   bool   `yaml:"splitByLevel"`   // 是否按日志级别分文件
	SplitByDate    bool   `yaml:"splitByDate"`    // 是否按日期分文件夹
	MaxSize        int    `yaml:"maxSize"`        // 单个日志文件最大大小（MB）
	MaxBackups     int    `yaml:"maxBackups"`     // 保留的旧日志文件数量
	MaxAge         int    `yaml:"maxAge"`         // 日志文件保留天数
	Compress       bool   `yaml:"compress"`       // 是否压缩旧日志文件
}

// Validate 验证日志配置
func (l *Log) Validate() error {
	return nil
}
