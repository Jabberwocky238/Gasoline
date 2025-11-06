package caesar

import (
	"wwww/transport"
)

type CaesarServer struct {
	upstream transport.TransportServer
	shift    int

	debugHook *func(bytein, byteout []byte, msg string)
}

func NewCaesarServer(shift int, upstream transport.TransportServer, debugHook *func(bytein, byteout []byte, msg string)) *CaesarServer {
	return &CaesarServer{
		shift:     shift,
		upstream:  upstream,
		debugHook: debugHook,
	}
}

func (t *CaesarServer) Listen(host string, port int) error {
	return t.upstream.Listen(host, port)
}

func (t *CaesarServer) Accept() (transport.TransportConn, error) {
	conn, err := t.upstream.Accept()
	if err != nil {
		return nil, err
	}
	return &CaesarConn{
		conn:      conn,
		debugHook: *t.debugHook,
		shift:     t.shift,
	}, nil
}

func (t *CaesarServer) Close() error {
	return t.upstream.Close()
}
