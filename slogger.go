package slogger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type SloggerSettings struct {
	LogLevel     LogLevel
	LogName      string
	LogDirectory string
	LogExtension string
}

type Slogger struct {
	mutex     sync.Mutex
	settings  SloggerSettings
	count     SloggerRecordCount
	task      *_SloggerWorker
	processor *SloggerProcessor
	logPath   string
	logFp     *os.File
}

func (self *Slogger) _SafeDo(f func() interface{}) interface{} {
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
	}()

	return f()
}

func (self *Slogger) Initialize(settings SloggerSettings, processor *SloggerProcessor) error {
	if nil == processor {
		return errors.New("Initialize failed. beacuse processor is nil.")
	}

	self.Close()

	if r, ok := self._SafeDo(
		func() interface{} {
			self.settings = settings
			self.count = SloggerRecordCount{}
			self.processor = processor

			//Create Worker.
			self.task = _CreateSloggerWorker(
				func(buf *SloggerData) {
					if nil == (*self.processor).Record(self.settings, buf) {
						//Record success.
						self.count._CountUpOnLogLevel(buf.logLevel)

						//Update logpath.
						if v := (*self.processor).GetLogPath(); nil != v {
							self.logPath = *v
						}
					}
				},
			)

			return nil
		},
	).(error); ok {
		return r
	}

	return nil
}

func (self *Slogger) Close() {
	self._SafeDo(
		func() interface{} {
			if nil != self.task {
				self.task._Shutdown()
				self.task = nil
			}

			if nil != self.processor {
				(*self.processor).Shutdown()
				self.processor = nil
			}

			return nil
		},
	)
}

func (self *Slogger) Settings() *SloggerSettings {
	if v, ok := self._SafeDo(
		func() interface{} {
			return self.settings
		},
	).(SloggerSettings); ok {
		return &v
	}

	return nil
}

func (self *Slogger) Counters() *SloggerRecordCount {
	if v, ok := self._SafeDo(
		func() interface{} {
			return self.count
		},
	).(SloggerRecordCount); ok {
		return &v
	}

	return nil
}

func (self *Slogger) GetLogPath() *string {
	if v, ok := self._SafeDo(
		func() interface{} {
			return self.logPath
		},
	).(string); ok {
		return &v
	}
	return nil
}

//output log.
func (self *Slogger) record(logLevel LogLevel, format string, v ...interface{}) {
	fileName, fileLine := _GetFileInfoFromStack(3)

	self._SafeDo(
		func() interface{} {
			//filter to loglevel.
			if logLevel < self.settings.LogLevel {
				return nil
			}

			//Enqueue.
			if nil != self.task {
				self.task._Offer(
					&SloggerData{
						logLevel:          logLevel,
						currentTimeMillis: GetCurrentTimeMillis(),
						logMessage:        fmt.Sprintf("%s(%d): ", fileName, fileLine) + fmt.Sprintf(format, v...),
					},
				)
				return nil
			} else {
				return errors.New("task is nil.")
			}
		},
	)
}

func (self *Slogger) Critical(format string, v ...interface{}) {
	self.record(CRITICAL, format, v...)
}

func (self *Slogger) Error(format string, v ...interface{}) {
	self.record(ERROR, format, v...)
}

func (self *Slogger) Warn(format string, v ...interface{}) {
	self.record(WARN, format, v...)
}

func (self *Slogger) Info(format string, v ...interface{}) {
	self.record(INFO, format, v...)
}

func (self *Slogger) Debug(format string, v ...interface{}) {
	self.record(DEBUG, format, v...)
}

func _GetFileInfoFromStack(depth int) (string, int) {
	_, file, line, _ := runtime.Caller(depth)
	_, fileName := filepath.Split(file)
	return fileName, line
}
