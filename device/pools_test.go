package device

import (
	"bytes"
	"testing"
	"unsafe"
)

// go test -v ./device -run TestPacketBufferSize
func TestPacketBufferSize(t *testing.T) {
	var pb PacketBuffer
	size := unsafe.Sizeof(pb)
	t.Logf("PacketBuffer 结构体大小: %d 字节", size)

	// 详细分析各个字段的大小
	bufferSize := unsafe.Sizeof(pb.buffer)
	lengthSize := unsafe.Sizeof(pb.length)
	t.Logf("  - buffer [1600]byte: %d 字节", bufferSize)
	t.Logf("  - length int: %d 字节", lengthSize)
	t.Logf("  总计: %d 字节 (期望: 1600 + 8 = 1608 字节)", bufferSize+lengthSize)

	// 验证大小是否符合预期
	expectedSize := 1600 + 8 // 1600字节数组 + 8字节int
	if size != uintptr(expectedSize) {
		t.Errorf("结构体大小不符合预期: 期望 %d 字节, 实际 %d 字节", expectedSize, size)
	}
}

// go test -v ./device -run TestPacketBufferCopy
func TestPacketBufferCopy(t *testing.T) {
	pools := NewPool()
	pb := pools.GetPacketBuffer()
	msg := []byte("Hello, world!")
	pb.Set(msg)
	copyBuf := pb.CopyPacket()
	pools.PutPacketBuffer(pb)
	if !bytes.Equal(copyBuf, msg) {
		t.Errorf("copyBuf 不符合预期: 期望 %s, 实际 %s", string(msg), string(copyBuf))
	}
}
