package slogger

import "sync"

type AtomicBool struct {
	mutex sync.Mutex
	value bool
}

func (p *AtomicBool) _SafeDo(f func() bool) bool {
	var r bool
	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
	}()
	r = f()
	return r
}

func (p *AtomicBool) Set(v bool) {
	p._SafeDo(
		func() bool {
			p.value = v
			return p.value
		},
	)
}

func (p *AtomicBool) Get() bool {
	return p._SafeDo(
		func() bool {
			return p.value
		},
	)
}
