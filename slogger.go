package slogger

import (
	"os"
	"sync"
)

type SloggerSettings struct {
	LogLevel     LogLevel
	LogName      string
	LogDirectory string
	LogFilepath  string
}

type Slogger struct {
	mutex    sync.Mutex
	settings SloggerSettings
	logFp    *os.File
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

//output log.
func (p *Slogger) record(logLevel LogLevel, format string, v ...interface{}) {
	p._SafeDo(
		func() interface{} {
			return nil
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
