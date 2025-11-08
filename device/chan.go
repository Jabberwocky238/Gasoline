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
	q := &genericQueue{
		c:  make(chan *PacketBuffer, 1024),
		wg: sync.WaitGroup{},
	}
	q.wg.Add(1)
	go func() {
		q.wg.Wait()
		close(q.c)
	}()
	return q
}
