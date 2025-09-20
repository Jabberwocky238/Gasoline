package device

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"wwww/config"
	"wwww/tun"
)

// PeerInfo 对端信息
type PeerInfo struct {
	UniqueID   string
	Ip         net.IP
	AllowedIPs net.IPNet
	Endpoint   net.Addr
}

// Device 设备结构体
type Device struct {
	// 配置信息
	config *config.Config

	// TUN 接口
	tunDevice tun.Device

	// 对端映射表 (公钥 -> PeerInfo)
	indexMap   map[string]*PeerInfo
	indexMutex sync.RWMutex

	// 对端 TLS 连接映射表 (公钥 -> TLS连接)
	connections map[string]*tls.Conn
	connMutex   sync.RWMutex

	// 网络监听
	listener  net.Listener
	tlsConfig *tls.Config

	// 控制通道
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewDevice 创建新的设备实例
func NewDevice(cfg *config.Config) (*Device, error) {
	device := &Device{
		config:      cfg,
		indexMap:    make(map[string]*PeerInfo),
		connections: make(map[string]*tls.Conn),
		stopChan:    make(chan struct{}),
	}

	// 初始化对端映射表
	device.indexMutex.Lock()
	defer device.indexMutex.Unlock()

	for _, peer := range device.config.Peers {
		ip, allowedIPs, err := net.ParseCIDR(peer.AllowedIPs)
		if err != nil {
			return nil, fmt.Errorf("解析允许IPs失败: %v", err)
		}
		endpoint, err := net.ResolveTCPAddr("tcp", peer.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("解析端点失败: %v", err)
		}
		device.indexMap[peer.UniqueID] = &PeerInfo{
			UniqueID:   peer.UniqueID,
			Ip:         ip,
			AllowedIPs: *allowedIPs,
			Endpoint:   endpoint,
		}
	}

	device.tlsConfig = &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	// 初始化 TUN 设备
	if err := device.initializeTUN(); err != nil {
		return nil, fmt.Errorf("初始化 TUN 设备失败: %v", err)
	}

	return device, nil
}

// Start 启动设备
func (d *Device) Start() error {
	// 检查是否为服务器模式（有ListenPort）
	if d.config.Interface.ListenPort > 0 {
		// 启动网络监听
		if err := d.startListener(); err != nil {
			return fmt.Errorf("启动监听失败: %v", err)
		}

		// 启动处理协程
		d.wg.Add(2)
		go d.handleNetworkConnections()
		go d.handleTUNData()
	} else {
		// 客户端模式，只启动TUN数据处理
		d.wg.Add(1)
		go d.handleTUNData()
	}

	return nil
}

// Stop 停止设备
func (d *Device) Stop() error {
	close(d.stopChan)

	// 关闭所有 peer 连接
	d.connMutex.Lock()
	for peerKey, conn := range d.connections {
		conn.Close()
		delete(d.connections, peerKey)
	}
	d.connMutex.Unlock()

	// 关闭监听器
	if d.listener != nil {
		d.listener.Close()
	}

	// 关闭 TUN 设备
	if d.tunDevice != nil {
		d.tunDevice.Close()
	}

	// 等待所有协程结束
	d.wg.Wait()

	return nil
}
