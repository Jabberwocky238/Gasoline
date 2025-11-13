package config

import (
	"context"
	"fmt"
	"os"

	"wwww/transport"
	"wwww/transport/caesar"
	"wwww/transport/tcp"
	"wwww/transport/tls"
	"wwww/transport/trojan"
)

// makeDependenciesList 递归构建依赖列表，确保无依赖的项在前，有依赖的项在后
// 返回的列表顺序：无依赖的配置 -> 有依赖的配置（最上层在前）
func makeDependenciesList(cfgMap map[string]Transport, cfgs []Transport) ([]Transport, error) {
	result := make([]Transport, 0)
	visited := make(map[string]bool)
	processing := make(map[string]bool) // 用于检测循环依赖

	// 递归函数：处理单个配置及其依赖
	var processCfg func(cfgID string) error
	processCfg = func(cfgID string) error {
		// 检查是否存在
		cfg, exists := cfgMap[cfgID]
		if !exists {
			return nil // 依赖不存在，跳过
		}

		// 检测循环依赖
		if processing[cfgID] {
			return fmt.Errorf("circular dependency detected: %s", cfgID)
		}

		// 如果已经处理过，跳过
		if visited[cfgID] {
			return nil
		}

		// 标记为正在处理
		processing[cfgID] = true

		// 先处理依赖（递归）
		if cfg.Underlying != "" {
			if err := processCfg(cfg.Underlying); err != nil {
				return err
			}
		}

		// 标记为已处理
		processing[cfgID] = false
		visited[cfgID] = true

		// 添加到结果列表（依赖已处理，现在可以安全添加）
		result = append(result, cfg)
		return nil
	}

	// 处理所有配置
	for _, cfg := range cfgs {
		if err := processCfg(cfg.ID); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func FromConfigServer(ctx context.Context, cfgs []Transport) (transport.TransportServer, error) {
	cfgMap := make(map[string]Transport)
	serverMap := make(map[string]transport.TransportServer)
	for _, cfg := range cfgs {
		cfgMap[cfg.ID] = cfg
	}
	dependenciesList, err := makeDependenciesList(cfgMap, cfgs)
	if err != nil {
		return nil, err
	}

	makeServer := func(cfg Transport) (transport.TransportServer, error) {
		switch cfg.Type {
		case "tcp":
			return tcp.NewTCPServer(), nil
		case "tls":
			tlsCfg := new(tls.TLSServerConfig)

			// 优先使用配置中直接提供的 CertPem 和 KeyPem
			if certPem, ok := cfg.Cfg["CertPem"].([]byte); ok && len(certPem) > 0 {
				tlsCfg.CertPem = certPem
			} else if certFile, ok := cfg.Cfg["CertFile"].(string); ok && certFile != "" {
				// 如果未直接提供，则从文件读取
				certPem, err := os.ReadFile(certFile)
				if err != nil {
					return nil, err
				}
				tlsCfg.CertPem = certPem
			}

			if keyPem, ok := cfg.Cfg["KeyPem"].([]byte); ok && len(keyPem) > 0 {
				tlsCfg.KeyPem = keyPem
			} else if keyFile, ok := cfg.Cfg["KeyFile"].(string); ok && keyFile != "" {
				// 如果未直接提供，则从文件读取
				keyPem, err := os.ReadFile(keyFile)
				if err != nil {
					return nil, err
				}
				tlsCfg.KeyPem = keyPem
			}

			// 设置其他字段
			if serverName, ok := cfg.Cfg["ServerName"].(string); ok {
				tlsCfg.ServerName = serverName
			}

			return tls.NewTLSServer(tlsCfg), nil
		case "caesar":
			caesarCfg := &caesar.CaesarConfig{
				Shift: int(cfg.Cfg["Shift"].(int64)),
			}
			underlyingServer := serverMap[cfg.Underlying]
			return caesar.NewCaesarServer(caesarCfg, underlyingServer, nil), nil

		case "trojan":
			trojanCfg, err := parseTrojanServerConfig(cfg)
			if err != nil {
				return nil, err
			}
			underlyingServer := serverMap[cfg.Underlying]
			return trojan.NewServer(ctx, underlyingServer, trojanCfg), nil
		}
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}

	for _, cfg := range dependenciesList {
		server, err := makeServer(cfg)
		if err != nil {
			return nil, err
		}
		if server == nil {
			return nil, fmt.Errorf("failed to make server for %s", cfg.ID)
		}
		serverMap[cfg.ID] = server
	}

	// 返回主transport
	for _, cfg := range cfgs {
		if cfg.Main {
			return serverMap[cfg.ID], nil
		}
	}
	return tcp.NewTCPServer(), fmt.Errorf("no main transport found, using tcp as fallback")
}

func FromConfigClient(ctx context.Context, cfgs []Transport) (transport.TransportClient, error) {
	cfgMap := make(map[string]Transport)
	clientMap := make(map[string]transport.TransportClient)
	for _, cfg := range cfgs {
		cfgMap[cfg.ID] = cfg
	}
	dependenciesList, err := makeDependenciesList(cfgMap, cfgs)
	if err != nil {
		return nil, err
	}

	makeClient := func(cfg Transport) transport.TransportClient {
		switch cfg.Type {
		case "tcp":
			return tcp.NewTCPClient()
		case "tls":
			tlsCfg := &tls.TLSClientConfig{
				ServerName:         cfg.Cfg["ServerName"].(string),
				SNI:                cfg.Cfg["SNI"].(bool),
				InsecureSkipVerify: cfg.Cfg["InsecureSkipVerify"].(bool),
			}
			if sni, ok := cfg.Cfg["SNI"].(bool); ok {
				tlsCfg.SNI = sni
			}
			if insecureSkipVerify, ok := cfg.Cfg["InsecureSkipVerify"].(bool); ok {
				tlsCfg.InsecureSkipVerify = insecureSkipVerify
			}
			return tls.NewTLSClient(tlsCfg)
		case "caesar":
			caesarCfg := &caesar.CaesarConfig{
				Shift: int(cfg.Cfg["Shift"].(int64)),
			}
			underlyingClient := clientMap[cfg.Underlying]
			return caesar.NewCaesarClient(caesarCfg, underlyingClient, nil)
		case "trojan":
			trojanCfg := &trojan.ClientConfig{
				Password: cfg.Cfg["Password"].(string),
			}
			underlyingClient := clientMap[cfg.Underlying]
			return trojan.NewClient(ctx, underlyingClient, trojanCfg)
		}
		return nil
	}

	for _, cfg := range dependenciesList {
		client := makeClient(cfg)
		if client == nil {
			return nil, fmt.Errorf("failed to make client for %s", cfg.ID)
		}
		clientMap[cfg.ID] = client
	}

	// 返回主transport
	for _, cfg := range cfgs {
		if cfg.Main {
			return clientMap[cfg.ID], nil
		}
	}
	return tcp.NewTCPClient(), fmt.Errorf("no main transport found, using tcp as fallback")
}

// ------------------------------------------------------------
// ------------------------------------------------------------

func parseTrojanServerConfig(cfg Transport) (*trojan.ServerConfig, error) {
	trojanCfg := &trojan.ServerConfig{}
	if passwdList, ok := cfg.Cfg["Passwords"]; ok {
		switch v := passwdList.(type) {
		case []string:
			trojanCfg.Passwords = v
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok {
					trojanCfg.Passwords = append(trojanCfg.Passwords, str)
				}
			}
		}
	}
	if len(trojanCfg.Passwords) == 0 {
		return nil, fmt.Errorf("passwords is empty")
	}
	if redirectHost, ok := cfg.Cfg["RedirectHost"].(string); ok {
		trojanCfg.RedirectHost = redirectHost
	} else {
		return nil, fmt.Errorf("redirect host is empty")
	}
	if redirectPort, ok := cfg.Cfg["RedirectPort"].(int64); ok {
		trojanCfg.RedirectPort = int(redirectPort)
	} else {
		return nil, fmt.Errorf("redirect port is empty")
	}
	return trojanCfg, nil
}
