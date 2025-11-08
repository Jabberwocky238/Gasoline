package tcp

import (
	"fmt"
	"net"

	"wwww/transport"
)

type TCPServer struct {
	listener *net.TCPListener

	connChan chan transport.TransportConn
}

func NewTCPServer() transport.TransportServer {
	return &TCPServer{
		connChan: make(chan transport.TransportConn, 1024),
	}
}

func (t *TCPServer) Listen(host string, port int) error {
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	t.listener = listener.(*net.TCPListener)
	go t.acceptLoop()
	return nil
}

func (t *TCPServer) acceptLoop() (net.Conn, error) {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return nil, err
		}
		// 优化TCP连接性能
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)                // 关闭Nagle算法，减少延迟
			tcpConn.SetReadBuffer(4 * 1024 * 1024)  // 4MB读缓冲区
			tcpConn.SetWriteBuffer(4 * 1024 * 1024) // 4MB写缓冲区
		}
		t.connChan <- conn.(*net.TCPConn)
	}
}

func (t *TCPServer) Accept() <-chan transport.TransportConn {
	return t.connChan
}

func (t *TCPServer) Close() error {
	close(t.connChan)
	return t.listener.Close()
}
