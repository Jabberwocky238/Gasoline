package config

// Interface 配置结构体
type Interface struct {
	PrivateKey  string `toml:"PrivateKey"`
	ListenPort  int    `toml:"ListenPort,omitempty"`
	Address     string `toml:"Address"`
	TransportID string `toml:"TransportID,omitempty"`
}

// Peer 配置结构体
type Peer struct {
	PublicKey   string `toml:"PublicKey"`
	AllowedIPs  string `toml:"AllowedIPs"`
	Endpoint    string `toml:"Endpoint"`
	TransportID string `toml:"TransportID,omitempty"`
}

type Transport struct {
	ID         string                 `toml:"ID"`
	Type       string                 `toml:"Type"`
	Underlying string                 `toml:"Underlying,omitempty"`
	Cfg        map[string]interface{} `toml:"Cfg,omitempty"`
}

// Config 主配置结构体
type Config struct {
	Interface  Interface   `toml:"Interface"`
	Peers      []Peer      `toml:"Peer"`
	Transports []Transport `toml:"Transport"`
}
