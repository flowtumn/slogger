package slogger

import (
	"reflect"
	"testing"
	"time"
)

const (
	_TEST_LOG_PATH = "\\1?test//?aa"
)

type _TestProcessor struct {
	callLogPath  int
	callRecord   int
	callShutdown int
	argsSetting  SloggerSettings
}

func (self *_TestProcessor) GetLogPath() *string {
	r := _TEST_LOG_PATH
	self.callLogPath = self.callLogPath + 1
	return &r
}

func (self *_TestProcessor) Record(settings SloggerSettings, _ *SloggerData) error {
	self.callRecord = self.callRecord + 1
	self.argsSetting = settings
	return nil
}

func (self *_TestProcessor) Shutdown() {
	self.callShutdown = self.callShutdown + 1
}

//実体はTestProcessorな、SloggerProcessorを作成。
func _CreateSloggerProcessorTest() (*SloggerProcessor, *_TestProcessor) {
	var r SloggerProcessor
	rt := &_TestProcessor{}
	r = rt
	return &r, rt
}

func Test_SLogger_Empty_Processor(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "dummy1",
		LogDirectory: "dumym2",
		LogExtension: "log",
	}
	processor, testProcessor := _CreateSloggerProcessorTest()

	r, err := CreateSlogger(
		SETTINGS,
		processor,
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	if v := (*processor).GetLogPath(); nil != v {
		if _TEST_LOG_PATH != *v {
			t.Errorf("dose not match. Expected: %s, Actual: %s", _TEST_LOG_PATH, *v)
		}
	} else {
		t.Errorf("return LogPath is nil.")
	}

	if !(1 == testProcessor.callLogPath && 0 == testProcessor.callRecord && 0 == testProcessor.callShutdown) {
		t.Errorf("GetLogPath counter does not match")
	}

	//SettingsのLogLevelは WARN なので、WARN以上のものが3回呼ばれる
	r.Debug("TEST")
	r.Info("TEST")
	r.Warn("TEST")
	r.Error("TEST")
	r.Critical("TEST")

	//少し待機
	time.Sleep(100 * time.Millisecond)

	if !(3 == testProcessor.callRecord && 0 == testProcessor.callShutdown) {
		t.Errorf("Record counter does not match")
	}

	//Shutdownを発火
	r.Close()
	if !(3 == testProcessor.callRecord && 1 == testProcessor.callShutdown) {
		t.Errorf("Shutdown counter does not match")
	}

	//Setting確認
	if !reflect.DeepEqual(
		testProcessor.argsSetting,
		SETTINGS,
	) {
		t.Errorf("Settings must be DATA.")
	}

	//以降Record/Shutdownは呼ばれない

	r.Debug("TEST")
	r.Info("TEST")
	r.Warn("TEST")
	r.Error("TEST")
	r.Critical("TEST")
	r.Close()

	if !(3 == testProcessor.callRecord && 1 == testProcessor.callShutdown) {
		t.Errorf("Shutdown and Record counter does not match")
	}

}
