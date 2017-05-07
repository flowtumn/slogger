package slogger

import "sync"

type _SloggerWorker struct {
	wg      sync.WaitGroup
	queue   chan *SloggerData
	handler func(*SloggerData)
	running AtomicBool
}

func _CreateSloggerWorker(handler func(*SloggerData)) *_SloggerWorker {
	p := &_SloggerWorker{
		queue:   make(chan *SloggerData, 1),
		handler: handler,
	}
	p._DoWork()
	return p
}

func (self *_SloggerWorker) _DoWork() {
	go self._Process()
	self.wg.Add(1)
	self.running.Set(true)
}

func (self *_SloggerWorker) _Offer(v *SloggerData) bool {
	if !self.running.Get() {
		return false
	}

	self.queue <- v
	return true
}

func (self *_SloggerWorker) IsRunning() bool {
	return self.running.Get()
}

func (self *_SloggerWorker) _Process() {
	defer func() {
		self.wg.Done()
	}()

	for {
		select {
		case v := <-self.queue:
			if nil == v {
				//Exit.
				return
			}
			self.handler(v)
		}
	}
}

func (self *_SloggerWorker) _Shutdown() {
	//End mark.
	self.queue <- nil
	self.wg.Wait()
	self.running.Set(false)
}
