package config

// Mysql MySQL数据库配置
type Mysql struct {
	// 基础连接配置
	Host     string `yaml:"host"`     // 数据库服务器地址
	Port     int    `yaml:"port"`     // 数据库服务器端口
	User     string `yaml:"user"`     // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	Dbname   string `yaml:"dbname"`   // 数据库名称

	// 连接参数配置
	Charset   string `yaml:"charset"`   // 字符集
	ParseTime string `yaml:"parseTime"` // 是否解析时间
	Loc       string `yaml:"loc"`       // 时区设置

	// GORM MySQL 驱动配置
	DefaultStringSize         int  `yaml:"defaultStringSize"`         // 字符串字段默认长度
	DisableDatetimePrecision  bool `yaml:"disableDatetimePrecision"`  // 是否禁用日期时间精度
	DontSupportRenameIndex    bool `yaml:"dontSupportRenameIndex"`    // 不使用 ALTER INDEX 重命名索引
	DontSupportRenameColumn   bool `yaml:"dontSupportRenameColumn"`   // 不使用 ALTER COLUMN 重命名列
	SkipInitializeWithVersion bool `yaml:"skipInitializeWithVersion"` // 是否跳过根据版本自动配置

	// 连接池配置
	MaxIdleConns    int `yaml:"maxIdleConns"`    // 最大空闲连接数
	MaxOpenConns    int `yaml:"maxOpenConns"`    // 最大打开连接数
	ConnMaxLifetime int `yaml:"connMaxLifetime"` // 连接最大生存时间（秒）
	ConnMaxIdleTime int `yaml:"connMaxIdleTime"` // 连接最大空闲时间（秒）

	// 日志配置
	LogEnabled     bool   `yaml:"logEnabled"`     // 是否启用日志
	LogLevel       string `yaml:"logLevel"`       // 日志级别: silent, error, warn, info
	LogSlowQuery   int    `yaml:"logSlowQuery"`   // 慢查询阈值（毫秒），0表示不记录慢查询
	LogColorful    bool   `yaml:"logColorful"`    // 是否启用彩色日志
	LogFile        string `yaml:"logFile"`        // 日志文件名（不含扩展名）
	LogIgnoreTrace bool   `yaml:"logIgnoreTrace"` // 是否忽略未找到记录的错误日志
}

// Validate 验证MySQL配置
func (m *Mysql) Validate() error {
	return nil
}
