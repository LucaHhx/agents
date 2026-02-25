package config

// System 系统配置
type System struct {
	ProjectRoot string `yaml:"projectRoot"` // 项目根目录（用于日志中显示相对路径）
	Migrate     bool   `yaml:"migrate"`     // 是否自动迁移数据库
	Redis       bool   `yaml:"redis"`       // 是否使用 Redis
	Name        string `yaml:"name"`
}

// Validate 验证系统配置
func (s *System) Validate() error {
	return nil
}
