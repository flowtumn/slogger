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
	Critical int64
	Error    int64
	Warn     int64
	Info     int64
	Debug    int64
}

type _SloggerBuffer struct {
	currentTimeMillis int64
	log               string
}

type Slogger struct {
	mutex                sync.Mutex
	settings             SloggerSettings
	count                SloggerOutputCount
	buffer               []_SloggerBuffer
	currentTimeStamp     string
	lastRecordTimeMillis int64
	logPath              string
	logFp                *os.File
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
	p._SafeDo(
		func() interface{} {
			p.settings = settings
			p.count = SloggerOutputCount{}
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

func (p *Slogger) _CanRecord() bool {
	if v, ok := p._SafeDo(
		func() interface{} {
			return GetCurrentTimeMillis()-p.lastRecordTimeMillis < p.settings.RecordCycleMillis
		},
	).(bool); ok {
		return v
	}
	return false
}

//output log.
func (p *Slogger) record(fs func(), logLevel LogLevel, format string, v ...interface{}) {
	//filter to loglevel.
	if logLevel < p.settings.LogLevel {
		return
	}

	// _ := GetTimeStamp(Full) + " " +
	// 	logLevel.toStr() + " " +
	// 	fmt.Sprintf(format, v...) + "\n"
	if nil == p._UpdateSink() {
		p._SafeDo(
			func() interface{} {
				p.logFp.WriteString(
					GetTimeStamp(Full) + " " +
						logLevel.toStr() + " " +
						fmt.Sprintf(format, v...) +
						"\n",
				)

				//Success.
				fs()

				//time update.
				p.lastRecordTimeMillis = GetCurrentTimeMillis()
				return nil
			},
		)
	}
}

func (p *Slogger) _UpdateSink() interface{} {
	return p._SafeDo(
		func() interface{} {
			tm := GetTimeStamp(Normal)
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
		},
	)
}

func (p *Slogger) Critical(format string, v ...interface{}) {
	p.record(
		func() {
			p.count.Critical = p.count.Critical + 1
		},
		CRITICAL,
		format,
		v...,
	)
}

func (p *Slogger) Error(format string, v ...interface{}) {
	p.record(
		func() {
			p.count.Error = p.count.Error + 1
		},
		ERROR,
		format,
		v...,
	)
}

func (p *Slogger) Warn(format string, v ...interface{}) {
	p.record(
		func() {
			p.count.Warn = p.count.Warn + 1
		},
		WARN,
		format,
		v...,
	)
}

func (p *Slogger) Info(format string, v ...interface{}) {
	p.record(
		func() {
			p.count.Info = p.count.Critical + 1
		},
		INFO,
		format,
		v...,
	)
}

func (p *Slogger) Debug(format string, v ...interface{}) {
	p.record(
		func() {
			p.count.Debug = p.count.Debug + 1
		},
		DEBUG,
		format,
		v...,
	)
}
