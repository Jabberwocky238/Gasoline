package device

import (
	"net"
)

const (
	IPv4offsetTotalLength = 2
	IPv4offsetSrc         = 12
	IPv4offsetDst         = IPv4offsetSrc + net.IPv4len
)

const (
	IPv6offsetPayloadLength = 4
	IPv6offsetSrc           = 8
	IPv6offsetDst           = IPv6offsetSrc + net.IPv6len
)

func (device *Device) RoutineReadFromTUN() {
	defer func() {
		device.log.Debugf("Routine: TUN reader - stopped")
	}()

	device.log.Debugf("Routine: TUN reader - started")

	for {
		// read packets
		pb := device.pools.GetPacketBuffer()
		length, readErr := device.tun.Read(pb.packet)
		if readErr != nil {
			if device.ctx.Err() != nil {
				break
			}
			device.log.Errorf("Failed to read packet from TUN device: %v", readErr)
			break
		}
		if length < 1 {
			device.log.Debugf("Received packet with length 0 from TUN device")
			continue
		}
		pb.length = length
		pb.Make()
		device.queue.routing.c <- pb
	}
}

func (device *Device) RoutineWriteToTUN() {
	defer func() {
		device.log.Debugf("Routine: TUN writer - stopped")
	}()

	device.log.Debugf("Routine: TUN writer - started")

	for pb := range device.queue.outbound.c {
		_, err := device.tun.Write(pb.CopyPacket())
		device.pools.PutPacketBuffer(pb)
		if err != nil {
			device.log.Errorf("Failed to write packet to TUN device: %v", err)
			break
		}
	}
}
