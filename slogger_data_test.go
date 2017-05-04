package slogger

import (
	"regexp"
	"testing"
)

func Test_SLogger_Data_Base(t *testing.T) {
	r := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}\:[0-9]{2}\:[0-9]{2} \[INFO\]`)

	if true != r.MatchString((&_SloggerData{
		logLevel:          INFO,
		currentTimeMillis: 0,
		logMessage:        "TEST",
	})._toLogMessage()) {
		t.Errorf("Broken log message.")
	}
}