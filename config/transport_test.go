package config

import (
	"context"
	"testing"
	"wwww/transport/trojan"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

// go test -v ./config -run TestFromConfig -timeout 5s
func TestFromConfig(t *testing.T) {
	cfg := `
[[Transport]]
ID = "tcp-id"
Type = "tcp"

[[Transport]]
ID = "caesar-id"
Type = "caesar"
Underlying = "tls-id"
Cfg.Shift = 3

[[Transport]]
ID = "trojan-id"
Type = "trojan"
Main = true
Underlying = "caesar-id"
Cfg.Password = "password"
Cfg.Passwords = ["passwordss"]

[[Transport]]
ID = "tls-id"
Type = "tls"
Cfg.ServerName = "localhost"
Cfg.CertFile = "../samples/cert.pem"
Cfg.KeyFile = "../samples/key.pem"
Cfg.InsecureSkipVerify = true
Cfg.SNI = true
	`
	parsed, err := ParseConfigFromString(cfg)
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}
	assert.Equal(t, parsed.Transports, []Transport{
		{
			ID:   "tcp-id",
			Type: "tcp",
		},
		{
			ID:         "caesar-id",
			Type:       "caesar",
			Underlying: "tls-id",
			Cfg: map[string]any{
				"Shift": int64(3),
			},
		},
		{
			ID:         "trojan-id",
			Type:       "trojan",
			Main:       true,
			Underlying: "caesar-id",
			Cfg: map[string]any{
				"Password":  "password",
				"Passwords": []interface{}{"passwordss"}, // []string
			},
		},
		{
			ID:   "tls-id",
			Type: "tls",
			Cfg: map[string]any{
				"ServerName":         "localhost",
				"CertFile":           "../samples/cert.pem",
				"KeyFile":            "../samples/key.pem",
				"InsecureSkipVerify": true,
				"SNI":                true,
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, err := FromConfigServer(ctx, parsed.Transports)
	if err != nil {
		t.Fatalf("创建服务器失败: %v", err)
	}
	client, err := FromConfigClient(ctx, parsed.Transports)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	if _, ok := server.(*trojan.TrojanServer); !ok {
		t.Fatalf("创建服务器失败: %v", err)
	}
	if _, ok := client.(*trojan.TrojanClient); !ok {
		t.Fatalf("创建客户端失败: %v", err)
	}
}

func TestStringListParsing(t *testing.T) {
	cfg := `
Cfg.Passwords = ["passwordss"]
`
	var config map[string]any
	if _, err := toml.Decode(cfg, &config); err != nil {
		t.Fatalf("解析 TOML 配置失败: %v", err)
	}

	// 验证 Cfg 存在
	cfgMap, ok := config["Cfg"].(map[string]any)
	assert.True(t, ok, "Cfg should be map[string]any")

	// 验证 Passwords 存在且类型为 []string
	passwords, ok := cfgMap["Passwords"]
	assert.True(t, ok, "Passwords should exist")

	// TOML 解析字符串数组为 []string 类型
	passwordsSlice, ok := passwords.([]string)
	assert.True(t, ok, "Passwords should be []string, got %T", passwords)
	assert.Equal(t, []string{"passwordss"}, passwordsSlice)
}
