package device

import (
	"sync"
)

type genericQueueElement struct {
	packet []byte
}

type genericQueueElementsContainer struct {
	// sync.Mutex
	elems []*genericQueueElement
}

type genericQueue struct {
	c  chan *PacketBuffer
	wg sync.WaitGroup
}

func newGenericQueue() *genericQueue {
	// 增加队列容量到4096，以应对高延迟场景
	// 高延迟下，TCP发送可能阻塞，需要更大的队列缓冲
	q := &genericQueue{
		c:  make(chan *PacketBuffer, 4096),
		wg: sync.WaitGroup{},
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()
	return q
}
