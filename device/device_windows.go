// go:build windows
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
	cmd := exec.Command("netsh", "interface", "ip", "set", "address",
		fmt.Sprintf("name=%s", interfaceName),
		fmt.Sprintf("source=static"),
		fmt.Sprintf("addr=%s", d.config.Interface.Address),
		fmt.Sprintf("gateway=none"))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("配置IP地址失败: %v, 输出: %s", err, string(output))
	}

	// 动态添加路由规则
	// 尝试获取LUID（Windows特定）
	if nativeTun, ok := d.tun.(interface{ LUID() uint64 }); ok {
		luid := nativeTun.LUID()
		if luid > 0 {
			// 解析配置中的地址
			_, ipNet, err := net.ParseCIDR(d.config.Interface.Address)
			if err == nil {
				networkAddr := ipNet.IP.String()
				mask := ipNet.Mask.String()

				// 添加路由规则
				cmd := exec.Command("route", "add", networkAddr, "mask", mask, "0.0.0.0", "if", fmt.Sprintf("%d", luid))
				cmd.CombinedOutput()
			}
		}
	}

	return nil
}
