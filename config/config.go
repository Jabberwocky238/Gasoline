package config

// Interface 配置结构体
type Interface struct {
	PrivateKey string `toml:"PrivateKey"`
	ListenPort int    `toml:"ListenPort"`
	Address    string `toml:"Address"`
}

// Peer 配置结构体
type Peer struct {
	PublicKey  string `toml:"PublicKey"`
	AllowedIPs string `toml:"AllowedIPs"`
	Endpoint   string `toml:"Endpoint"`
}

type Transport struct {
	ID         string                 `toml:"ID"`
	Type       string                 `toml:"Type"`
	Main       bool                   `toml:"Main,omitempty"`
	Underlying string                 `toml:"Underlying,omitempty"`
	Cfg        map[string]interface{} `toml:"Cfg,omitempty"`
}

// Config 主配置结构体
type Config struct {
	Interface  Interface   `toml:"Interface"`
	Peers      []Peer      `toml:"Peer"`
	Transports []Transport `toml:"Transport"`
}
