package udp

import (
	"net"

	"wwww/transport"
)

type UDPServer struct {
	bind     *net.UDPConn
	connChan chan transport.TransportConn
	conns    map[string]*UDPConn // 使用地址字符串作为键
}

func NewUDPServer() *UDPServer {
	return &UDPServer{
		connChan: make(chan transport.TransportConn, 1024),
		conns:    make(map[string]*UDPConn),
	}
}

func (t *UDPServer) Listen(host string, port int) error {
	addr := net.UDPAddr{IP: net.ParseIP(host), Port: port}
	listener, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	t.bind = listener
	go t.acceptLoop(&addr)
	return nil
}

func (t *UDPServer) acceptLoop(laddr *net.UDPAddr) {
	buf := make([]byte, 65535)
	for {
		n, raddr, err := t.bind.ReadFromUDP(buf)
		if err != nil {
			break
		}
		// 使用地址字符串作为键，确保相同地址的客户端使用同一个连接
		addrKey := raddr.String()
		if _, ok := t.conns[addrKey]; !ok {
			packetConn := net.PacketConn(t.bind)
			conn := NewUDPConn(laddr, raddr, packetConn)
			t.conns[addrKey] = conn
			t.connChan <- conn
		}
		// 复制数据到新的 slice，避免被下次读取覆盖
		data := make([]byte, n)
		copy(data, buf[:n])
		t.conns[addrKey].packetChan <- data
	}
}

func (t *UDPServer) Accept() <-chan transport.TransportConn {
	return t.connChan
}

func (t *UDPServer) Close() error {
	close(t.connChan)
	return t.bind.Close()
}
