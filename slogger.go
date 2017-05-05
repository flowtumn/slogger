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
	mutex            sync.Mutex
	settings         SloggerSettings
	count            SloggerRecordCount
	task             *_SloggerWorker
	processor        *SloggerProcessor
	isRuning         bool
	currentTimeStamp string
	logFp            *os.File
}

func _CreateLogFileName(prefix string, suffix string) string {
	return prefix + "-" + GetTimeStamp(Normal) + "." + suffix
}

func (p *Slogger) _SafeDo(f func() interface{}) interface{} {
	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
	}()

	return f()
}

func (p *Slogger) Initialize(settings SloggerSettings, processor *SloggerProcessor) error {
	if nil == processor {
		return errors.New("Initialize failed. beacuse processor is nil.")
	}

	p.Close()

	if r, ok := p._SafeDo(
		func() interface{} {
			p.settings = settings
			p.count = SloggerRecordCount{}
			p.isRuning = false
			p.processor = processor

			//Create Worker.
			p.task = _CreateSloggerWorker(
				func(buf *SloggerData) {
					if nil == (*p.processor).Record(buf) {
						//Record success.
						p.count._CountUpOnLogLevel(buf.logLevel)
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

func (p *Slogger) Close() {
	p._SafeDo(
		func() interface{} {
			if nil != p.task {
				p.task._Shutdown()
				p.task = nil
			}

			if nil != p.processor {
				(*p.processor).Shutdown()
				p.processor = nil
			}

			return nil
		},
	)
}

func (p *Slogger) Settings() *SloggerSettings {
	if v, ok := p._SafeDo(
		func() interface{} {
			return p.settings
		},
	).(SloggerSettings); ok {
		return &v
	}

	return nil
}

func (p *Slogger) Counters() *SloggerRecordCount {
	if v, ok := p._SafeDo(
		func() interface{} {
			return p.count
		},
	).(SloggerRecordCount); ok {
		return &v
	}

	return nil
}

func (p *Slogger) GetLogPath() *string {
	if v, ok := p._SafeDo(
		func() interface{} {
			if nil != p.processor {
				return (*p.processor).GetLogPath()
			}
			return nil
		},
	).(*string); ok {
		return v
	}
	return nil
}

//output log.
func (p *Slogger) record(logLevel LogLevel, format string, v ...interface{}) {
	//filter to loglevel.
	if logLevel < p.settings.LogLevel {
		return
	}

	fileName, fileLine := _GetFileInfoFromStack(3)

	p._SafeDo(
		func() interface{} {
			//Enqueue.
			if nil != p.task {
				p.task._Offer(&SloggerData{
					logLevel:          logLevel,
					currentTimeMillis: GetCurrentTimeMillis(),
					logMessage:        fmt.Sprintf("%s(%d): ", fileName, fileLine) + fmt.Sprintf(format, v...) + "\n",
				})
				return nil
			} else {
				return errors.New("task is nil.")
			}
		},
	)
}

func (p *Slogger) Critical(format string, v ...interface{}) {
	p.record(CRITICAL, format, v...)
}

func (p *Slogger) Error(format string, v ...interface{}) {
	p.record(ERROR, format, v...)
}

func (p *Slogger) Warn(format string, v ...interface{}) {
	p.record(WARN, format, v...)
}

func (p *Slogger) Info(format string, v ...interface{}) {
	p.record(INFO, format, v...)
}

func (p *Slogger) Debug(format string, v ...interface{}) {
	p.record(DEBUG, format, v...)
}

func _GetFileInfoFromStack(depth int) (string, int) {
	_, file, line, _ := runtime.Caller(depth)
	_, fileName := filepath.Split(file)
	return fileName, line
}
