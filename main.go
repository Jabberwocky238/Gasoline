package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wwww/config"
	"wwww/device"
)

func main() {
	// 命令行参数
	var configPath = flag.String("config", "", "配置文件路径")
	var help = flag.Bool("help", false, "显示帮助信息")
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *configPath == "" {
		fmt.Println("错误：必须指定配置文件路径")
		fmt.Println("使用 -help 查看帮助信息")
		os.Exit(1)
	}

	// 解析配置文件
	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 创建设备实例
	dev, err := device.NewDevice(cfg)
	if err != nil {
		log.Fatalf("创建设备失败: %v", err)
	}

	// 启动设备
	if err := dev.Start(); err != nil {
		log.Fatalf("启动设备失败: %v", err)
	}

	// 显示启动信息
	showStartupInfo(cfg)

	// 等待中断信号
	waitForInterrupt()

	// 停止设备
	fmt.Println("\n正在停止设备...")
	if err := dev.Stop(); err != nil {
		log.Printf("停止设备失败: %v", err)
	}

	fmt.Println("设备已停止")
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("WireGuard-like VPN 设备")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Printf("  %s -config <配置文件路径>\n", os.Args[0])
	fmt.Println()
	fmt.Println("参数:")
	fmt.Println("  -config string")
	fmt.Println("        配置文件路径 (必需)")
	fmt.Println("  -help")
	fmt.Println("        显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Printf("  %s -config tests/server.toml\n", os.Args[0])
	fmt.Printf("  %s -config tests/client.toml\n", os.Args[0])
}

// showStartupInfo 显示启动信息
func showStartupInfo(cfg *config.Config) {
	fmt.Println("=== WireGuard-like VPN 设备已启动 ===")
	fmt.Printf("模式: %s\n", getModeString(cfg))
	fmt.Printf("地址: %s\n", cfg.Interface.Address)

	if cfg.Interface.ListenPort > 0 {
		fmt.Printf("监听端口: %d\n", cfg.Interface.ListenPort)
		fmt.Println("状态: 服务器模式 - 等待客户端连接")
	} else {
		fmt.Println("状态: 客户端模式 - 准备连接服务器")
	}

	fmt.Printf("对端数量: %d\n", len(cfg.Peers))
	for i, peer := range cfg.Peers {
		fmt.Printf("  对端 %d: %s (%s)\n", i+1, peer.UniqueID, peer.Endpoint)
	}

	fmt.Println("=====================================")
	fmt.Println("按 Ctrl+C 停止设备")
}

// getModeString 获取模式字符串
func getModeString(cfg *config.Config) string {
	if cfg.Interface.ListenPort > 0 {
		return "服务器"
	}
	return "客户端"
}

// waitForInterrupt 等待中断信号
func waitForInterrupt() {
	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigChan
}
