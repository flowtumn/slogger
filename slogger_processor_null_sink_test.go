package slogger

import "testing"

func Test_SLogger_Null_Sink_Processor(t *testing.T) {
	processor := CreateSloggerProcessorNullSink()

	defer func() {
		(*processor).Shutdown()
	}()

	if nil == processor {
		t.Fatalf("Null sink is nil.")
	}

	if nil != (*processor).GetLogPath() {
		t.Errorf("return LogPath is nil.")
	}

	if nil != (*processor).Record(SloggerSettings{}, nil) {
		t.Fatalf("Record error.")
	}
}
