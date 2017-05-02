package slogger

import (
	"fmt"
	"os"
	"sync"
)

type SloggerSettings struct {
	LogLevel          LogLevel
	LogName           string
	LogDirectory      string
	LogExtension      string
	RecordCycleMillis int64
}

type SloggerOutputCount struct {
	Debug    int64
	Info     int64
	Warn     int64
	Error    int64
	Critical int64
}

type _SloggerBuffer struct {
	logLevel          LogLevel
	currentTimeMillis int64
	logMessage        string
}

type Slogger struct {
	mutex                sync.Mutex
	settings             SloggerSettings
	count                SloggerOutputCount
	buffer               []_SloggerBuffer
	isRuning             bool
	currentTimeStamp     string
	lastRecordTimeMillis int64
	logPath              string
	logFp                *os.File
}

func _CreateLogFileName(prefix string, suffix string) string {
	return prefix + "-" + GetTimeStamp(Normal) + "." + suffix
}

func (p *_SloggerBuffer) _toLogMessage() string {
	return ConvertTimeStamp(p.currentTimeMillis, Full) + " " + p.logLevel.toStr() + " " + p.logMessage
}

func (p *Slogger) _SafeDo(f func() interface{}) interface{} {
	p.mutex.Lock()
	defer func() {
		p.mutex.Unlock()
	}()

	return f()
}

func (p *Slogger) Initialize(settings SloggerSettings) {
	p._SafeDo(
		func() interface{} {
			p.settings = settings
			p.count = SloggerOutputCount{}
			p.buffer = []_SloggerBuffer{}
			p.isRuning = false
			if nil != p.logFp {
				p.logFp.Close()
				p.logFp = nil
			}
			return nil
		},
	)
}

func (p *Slogger) Close() {
	p._SafeDo(
		func() interface{} {
			//Flush.
			p._RecordProcess()

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

	logData := _SloggerBuffer{
		logLevel:          logLevel,
		currentTimeMillis: GetCurrentTimeMillis(),
		logMessage:        fmt.Sprintf(format, v...),
	}

	p._SafeDo(
		func() interface{} {
			//Buffering.
			p.buffer = append(p.buffer, logData)

			if 0 < p.settings.RecordCycleMillis && GetCurrentTimeMillis()-p.lastRecordTimeMillis < p.settings.RecordCycleMillis {
				//Cycle time not exceeded.
				return nil
			}

			return p._RecordProcess()
		},
	)
}

func (p *Slogger) _RecordProcess() interface{} {
	for n, v := range p.buffer {
		if err := p._UpdateSink(v.currentTimeMillis); nil != err {
			p.buffer = p.buffer[:n]
			return err
		}

		//write log.
		p.logFp.WriteString(v._toLogMessage())

		//Count up.
		p._CountUpOnLogLevel(v.logLevel)
	}

	//init.
	p.buffer = []_SloggerBuffer{}

	//time update.
	p.lastRecordTimeMillis = GetCurrentTimeMillis()
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
	p.record(
		CRITICAL,
		format,
		v...,
	)
}

func (p *Slogger) Error(format string, v ...interface{}) {
	p.record(
		ERROR,
		format,
		v...,
	)
}

func (p *Slogger) Warn(format string, v ...interface{}) {
	p.record(
		WARN,
		format,
		v...,
	)
}

func (p *Slogger) Info(format string, v ...interface{}) {
	p.record(
		INFO,
		format,
		v...,
	)
}

func (p *Slogger) Debug(format string, v ...interface{}) {
	p.record(DEBUG, format, v...)
}
