package slogger

import "testing"

func Test_LogLevel_Base(t *testing.T) {
	if "DEBUG" != DEBUG.toStr() {
		t.Errorf("Error: Debug to string.")
	}

	if "INFO" != INFO.toStr() {
		t.Errorf("Error: Info to string.")
	}

	if "WARN" != WARN.toStr() {
		t.Errorf("Error: Warn to string.")
	}

	if "ERROR" != ERROR.toStr() {
		t.Errorf("Error: Error to string.")
	}

	if "CRITICAL" != CRITICAL.toStr() {
		t.Errorf("Error: Critical to string.")
	}
}
