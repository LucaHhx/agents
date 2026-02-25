package config

// Config 应用程序总配置
type Config struct {
	System System `yaml:"system"` // 系统配置
	Server Server `yaml:"server"` // 服务器配置
	Mysql  Mysql  `yaml:"mysql"`  // MySQL数据库配置
	Redis  Redis  `yaml:"redis"`  // Redis缓存配置
	Log    Log    `yaml:"log"`    // 日志配置
	CORS   CORS   `yaml:"cors"`   // 跨域配置
	Upload Upload `yaml:"upload"` // 文件上传配置
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	// 验证服务器配置
	if err := c.Server.Validate(); err != nil {
		return err
	}
	// 验证MySQL配置
	if err := c.Mysql.Validate(); err != nil {
		return err
	}
	// 验证Redis配置
	if err := c.Redis.Validate(); err != nil {
		return err
	}
	return nil
}
