package slogger

import (
	"regexp"
	"testing"
)

func Test_SLogger_Data_Base(t *testing.T) {
	r := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}\:[0-9]{2}\:[0-9]{2} \[INFO\]`)

	if true != r.MatchString((&SloggerData{
		LogLevel:          INFO,
		CurrentTimeMillis: 0,
		LogMessage:        "TEST",
	}).ToLogMessage()) {
		t.Errorf("Broken log message.")
	}
}
