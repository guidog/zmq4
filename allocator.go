package zmq4

import (
	"sync"
)

type slab = []byte

type slabPool = chan *slab

type slabber struct {
	sync.Mutex
	slabs map[int]slabPool
}

const tooLarge = 1024 * 1024

var allocator slabber

func init() {
	allocator.slabs = make(map[int]slabPool)
	for _, x := range []int{64, 256, 512, 1024, 4096, 16384, 65536, 131072} {
		sz := x
		allocator.slabs[sz] = make(slabPool, 2500)
	}
}

func (s slabber) getPoolSize(sz int) (slabPool, int) {
	var (
		p slabPool
		n int
	)
	switch {
	case sz <= 64:
		p = s.slabs[64]
		n = 64
	case sz <= 256:
		p = s.slabs[256]
		n = 256
	case sz <= 512:
		p = s.slabs[512]
		n = 512
	case sz <= 1024:
		p = s.slabs[1024]
		n = 1024
	case sz <= 4096:
		p = s.slabs[4096]
		n = 4096
	case sz <= 16384:
		p = s.slabs[16384]
		n = 16384
	case sz <= 65536:
		p = s.slabs[65536]
		n = 65536
	case sz <= 131072:
		p = s.slabs[131072]
		n = 131072
	default:
		// tooLarge
		p = nil
		n = sz
	}
	return p, n
}

func (s slabber) getSlab(p slabPool, sz int) *slab {
	// s.Lock()
	// defer s.Unlock()
	select {
	case ret := <-p:
		// s.cached[sz] = s.cached[sz] + 1
		return ret
	default:
		ret := make(slab, sz)
		// s.new[sz] = s.new[sz] + 1
		return &ret
	}
}

func (s slabber) alloc(sz int) *slab {
	p, n := s.getPoolSize(sz)
	if p == nil {
		// too large
		ret := make(slab, sz)
		return &ret
	}
	buf := s.getSlab(p, n)
	// set len to given size
	*buf = (*buf)[:sz]
	return buf
}

func (s slabber) free(sl *slab) {
	sz := cap(*sl)
	p, _ := s.getPoolSize(sz)
	if p == nil {
		// too large, do not keep
		return
	}
	p <- sl
}
