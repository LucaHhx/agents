package config

// Server 服务器配置
type Server struct {
	Host           string `yaml:"host"`           // 服务器监听地址
	Port           int    `yaml:"port"`           // 服务器监听端口
	Mode           string `yaml:"mode"`           // 运行模式: debug, release, test
	ReadTimeout    int    `yaml:"readTimeout"`    // 读取超时时间（秒）
	WriteTimeout   int    `yaml:"writeTimeout"`   // 写入超时时间（秒）
	MaxHeaderBytes int    `yaml:"maxHeaderBytes"` // 最大请求头大小（MB）
}

// Validate 验证服务器配置
func (s *Server) Validate() error {
	return nil
}
