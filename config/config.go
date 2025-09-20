package config

import "log"

// Interface 配置结构体
type Interface struct {
	UniqueID   string `toml:"UniqueID"`
	ListenPort int    `toml:"ListenPort"`
	Address    string `toml:"Address"`
}

// Peer 配置结构体
type Peer struct {
	UniqueID   string `toml:"UniqueID"`
	AllowedIPs string `toml:"AllowedIPs"`
	Endpoint   string `toml:"Endpoint"`
}

// Config 主配置结构体
type Config struct {
	Interface Interface `toml:"Interface"`
	Peers     []Peer    `toml:"Peer"`
}

func NewConfig(configPath string) *Config {
	config, err := ParseConfig(configPath)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	return config
}
