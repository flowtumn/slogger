package slogger

import (
	"sync"
	"time"
)

const (
	//wait a 3ms.
	TASK_WAIT_TIME = 3
)

type _CacheData struct {
	setting *SloggerSettings
	data    *SloggerData
}

type _CacheDatas []*_CacheData

type SloggerProcessorCacheFile struct {
	SloggerProcessorFile
	mutex    sync.Mutex
	buffers  *_CacheDatas
	task     *SimpleWork
	shutdown AtomicBool
}

func (self *SloggerProcessorCacheFile) _SafeDo(f func() interface{}) interface{} {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return f()
}

func (self *SloggerProcessorCacheFile) Offer(setting *SloggerSettings, data *SloggerData) {
	self._SafeDo(
		func() interface{} {
			tmp := append(*self.buffers, &_CacheData{setting: setting, data: data})
			self.buffers = &tmp
			return nil
		},
	)
}

func (self *SloggerProcessorCacheFile) Poll() *_CacheDatas {
	if p, ok := self._SafeDo(
		func() interface{} {
			p := self.buffers
			self.buffers = &_CacheDatas{}
			return p
		},
	).(*_CacheDatas); ok {
		return p
	}
	return nil
}

func (self *SloggerProcessorCacheFile) Record(setting SloggerSettings, data *SloggerData) error {
	//Append in queue.
	self.Offer(&setting, data)
	return nil
}

func (self *SloggerProcessorCacheFile) Shutdown() {
	self.shutdown.Set(true)
	for self.task.IsRunning() {
		time.Sleep(TASK_WAIT_TIME >> 1)
	}

	//flush.
	self._Write()

	//File Shutdown.
	self.SloggerProcessorFile.Shutdown()
}

func (self *SloggerProcessorCacheFile) _Write() bool {
	if datas := self.Poll(); nil != datas {
		for _, v := range *datas {
			if err := self.SloggerProcessorFile._UpdateSink(v.setting, v.data.CurrentTimeMillis); nil != err {
				//error.
				return true
			}

			self.logFp.WriteString(v.data.ToLogMessage() + "\n")
		}
	}

	return true && !self.shutdown.Get()
}

func CreateSloggerProcessorCacheFile() *SloggerProcessor {
	r := &SloggerProcessorCacheFile{
		buffers: &_CacheDatas{},
	}

	r.task = CreateWork(
		func() bool {
			return r._Write()
		},
		TASK_WAIT_TIME,
	)

	var rr SloggerProcessor
	rr = r
	return &rr
}
