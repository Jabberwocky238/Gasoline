package tls

import (
	"crypto/tls"
	"fmt"
	"net"

	"wwww/transport"
)

type TLSServer struct {
	cfg    *TLSServerConfig
	tlsCfg *tls.Config

	listener net.Listener
}

func NewTLSServer(cfg *TLSServerConfig) *TLSServer {
	tlsCfg, err := cfg.ToTlsConfig()
	if err != nil {
		return nil
	}
	return &TLSServer{
		cfg:    cfg,
		tlsCfg: tlsCfg,
	}
}

func (t *TLSServer) Listen(host string, port int) error {
	baseListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	t.listener = tls.NewListener(baseListener, t.tlsCfg)
	return nil
}

func (t *TLSServer) Accept() (transport.TransportConn, error) {
	conn, err := t.listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (t *TLSServer) Close() error {
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
