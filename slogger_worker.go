package slogger

import "sync"

type _SloggerWorker struct {
	wg      sync.WaitGroup
	queue   chan *_SloggerData
	handler func(*_SloggerData)
	running AtomicBool
}

func _CreateSloggerWorker(handler func(*_SloggerData)) *_SloggerWorker {
	p := &_SloggerWorker{
		queue:   make(chan *_SloggerData, 1),
		handler: handler,
	}
	p._DoWork()
	return p
}

func (p *_SloggerWorker) _DoWork() {
	go p._Process()
	p.wg.Add(1)
	p.running.Set(true)
}

func (p *_SloggerWorker) _Offer(v *_SloggerData) bool {
	if !p.running.Get() {
		return false
	}

	p.queue <- v
	return true
}

func (p *_SloggerWorker) IsRunning() bool {
	return p.running.Get()
}

func (p *_SloggerWorker) _Process() {
	defer func() {
		p.wg.Done()
	}()

	for {
		select {
		case v := <-p.queue:
			if nil == v {
				//Exit.
				return
			}
			p.handler(v)
		}
	}
}

func (p *_SloggerWorker) _Shutdown() {
	//End mark.
	p.queue <- nil
	p.wg.Wait()
	p.running.Set(false)
}
