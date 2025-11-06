package caesar

import (
	"wwww/transport"
)

type CaesarClient struct {
	shift    int
	upstream transport.TransportClient

	debugHook *func(bytein, byteout []byte, msg string)
}

func NewCaesarClient(shift int, upstream transport.TransportClient, debugHook *func(bytein, byteout []byte, msg string)) *CaesarClient {
	return &CaesarClient{shift: shift, upstream: upstream, debugHook: debugHook}
}

func (t *CaesarClient) Dial(endpoint string) (transport.TransportConn, error) {
	conn, err := t.upstream.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &CaesarConn{
		conn:      conn,
		debugHook: *t.debugHook,
		shift:     t.shift,
	}, nil
}
