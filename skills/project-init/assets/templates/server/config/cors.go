package config

// CORS 跨域配置
type CORS struct {
	Enabled          bool     `yaml:"enabled"`          // 是否启用跨域
	AllowOrigins     []string `yaml:"allowOrigins"`     // 允许的来源
	AllowMethods     []string `yaml:"allowMethods"`     // 允许的HTTP方法
	AllowHeaders     []string `yaml:"allowHeaders"`     // 允许的请求头
	ExposeHeaders    []string `yaml:"exposeHeaders"`    // 暴露的响应头
	AllowCredentials bool     `yaml:"allowCredentials"` // 是否允许携带认证信息
	MaxAge           int      `yaml:"maxAge"`           // 预检请求缓存时间（秒）
}

// Validate 验证跨域配置
func (c *CORS) Validate() error {
	return nil
}
