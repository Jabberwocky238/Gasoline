package device

import (
	"context"
	"sync"
	"wwww/config"
	"wwww/transport"
)

type DeviceListener struct {
	mu     sync.Mutex
	server transport.TransportServer
}

func NewDeviceListener(ctx context.Context, cfg *config.Config) (*DeviceListener, error) {
	listener := new(DeviceListener)
	var err error
	listener.server, err = config.FromConfigServer(ctx, cfg.Transports, cfg.Interface.TransportID)
	if err != nil {
		if listener.server == nil {
			return nil, err
		}
		log.Warnf("Failed to create server: %v", err)
	}
	return listener, nil
}

func (d *DeviceListener) Accept() <-chan transport.TransportConn {
	connChan := make(chan transport.TransportConn)
	go func() {
		for {
			conn, err := d.server.Accept()
			if err != nil {
				return
			}
			connChan <- conn
		}
	}()
	return connChan
}

func (d *DeviceListener) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.server.Close()
}

func (d *DeviceListener) Listen(host string, port int) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.server.Listen(host, port)
}
