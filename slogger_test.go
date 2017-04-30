package slogger

import (
	"reflect"
	"testing"
)

func Test_SLogger_Base(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "dummy1",
		LogDirectory: "dumym2",
		LogFilepath:  "dummy3",
	}
	r := Slogger{}

	if reflect.DeepEqual(r.Settings(), DATA) {
		t.Errorf("It is equal to DATA. This is incorrect")
	}

	r.Initialize(DATA)

	if !reflect.DeepEqual(r.Settings(), DATA) {
		t.Errorf("Settings must be DATA.")
	}

}
