package slogger

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func _writeLog(p *Slogger) {
	if nil == p {
		return
	}

	(*p).Debug("DEBUG %d", 100)
	(*p).Info("INFO %f", 1.2345)
	(*p).Warn("WARN %s", "warn message")
	(*p).Error("ERROR")
	(*p).Critical("Critical %+v", errors.New("Critical"))
}

func Test_SLogger_Base(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          WARN,
		LogName:           "dummy1",
		LogDirectory:      "dumym2",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}
	r := Slogger{}

	if reflect.DeepEqual(*r.Settings(), DATA) {
		t.Errorf("It is equal to DATA. This is incorrect")
	}

	r.Initialize(DATA)

	if !reflect.DeepEqual(*r.Settings(), DATA) {
		t.Errorf("Settings must be DATA.")
	}

}

func Test_SLogger_Debug(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          DEBUG,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	//Debugなら、全て記録される。
	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    1,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Info(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          INFO,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_WARN(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          WARN,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Error(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          ERROR,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Critical(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          CRITICAL,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}
