package device

import "encoding/binary"

const (
	TransportMsgHeaderSize = 2
	TransportMsgMaxSize    = 1 << 22 // 4MB
)

type TransportMsg struct {
	size   uint16 // maxsize: 2^20 = 1MB
	packet []byte
}

func (msg *TransportMsg) Marshal(buf []byte) {
	binary.LittleEndian.PutUint16(buf[:TransportMsgHeaderSize], msg.size)
	copy(buf[TransportMsgHeaderSize:], msg.packet)
}

func (msg *TransportMsg) Unmarshal(buf []byte) {
	msg.size = binary.LittleEndian.Uint16(buf[:TransportMsgHeaderSize])
	msg.packet = buf[TransportMsgHeaderSize:]
}
