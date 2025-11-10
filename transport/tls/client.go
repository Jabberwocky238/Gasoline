package tls

import (
	"crypto/tls"
	"wwww/transport"
)

type TLSClient struct {
	cfg           *TLSClientConfig
	transportConn *TransportTLSConn
}

func NewTLSClient(cfg *TLSClientConfig) transport.TransportClient {
	return &TLSClient{
		cfg:           cfg,
		transportConn: nil,
	}
}

func (t *TLSClient) Dial(endpoint string) (transport.TransportConn, error) {
	tlsCfg := t.cfg.TLSConfig
	if tlsCfg == nil {
		tlsCfg = &tls.Config{
			InsecureSkipVerify: t.cfg.InsecureSkipVerify,
			ServerName:         t.cfg.ServerName,
			RootCAs:            t.cfg.RootCAs,
			Certificates:       t.cfg.Certificates,
		}
	}

	conn, err := tls.Dial("tcp", endpoint, tlsCfg)
	if err != nil {
		return nil, err
	}
	t.transportConn = conn
	return t.transportConn, nil
}

func (t *TLSClient) Close() error {
	if t.transportConn != nil {
		return t.transportConn.Close()
	}
	return nil
}
