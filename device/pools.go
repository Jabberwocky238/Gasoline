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
	length int
}

func (p *PacketBuffer) SetPacket(buf []byte) {
	copy(p.buffer[MessageHeaderSize:], buf)
	p.length = len(buf)
	binary.LittleEndian.PutUint16(p.buffer[:MessageHeaderSize], uint16(p.length))
}

func (p *PacketBuffer) Packet() []byte {
	return p.buffer[MessageHeaderSize : MessageHeaderSize+p.length]
}

func (p *PacketBuffer) Message() []byte {
	return p.buffer[:MessageHeaderSize+p.length]
}

func (p *PacketBuffer) IpVersion() int {
	return int(p.Packet()[0] >> 4)
}

func NewPool() *Pools {
	return &Pools{
		PacketBuffers: NewWaitPool(0, func() any {
			return &PacketBuffer{
				length: 0,
			}
		}),
	}
}

func (p *Pools) GetPacketBuffer() *PacketBuffer {
	pb := p.PacketBuffers.Get().(*PacketBuffer)
	pb.length = 0
	return pb
}

func (p *Pools) PutPacketBuffer(pb *PacketBuffer) {
	pb.length = 0
	p.PacketBuffers.Put(pb)
}
