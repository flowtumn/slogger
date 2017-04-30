package slogger

import "testing"

func Test_LogLevel_Base(t *testing.T) {
	if logLevelTable[DEBUG] != DEBUG.toStr() {
		t.Errorf("Error: Debug to string.")
	}

	if logLevelTable[INFO] != INFO.toStr() {
		t.Errorf("Error: Info to string.")
	}

	if logLevelTable[WARN] != WARN.toStr() {
		t.Errorf("Error: Warn to string.")
	}

	if logLevelTable[ERROR] != ERROR.toStr() {
		t.Errorf("Error: Error to string.")
	}

	if logLevelTable[CRITICAL] != CRITICAL.toStr() {
		t.Errorf("Error: Critical to string.")
	}
}
