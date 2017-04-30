package slogger

import (
	"errors"
	"reflect"
	"testing"
)

func Test_SLogger_Base(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "dummy1",
		LogDirectory: "dumym2",
		LogExtension: ".log",
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
		LogLevel:     DEBUG,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: ".log",
	}

	r := Slogger{}
	r.Initialize(DATA)
	r.Debug("DEBUG %d", 100)
	r.Info("INFO %f", 1.2345)
	r.Warn("WARN %s", "warn message")
	r.Error("ERROR")
	r.Critical("Critical %+v", errors.New("Critical"))
}
