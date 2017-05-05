package slogger

import (
	"reflect"
	"testing"
)

func Test_SLogger_RecordCount_Base(t *testing.T) {
	v := SloggerRecordCount{}

	v._CountUpOnLogLevel(DEBUG)
	if !reflect.DeepEqual(
		v,
		SloggerRecordCount{
			Debug:    1,
			Info:     0,
			Warn:     0,
			Error:    0,
			Critical: 0,
		},
	) {
		t.Errorf("Record debug count dose not match.")
	}

	v._CountUpOnLogLevel(INFO)
	if !reflect.DeepEqual(
		v,
		SloggerRecordCount{
			Debug:    1,
			Info:     1,
			Warn:     0,
			Error:    0,
			Critical: 0,
		},
	) {
		t.Errorf("Record info count dose not match.")
	}

	v._CountUpOnLogLevel(WARN)
	if !reflect.DeepEqual(
		v,
		SloggerRecordCount{
			Debug:    1,
			Info:     1,
			Warn:     1,
			Error:    0,
			Critical: 0,
		},
	) {
		t.Errorf("Record warn count dose not match.")
	}

	v._CountUpOnLogLevel(ERROR)
	if !reflect.DeepEqual(
		v,
		SloggerRecordCount{
			Debug:    1,
			Info:     1,
			Warn:     1,
			Error:    1,
			Critical: 0,
		},
	) {
		t.Errorf("Record error count dose not match.")
	}

	v._CountUpOnLogLevel(CRITICAL)
	if !reflect.DeepEqual(
		v,
		SloggerRecordCount{
			Debug:    1,
			Info:     1,
			Warn:     1,
			Error:    1,
			Critical: 1,
		},
	) {
		t.Errorf("Record critical count dose not match..")
	}
}
