package device

import (
	"crypto/tls"
	"fmt"
)

// HandleOutbound 处理从TUN设备读取的数据包（出站流量）
func (d *Device) HandleOutbound(data []byte) {
	// 解析IP头部获取连接信息
	ipHeader, err := ParseIPHeader(data)
	if err != nil {
		// 静默跳过无法解析的包（可能是IPv6或其他格式）
		return
	}

	// 首先检查目标IP是否在peers的AllowedIPs范围内
	if !d.IsIPInAllowedRange(ipHeader.DestIP) {
		return
	}

	// 检查是否为广播或多播包，如果是则跳过
	if IsBroadcastOrMulticast(ipHeader.DestIP) {
		// 静默跳过广播/多播包
		return
	}

	// 创建或更新连接元数据
	metadata := d.connectionManager.GetOrCreateConnection(
		ipHeader.SourceIP,
		ipHeader.DestIP,
		ipHeader.Protocol,
	)

	// 设置端口信息
	metadata.SourcePort = ipHeader.GetSourcePort(data)
	metadata.DestPort = ipHeader.GetDestPort(data)

	// 更新统计信息
	metadata.UpdateStats(0, uint64(len(data)))

	// 根据目标IP查找对应的对端连接
	targetConn := d.findPeerByIP(ipHeader.DestIP)
	if targetConn == nil {
		return
	}

	// 发送到目标对端
	d.sendToPeer(targetConn, data)
}

// sendToPeer 发送数据到指定的对端连接
func (d *Device) sendToPeer(targetConn *tls.Conn, data []byte) error {
	// 创建协议消息
	msg := NewProtocolMessage(Data, nil, nil, data)

	// 序列化协议消息
	serializedData := msg.Serialize()

	// 发送到目标对端
	_, err := targetConn.Write(serializedData)
	if err != nil {
		return fmt.Errorf("发送数据到对端失败: %v", err)
	}

	return nil
}
