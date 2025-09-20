//go:build windows

package device

import (
	"fmt"
	"net"
	"os/exec"
	"time"

	"wwww/tun"
)

// initializeTUN 初始化 TUN 设备
func (d *Device) initializeTUN(tunName string) error {
	// 创建TUN设备
	tunDevice, err := tun.CreateTUN(tunName, 0)
	if err != nil {
		return err
	}
	d.tun = tunDevice

	// 等待设备创建完成
	time.Sleep(50 * time.Millisecond)

	// 获取接口名称
	interfaceName, err := d.tun.Name()
	if err != nil {
		return fmt.Errorf("获取接口名称失败: %v", err)
	}

	// 配置IP地址
	if err := d.configureTUNInterface(interfaceName); err != nil {
		tunDevice.Close()
		return fmt.Errorf("配置TUN接口失败: %v", err)
	}

	// 配置防火墙规则
	if err := d.configureWindowsFirewall(interfaceName); err != nil {
		return fmt.Errorf("配置防火墙失败: %v", err)
	}

	// 添加路由规则，确保VPN网段的流量通过TUN接口
	if err := d.addWindowsRoutes(interfaceName); err != nil {
		fmt.Printf("添加路由规则失败: %v\n", err)
	}

	return nil
}

// configureTUNInterface 配置TUN接口的IP地址
func (d *Device) configureTUNInterface(interfaceName string) error {
	// 配置IP地址
	cmd := exec.Command("netsh", "interface", "ip", "set", "address",
		fmt.Sprintf("name=%s", interfaceName),
		"source=static",
		fmt.Sprintf("addr=%s", d.config.Interface.Address),
		"gateway=none")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("配置IP地址失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// configureWindowsFirewall 配置Windows防火墙规则
func (d *Device) configureWindowsFirewall(interfaceName string) error {
	// 解析接口地址
	_, ipNet, err := net.ParseCIDR(d.config.Interface.Address)
	if err != nil {
		return fmt.Errorf("解析接口地址失败: %v", err)
	}

	networkAddr := ipNet.String()

	// 添加入站规则 - 允许TUN接口的所有入站流量
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s-Inbound", interfaceName),
		"dir=in",
		"action=allow",
		"profile=any",
		fmt.Sprintf("localip=%s", networkAddr),
		"protocol=any")

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("添加入站防火墙规则失败: %v, 输出: %s\n", err, string(output))
	} else {
		fmt.Printf("成功添加入站防火墙规则: %s\n", networkAddr)
	}

	// 添加出站规则 - 允许TUN接口的所有出站流量
	cmd = exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s-Outbound", interfaceName),
		"dir=out",
		"action=allow",
		"profile=any",
		fmt.Sprintf("localip=%s", networkAddr),
		"protocol=any")

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("添加出站防火墙规则失败: %v, 输出: %s\n", err, string(output))
	} else {
		fmt.Printf("成功添加出站防火墙规则: %s\n", networkAddr)
	}

	return nil
}

// addWindowsRoutes 添加Windows路由规则
func (d *Device) addWindowsRoutes(interfaceName string) error {
	// 解析接口地址
	_, ipNet, err := net.ParseCIDR(d.config.Interface.Address)
	if err != nil {
		return fmt.Errorf("解析接口地址失败: %v", err)
	}

	// 获取TUN接口的LUID
	nativeTun, ok := d.tun.(interface{ LUID() uint64 })
	if !ok {
		return fmt.Errorf("无法获取TUN接口LUID")
	}

	luid := nativeTun.LUID()
	if luid == 0 {
		return fmt.Errorf("TUN接口LUID为0")
	}

	networkAddr := ipNet.IP.String()
	mask := ipNet.Mask.String()

	// 添加VPN网段路由，确保流量通过TUN接口
	cmd := exec.Command("route", "add", networkAddr, "mask", mask, "0.0.0.0", "metric", "1", "if", fmt.Sprintf("%d", luid))
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("添加VPN网段路由失败: %v, 输出: %s\n", err, string(output))
	} else {
		fmt.Printf("成功添加VPN网段路由: %s/%s\n", networkAddr, mask)
	}

	return nil
}
