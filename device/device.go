package device

import (
	"bytes"
	"context"
	"net"
	"net/netip"
	"time"

	// "time"
	"wwww/config"

	singTun "github.com/jabberwocky238/sing-tun"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type Device struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	tun    singTun.Tun

	key struct {
		privateKey PrivateKey
		publicKey  PublicKey
	}

	peers      map[PublicKey]*Peer
	allowedips AllowedIPs

	listener *DeviceListener
	endpoint struct {
		local net.IP
	}

	queue struct {
		outbound *genericQueue // 进入TUN的包
		routing  *genericQueue // 需要路由的包
	}

	pools    *Pools // 内存池
	debugger *Debugger
}

func NewDevice(cfg *config.Config, tun singTun.Tun) *Device {
	var err error
	device := new(Device)
	device.ctx, device.cancel = context.WithCancel(context.Background())

	device.debugger = NewDebugger(device)
	device.cfg = cfg
	device.tun = tun
	device.pools = NewPool()
	if device.cfg.Interface.ListenPort > 0 {
		device.listener, err = NewDeviceListener(device.ctx, device.cfg)
		if err != nil {
			log.Errorf("Failed to create listener: %v", err)
			return nil
		}
	}

	var privateKey PrivateKey
	if err := privateKey.FromBase64(device.cfg.Interface.PrivateKey); err != nil {
		log.Errorf("Failed to parse private key: %v", err)
		return nil
	}
	device.key.privateKey = privateKey
	device.key.publicKey = privateKey.PublicKey()

	prefix, err := netip.ParsePrefix(device.cfg.Interface.Address)
	if err != nil {
		log.Errorf("Failed to parse interface address: %v", err)
		return nil
	}
	localIp := prefix.Addr().AsSlice()
	device.endpoint.local = localIp
	log.Debugf("Device local IP %s", net.IP(localIp).String())

	device.peers = make(map[PublicKey]*Peer)
	for _, peerConfig := range device.cfg.Peers {
		peer, err := device.NewPeer(&peerConfig)
		if err != nil {
			log.Errorf("Failed to create peer: %v", err)
			continue
		}
		device.peers[peer.key.publicKey] = peer
		device.allowedips.Insert(peer.allowedIPs, peer)
	}

	return device
}

func (device *Device) Start() error {
	log.Debugf("Starting device")
	err := device.tun.Start()
	if err != nil {
		log.Errorf("Failed to start tun: %v", err)
		return err
	}

	go device.debugger.Start() // lifetime listen to device
	device.queue.outbound = newGenericQueue()
	device.queue.routing = newGenericQueue()

	for _, peer := range device.peers {
		err := peer.Start()
		if err != nil {
			log.Errorf("Failed to start peer %s: %v", peer.endpoint.local.String(), err)
			continue
		}
	}

	if device.cfg.Interface.ListenPort > 0 {
		go device.RoutineListenPort() // lifetime listen to device.listener
	}

	go device.RoutineRoutingPackets() // lifetime listen to device.queue.routing
	go device.RoutineReadFromTUN(0)   // lifetime listen to device.tun
	go device.RoutineWriteToTUN()     // lifetime listen to device.tun and device.queue.outbound

	// go device.RoutineBoardcast()

	return nil
}

func (device *Device) Close() {
	log.Debugf("Closing device")
	device.queue.outbound.wg.Done()
	device.queue.routing.wg.Done()
	device.listener.Close()
	device.cancel()
	device.tun.Close()
}

func (device *Device) RoutineListenPort() error {
	defer func() {
		log.Debugf("Routine: listen port - stopped")
	}()
	log.Debugf("Routine: listen port - started")

	host := "0.0.0.0"
	port := device.cfg.Interface.ListenPort

	err := device.listener.Listen(host, port)
	if err != nil {
		log.Errorf("Failed to listen on port %d: %v", port, err)
		return err
	}

	for conn := range device.listener.Accept() {
		log.Debugf("Accepted connection from %s", conn.RemoteAddr().String())
		go func() {
			handshake := NewHandshake(conn, device, nil)
			publicKey, err := handshake.ReceiveHandshake()
			if err != nil {
				log.Errorf("Failed to receive handshake: %v", err)
				return
			}
			peer := device.peers[*publicKey]
			if peer == nil {
				log.Errorf("Peer not found for public key %s", publicKey)
				return
			}
			peer.conn.mu.Lock()
			peer.conn.conn = conn
			peer.conn.handshake = handshake
			peer.conn.isConnected = true
			peer.conn.mu.Unlock()

			log.Debugf("Connected to peer %s", peer.endpoint.local.String())
			go peer.RoutineSequentialSender()
			go peer.RoutineSequentialReceiver()
		}()
	}
	return nil
}

func (device *Device) RoutineRoutingPackets() {
	defer func() {
		log.Debugf("Routine: routing packets - stopped")
	}()

	log.Debugf("Routine: routing packets - started")

	var (
		ticker               = time.NewTicker(1 * time.Second)
		correctPacketCount   = 0
		incorrectPacketCount = 0
	)

	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Debugf("Routine: routing packets - correct: %d, incorrect: %d", correctPacketCount, incorrectPacketCount)
		}
	}()

	for pb := range device.queue.routing.c {
		// lookup peer
		var peer *Peer
		var dst []byte
		switch pb.IpVersion() {
		case 4:
			if pb.length < ipv4.HeaderLen {
				incorrectPacketCount++
				device.pools.PutPacketBuffer(pb)
				continue
			}
			// device.debugger.LogPacket(pb.CopyPacket(), 4)
			dst = pb.Packet()[IPv4offsetDst : IPv4offsetDst+net.IPv4len]
		case 6:
			if pb.length < ipv6.HeaderLen {
				incorrectPacketCount++
				device.pools.PutPacketBuffer(pb)
				continue
			}
			// device.debugger.LogPacket(pb.CopyPacket(), 6)
			dst = pb.Packet()[IPv6offsetDst : IPv6offsetDst+net.IPv6len]
		default:
			log.Debugf("Received packet with unknown IP version")
			device.pools.PutPacketBuffer(pb)
			incorrectPacketCount++
			continue
		}

		correctPacketCount++
		// 判断接收者是不是自己
		if bytes.Equal(dst, device.endpoint.local) {
			device.queue.outbound.c <- pb
			continue
		}
		// 查找peer
		peer = device.allowedips.Lookup(dst)
		if peer == nil {
			// device.log.Errorf("Peer not found for IP %s", net.IP(dst).String())
			device.pools.PutPacketBuffer(pb)
			incorrectPacketCount++
			continue
		}
		peer.queue.inbound.c <- pb
	}
}
