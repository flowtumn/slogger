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
	mutex               sync.Mutex
	settings            SloggerSettings
	count               SloggerRecordCount
	task                *_SloggerWorker
	processor           *SloggerProcessor
	isRuning            bool
	currentTimeStamp    string
	lastRecordTimeNanos int64
	logPath             string
	logFp               *os.File
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

func (p *Slogger) Initialize(settings SloggerSettings, processor *SloggerProcessor) {
	p.Close()

	p._SafeDo(
		func() interface{} {
			p.settings = settings
			p.count = SloggerRecordCount{}
			p.isRuning = false
			p.logFp = nil

			p.processor = processor

			//Create Worker.
			p.task = _CreateSloggerWorker(
				func(buf *SloggerData) {
					if nil != p.processor {
						(*p.processor).Record(buf)
					}
				},
			)

			return nil
		},
	)
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

func (p *Slogger) _RecordProcess(v *SloggerData) interface{} {
	if nil == v {
		return errors.New("SloggerBuffer is nil.")
	}

	if err := p._UpdateSink(v.currentTimeMillis); nil != err {
		return err
	}

	//Log write.
	if _, err := p.logFp.WriteString(v._toLogMessage()); nil != err {
		return err
	}

	p.count._CountUpOnLogLevel(v.logLevel)

	return nil
}

func (p *Slogger) _UpdateSink(currentTimeMillis int64) interface{} {
	tm := ConvertTimeStamp(currentTimeMillis, Normal)
	if p.currentTimeStamp == tm {
		return nil
	}

	//Update a currentTimeStamp.
	p.currentTimeStamp = tm
	if nil != p.logFp {
		p.logFp.Close()
	}

	p.logPath = _CreateLogFileName(p.settings.LogName, p.settings.LogExtension)
	if fp, err := os.OpenFile(p.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); nil == err {
		p.logFp = fp
		return nil
	} else {
		return err
	}
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
