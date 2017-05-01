package slogger

import "testing"

func Test_LogTime_Base(t *testing.T) {
	if "2006-01-02" != logLayoutTable[Normal] {
		t.Errorf("Error: Log layout Normal.")
	}

	if "2006-01-02 15:04:05" != logLayoutTable[Full] {
		t.Errorf("Error: Log layout Full.")
	}
}
