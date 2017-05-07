package slogger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Slogger struct {
	mutex     sync.Mutex
	settings  SloggerSettings
	count     SloggerRecordCount
	task      *_SloggerWorker
	processor *SloggerProcessor
	logPath   string
	logFp     *os.File
}

func CreateSlogger(settings SloggerSettings, processor *SloggerProcessor) (*Slogger, error) {
	if nil == processor {
		return nil, errors.New("Initialize failed. beacuse processor is nil.")
	}

	settings.Trim()
	p := &Slogger{
		settings:  settings,
		processor: processor,
	}

	//Create Worker.
	p.task = _CreateSloggerWorker(
		func(data *SloggerData) {
			if nil == (*p.processor).Record(p.settings, data) {
				//Record success.
				p.count._CountUpOnLogLevel(data.LogLevel)

				//Update logpath.
				if v := (*p.processor).GetLogPath(); nil != v {
					p.logPath = *v
				}
			}
		},
	)

	if nil == p.task {
		return nil, errors.New("Internal Error. Failed to create worker.")
	}

	return p, nil
}

func (self *Slogger) _SafeDo(f func() interface{}) interface{} {
	self.mutex.Lock()
	defer func() {
		self.mutex.Unlock()
	}()

	return f()
}

func (self *Slogger) IsRunning() bool {
	return self.task.IsRunning()
}

func (self *Slogger) Close() {
	self._SafeDo(
		func() interface{} {
			if self.task.IsRunning() {
				self.task._Shutdown()
				(*self.processor).Shutdown()
			}
			return nil
		},
	)
}

func (self *Slogger) Settings() SloggerSettings {
	return self.settings
}

func (self *Slogger) RecordCounter() *SloggerRecordCount {
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
	//filter to loglevel.
	if logLevel < self.settings.LogLevel {
		return
	}

	fileName, fileLine := _GetFileInfoFromStack(3)

	self.task._Offer(
		&SloggerData{
			CurrentTimeMillis: GetCurrentTimeMillis(),
			LogLevel:          logLevel,
			LogMessage:        fmt.Sprintf(format, v...),
			SourceName:        fileName,
			SourceLine:        fileLine,
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
