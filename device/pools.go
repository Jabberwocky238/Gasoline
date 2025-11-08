package device

import (
	"encoding/binary"
	"sync"
)

type WaitPool struct {
	pool  sync.Pool
	cond  sync.Cond
	lock  sync.Mutex
	count uint32 // Get calls not yet Put back
	max   uint32
}

func NewWaitPool(max uint32, new func() any) *WaitPool {
	p := &WaitPool{pool: sync.Pool{New: new}, max: max}
	p.cond = sync.Cond{L: &p.lock}
	return p
}

func (p *WaitPool) Get() any {
	if p.max != 0 {
		p.lock.Lock()
		for p.count >= p.max {
			p.cond.Wait()
		}
		p.count++
		p.lock.Unlock()
	}
	return p.pool.Get()
}

func (p *WaitPool) Put(x any) {
	p.pool.Put(x)
	if p.max == 0 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.count--
	p.cond.Signal()
}

type Pools struct {
	PacketBuffers *WaitPool
}

const (
	MaxMtu            = 1600
	MessageHeaderSize = 2
)

type PacketBuffer struct {
	buffer [MessageHeaderSize + MaxMtu]byte // transport message bufferF

	packet    []byte // raw packet
	message   []byte // Encoded transport message
	length    int
	ipVersion int
}

func (p *PacketBuffer) Set(buf []byte) {
	p.length = len(buf)
	binary.LittleEndian.PutUint16(p.buffer[:MessageHeaderSize], uint16(p.length))
	copy(p.buffer[MessageHeaderSize:], buf)
	p.packet = p.buffer[MessageHeaderSize : MessageHeaderSize+p.length]
	p.message = p.buffer[:MessageHeaderSize+p.length]
	p.ipVersion = int(p.packet[0] >> 4)
}

func (p *PacketBuffer) CopyPacket() []byte {
	copyBuf := make([]byte, p.length)
	copy(copyBuf, p.packet)
	return copyBuf
}

func (p *PacketBuffer) CopyMessage() []byte {
	copyBuf := make([]byte, p.length+MessageHeaderSize)
	copy(copyBuf, p.message)
	return copyBuf
}

func NewPool() *Pools {
	return &Pools{
		PacketBuffers: NewWaitPool(0, func() any {
			return &PacketBuffer{
				length: 0,
				packet: nil,
			}
		}),
	}
}

func (p *Pools) GetPacketBuffer() *PacketBuffer {
	return p.PacketBuffers.Get().(*PacketBuffer)
}

func (p *Pools) PutPacketBuffer(pb *PacketBuffer) {
	pb.length = 0
	pb.packet = nil
	p.PacketBuffers.Put(pb)
}
