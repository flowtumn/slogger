package slogger

import (
	"reflect"
	"testing"
)

func Test_SLogger_Worker_Base(t *testing.T) {
	DATA := SloggerData{
		logLevel:          CRITICAL,
		currentTimeMillis: 12345,
		logMessage:        "TEST",
	}
	REPEAT := 10
	var called = 0

	worker := _CreateSloggerWorker(
		func(v *SloggerData) {
			called = called + 1
			if nil != v && reflect.DeepEqual(*v, DATA) {
				return
			}
			t.Errorf("does not match.")
		},
	)

	//Workerは動いている
	if true != worker.IsRunning() {
		t.Errorf("Worker not running.")
	}

	for i := 0; i < REPEAT; i++ {
		//DATAをpush
		worker._Offer(&DATA)
	}

	//Workerを停止
	worker._Shutdown()

	if false != worker.IsRunning() {
		t.Errorf("Worker not stopped.")
	}

	//停止してもpush
	worker._Offer(&DATA)
	worker._Offer(&DATA)
	worker._Offer(&DATA)

	//呼ばれる回数は REPEAT と同じでないとNG
	if REPEAT != called {
		t.Errorf("the called does not match. Expected: %d, Actual: %d", called, REPEAT)
	}
}
