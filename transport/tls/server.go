package tls

import (
	"crypto/tls"
	"fmt"
	"net"

	"wwww/transport"
)

type TLSServer struct {
	cfg      *TLSServerConfig
	listener net.Listener

	connChan chan transport.TransportConn
}

func NewTLSServer(cfg *TLSServerConfig) *TLSServer {
	return &TLSServer{
		cfg:      cfg,
		connChan: make(chan transport.TransportConn, 1024),
	}
}

func (t *TLSServer) Listen(host string, port int) error {
	baseListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	tlsCfg := t.cfg.TLSConfig
	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
		if len(t.cfg.CertPEM) > 0 && len(t.cfg.KeyPEM) > 0 {
			cert, err := tls.X509KeyPair(t.cfg.CertPEM, t.cfg.KeyPEM)
			if err != nil {
				baseListener.Close()
				return err
			}
			tlsCfg.Certificates = []tls.Certificate{cert}
		}
		if t.cfg.ClientCAs != nil {
			tlsCfg.ClientCAs = t.cfg.ClientCAs
			if t.cfg.RequireClientCert {
				tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
			}
		}
	}

	t.listener = tls.NewListener(baseListener, tlsCfg)
	go t.acceptLoop()
	return nil
}

func (t *TLSServer) acceptLoop() (net.Conn, error) {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return nil, err
		}
		t.connChan <- conn
	}
}

func (t *TLSServer) Accept() <-chan transport.TransportConn {
	return t.connChan
}

func (t *TLSServer) Close() error {
	close(t.connChan)
	if t.listener != nil {
		return t.listener.Close()
	}
	return nil
}
