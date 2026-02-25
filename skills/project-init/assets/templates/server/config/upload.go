package config

// Upload 文件上传配置
type Upload struct {
	MaxSize    int      `yaml:"maxSize"`    // 最大文件大小（MB）
	AllowTypes []string `yaml:"allowTypes"` // 允许的文件类型
	SavePath   string   `yaml:"savePath"`   // 文件保存路径
}

// Validate 验证文件上传配置
func (u *Upload) Validate() error {
	return nil
}
