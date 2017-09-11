package slogger

import "testing"

func Test_SLogger_Stdout_Processor(t *testing.T) {
	processor := CreateSloggerProcessorStdout()

	if nil == processor {
		t.Fatalf("Stdout processor is nil.")
	}

	defer func() {
		(*processor).Shutdown()
	}()

	if nil != (*processor).GetLogPath() {
		t.Errorf("return LogPath not is nil.")
	}

	if nil != (*processor).Record(SloggerSettings{}, &SloggerData{LogLevel: DEBUG, LogMessage: "Test_SLogger_Stdout_Processor"}) {
		t.Fatalf("Record error.")
	}
}
