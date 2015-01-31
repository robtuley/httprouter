package proxy

import (
	"container/ring"
	"net/http"
	"sync"
)

func newRoundRobin(fallback http.Handler) roundRobin {
	return roundRobin{fallback: fallback}
}

type roundRobin struct {
	mu       sync.RWMutex
	ring     *ring.Ring
	fallback http.Handler
}

func (rr *roundRobin) Choose() http.Handler {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if rr.ring == nil {
		return rr.fallback
	}

	h := rr.ring.Value.(http.Handler)
	rr.ring = rr.ring.Next()
	return h
}

func (rr *roundRobin) Add(h http.Handler) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	nr := &ring.Ring{Value: h}
	if rr.ring == nil {
		rr.ring = nr
	} else {
		rr.ring = rr.ring.Link(nr).Next()
	}
}

func (rr *roundRobin) Remove(h http.Handler) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if rr.ring == nil {
		return
	}

	if rr.ring.Len() == 1 {
		if h == rr.ring.Value {
			rr.ring = ring.New(0)
		}
		return
	}

	for i := rr.ring.Len(); i > 0; i-- {
		rr.ring = rr.ring.Next()
		if h == rr.ring.Value {
			rr.ring = rr.ring.Unlink(1)
			return
		}
	}
}
