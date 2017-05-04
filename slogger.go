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

type SloggerOutputCount struct {
	Debug    int64
	Info     int64
	Warn     int64
	Error    int64
	Critical int64
}

type Slogger struct {
	mutex               sync.Mutex
	settings            SloggerSettings
	count               SloggerOutputCount
	task                *_SloggerWorker
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

func (p *Slogger) Initialize(settings SloggerSettings) {
	p.Close()

	p._SafeDo(
		func() interface{} {
			p.settings = settings
			p.count = SloggerOutputCount{}
			p.isRuning = false
			p.logFp = nil

			//Create Worker.
			p.task = _CreateSloggerWorker(
				func(buf *_SloggerData) {
					p._RecordProcess(buf)
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

			if nil != p.logFp {
				p.logFp.Close()
				p.logFp = nil
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

func (p *Slogger) Counters() *SloggerOutputCount {
	if v, ok := p._SafeDo(
		func() interface{} {
			return p.count
		},
	).(SloggerOutputCount); ok {
		return &v
	}

	return nil
}

func (p *Slogger) GetLogPath() *string {
	if v, ok := p._SafeDo(
		func() interface{} {
			return p.logPath
		},
	).(string); ok {
		return &v
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
				p.task._Offer(&_SloggerData{
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

func (p *Slogger) _RecordProcess(v *_SloggerData) interface{} {
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

	p._CountUpOnLogLevel(v.logLevel)

	return nil
}

func (p *Slogger) _CountUpOnLogLevel(logLevel LogLevel) {
	switch logLevel {
	default:
	case DEBUG:
		p.count.Debug = p.count.Debug + 1
	case INFO:
		p.count.Info = p.count.Info + 1
	case WARN:
		p.count.Warn = p.count.Warn + 1
	case ERROR:
		p.count.Error = p.count.Error + 1
	case CRITICAL:
		p.count.Critical = p.count.Critical + 1
	}
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
