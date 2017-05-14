package slogger

import "sync"

type _CacheData struct {
	setting *SloggerSettings
	data    *SloggerData
}

type _CacheDatas []*_CacheData

type SloggerProcessorCacheFile struct {
	SloggerProcessorFile
	mutex   sync.Mutex
	buffers *_CacheDatas
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
	self.Offer(&setting, data)
	return nil
}

func (self *SloggerProcessorCacheFile) Shutdown() {
	self.SloggerProcessorFile.Shutdown()
}

func CreateSloggerProcessorCacheFile() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorCacheFile{
		buffers: &_CacheDatas{},
	}
	return &r
}
