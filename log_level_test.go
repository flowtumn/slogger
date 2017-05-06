package slogger

import "testing"

func Test_LogLevel_Base(t *testing.T) {
	if logLevelTable[DEBUG] != DEBUG.ToString() {
		t.Errorf("Error: Debug to string.")
	}

	if logLevelTable[INFO] != INFO.ToString() {
		t.Errorf("Error: Info to string.")
	}

	if logLevelTable[WARN] != WARN.ToString() {
		t.Errorf("Error: Warn to string.")
	}

	if logLevelTable[ERROR] != ERROR.ToString() {
		t.Errorf("Error: Error to string.")
	}

	if logLevelTable[CRITICAL] != CRITICAL.ToString() {
		t.Errorf("Error: Critical to string.")
	}
}
