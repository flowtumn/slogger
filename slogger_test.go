package slogger

import (
	"errors"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	DEFAULT_TEST_WAIT_TIMES = 500 * time.Millisecond
)

func _writeLog(p *Slogger, waitTimes time.Duration) {
	if nil == p {
		return
	}

	(*p).Debug("DEBUG %d", 100)
	(*p).Info("INFO %f", 1.2345)
	(*p).Warn("WARN %s", "warn message")
	(*p).Error("ERROR")
	(*p).Critical("Critical %+v", errors.New("Critical"))

	if 0 < waitTimes {
		time.Sleep(waitTimes)
	}
}

func Test_SLogger_Base(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "dummy1",
		LogDirectory: "dumym2",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	if true != r.IsRunning() {
		t.Fatalf("Slogger not running..")
	}

	if !reflect.DeepEqual(r.Settings(), SETTINGS) {
		t.Errorf("Settings must be SETTINGS.")
	}

	r.Close()
	if false != r.IsRunning() {
		t.Fatalf("Slogger Worker not stopped.")
	}
}

func Test_SLogger_Fail(t *testing.T) {
	r, _ := CreateSlogger(
		SloggerSettings{},
		nil,
	)

	if r != nil {
		t.Fatalf("Error: Slogger Object not nil.")
	}
}

func Test_SLogger_Close(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "dummy1",
		LogDirectory: "dumym2",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	r.Close()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	//Close後に記録されてはいけない
	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 0,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("It is written after closed.")
	}
}

func Test_SLogger_Debug(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     DEBUG,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	//Debugなら、全て記録される。
	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    1,
		},
	) {
		t.Errorf("Record count does not match.")
	}
}

func Test_SLogger_Info(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     INFO,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    0,
		},
	) {
		t.Errorf("Record count does not match.")
	}
}

func Test_SLogger_WARN(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record count does not match.")
	}
}

func Test_SLogger_Error(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     ERROR,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 1,
			Error:    1,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record count does not match.")
	}
}

func Test_SLogger_Critical(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     CRITICAL,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorFile(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	_writeLog(r, DEFAULT_TEST_WAIT_TIMES)

	if !reflect.DeepEqual(
		*r.RecordCounter(),
		SloggerRecordCount{
			Critical: 1,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record count does not match.")
	}

}

func Test_SLogger_MT(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     DEBUG,
		LogName:      "TEST",
		LogDirectory: "./",
		LogExtension: "log",
	}
	WORKER_COUNT := (int64)(8)
	genRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var TEST_CHECK = SloggerRecordCount{}

	r, err := CreateSlogger(
		SETTINGS,
		CreateSloggerProcessorNullSink(),
	)

	if nil != err {
		t.Fatalf("Failed to CreateSlogger.")
	}

	defer func() {
		r.Close()
	}()

	var waiter sync.WaitGroup
	for i := 0; i < (int)(WORKER_COUNT); i++ {
		waiter.Add(1)

		//書き込む数は最大10万回。
		count := genRand.Int63n(100000)
		TEST_CHECK.Debug = TEST_CHECK.Debug + count
		TEST_CHECK.Info = TEST_CHECK.Info + count
		TEST_CHECK.Warn = TEST_CHECK.Warn + count
		TEST_CHECK.Error = TEST_CHECK.Error + count
		TEST_CHECK.Critical = TEST_CHECK.Critical + count

		go func(writeCount int) {
			for ii := 0; ii < writeCount; ii++ {
				//Wait無しで書き込み続ける
				_writeLog(r, 0)
			}
			waiter.Done()
		}((int)(count))
	}

	waiter.Wait()

	//Flush
	r.Close()

	if !reflect.DeepEqual(
		*r.RecordCounter(),
		TEST_CHECK,
	) {
		t.Errorf("Record count does not match. Actual: %+v,  Expected: %+v", *r.RecordCounter(), TEST_CHECK)
	}
}
