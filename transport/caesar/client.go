package caesar

import (
	"wwww/transport"
)

type CaesarClient struct {
	cfg      *CaesarConfig
	upstream transport.TransportClient

	debugHook func(bytein, byteout []byte, msg string)
}

func NewCaesarClient(cfg *CaesarConfig, upstream transport.TransportClient, debugHook func(bytein, byteout []byte, msg string)) *CaesarClient {
	return &CaesarClient{
		cfg:       cfg,
		upstream:  upstream,
		debugHook: debugHook,
	}
}

func (t *CaesarClient) Dial(endpoint string) (transport.TransportConn, error) {
	conn, err := t.upstream.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &CaesarConn{
		conn:      conn,
		debugHook: t.debugHook,
		shift:     t.cfg.Shift,
	}, nil
}

func (t *CaesarClient) Close() error {
	return t.upstream.Close()
}
