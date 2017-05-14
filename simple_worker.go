package slogger

import (
	"sync"
	"time"
)

type SimpleWork struct {
	wg         sync.WaitGroup
	_IsRunning AtomicBool
}

func (self *SimpleWork) IsRunning() bool {
	return self._IsRunning.Get()
}

func _DoProcess(f func() bool, wait func(), exit func()) {
	for f() {
		wait()
	}
	exit()
}

func CreateWork(f func() bool, cycleMillis int64) *SimpleWork {
	r := &SimpleWork{}

	r.wg.Add(1)
	r._IsRunning.Set(true)

	go _DoProcess(
		f,
		func() {
			time.Sleep((time.Duration)(cycleMillis) * time.Millisecond)
		},
		func() {
			r.wg.Done()
			r._IsRunning.Set(false)
		},
	)

	return r
}
