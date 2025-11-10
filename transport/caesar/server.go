package caesar

import (
	"wwww/transport"
)

type CaesarServer struct {
	upstream transport.TransportServer
	cfg      *CaesarConfig

	debugHook func(bytein, byteout []byte, msg string)

	connChan chan transport.TransportConn
}

func NewCaesarServer(cfg *CaesarConfig, upstream transport.TransportServer, debugHook func(bytein, byteout []byte, msg string)) *CaesarServer {
	return &CaesarServer{
		cfg:       cfg,
		upstream:  upstream,
		debugHook: debugHook,
		connChan:  make(chan transport.TransportConn),
	}
}

func (t *CaesarServer) Listen(host string, port int) error {
	return t.upstream.Listen(host, port)
}

func (t *CaesarServer) Accept() <-chan transport.TransportConn {
	go func() {
		for conn := range t.upstream.Accept() {
			t.connChan <- &CaesarConn{
				conn:      conn,
				debugHook: t.debugHook,
				shift:     t.cfg.Shift,
			}
		}
	}()
	return t.connChan
}

func (t *CaesarServer) Close() error {
	close(t.connChan)
	return t.upstream.Close()
}
