package config

// Redis 缓存配置
type Redis struct {
	Host         string `yaml:"host"`         // Redis服务器地址
	Port         int    `yaml:"port"`         // Redis服务器端口
	Password     string `yaml:"password"`     // Redis密码
	DB           int    `yaml:"db"`           // 使用的数据库编号
	PoolSize     int    `yaml:"poolSize"`     // 连接池大小
	MinIdleConns int    `yaml:"minIdleConns"` // 最小空闲连接数
	MaxRetries   int    `yaml:"maxRetries"`   // 最大重试次数
	DialTimeout  int    `yaml:"dialTimeout"`  // 连接超时时间（秒）
	ReadTimeout  int    `yaml:"readTimeout"`  // 读取超时时间（秒）
	WriteTimeout int    `yaml:"writeTimeout"` // 写入超时时间（秒）
}

// Validate 验证Redis配置
func (r *Redis) Validate() error {
	return nil
}
