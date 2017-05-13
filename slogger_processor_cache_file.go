package slogger

import "sync"

type _CacheData struct {
	setting *SloggerSettings
	data    *SloggerData
}

type SloggerProcessorCacheFile struct {
	SloggerProcessorFile
	mutex   sync.Mutex
	buffers *[]*_CacheData
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

func (self *SloggerProcessorCacheFile) Poll() *[]*_CacheData {
	if p, ok := self._SafeDo(
		func() interface{} {
			p := self.buffers
			self.buffers = &[]*_CacheData{}
			return p
		},
	).(*[]*_CacheData); ok {
		return p
	}
	return nil
}

func (self *SloggerProcessorCacheFile) Record(setting SloggerSettings, data *SloggerData) error {
	defer func() {

	}()
	return nil
	if err := self._UpdateSink(&setting, data.CurrentTimeMillis); nil != err {
		return err
	}

	//Log write.
	if _, err := self.logFp.WriteString(data.ToLogMessage() + "\n"); nil != err {
		return err
	}

	return nil
}

func CreateSloggerProcessorCacheFile() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorCacheFile{
		buffers: &[]*_CacheData{},
	}
	return &r
}
