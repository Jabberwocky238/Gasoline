// go:build !windows

package device

import (
	"fmt"
	"os/exec"
	"time"
	"wwww/tun"
)

// initializeTUN 初始化 TUN 设备
func (d *Device) initializeTUN(tunName string) error {
	// 这里需要根据实际的 TUN 实现来调用
	// 假设有一个 CreateTUN 函数
	tunDevice, err := tun.CreateTUN(tunName, 1420)
	if err != nil {
		return err
	}
	d.tun = tunDevice

	// 等待设备创建完成
	time.Sleep(50 * time.Millisecond)

	// 配置IP地址
	if err := d.configureTUNInterface(); err != nil {
		tunDevice.Close()
		return fmt.Errorf("配置TUN接口失败: %v", err)
	}

	return nil
}

// configureTUNInterface 配置TUN接口的IP地址
func (d *Device) configureTUNInterface() error {
	// 获取接口名称
	interfaceName, err := d.tun.Name()
	if err != nil {
		return fmt.Errorf("获取接口名称失败: %v", err)
	}

	// 配置IP地址
	cmd := exec.Command("ip", "addr", "add", d.config.Interface.Address, "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("配置IP地址失败: %v", err)
	}

	return nil
}
