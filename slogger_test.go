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

func _writeLog(p *Slogger) {
	if nil == p {
		return
	}

	(*p).Debug("DEBUG %d", 100)
	(*p).Info("INFO %f", 1.2345)
	(*p).Warn("WARN %s", "warn message")
	(*p).Error("ERROR")
	(*p).Critical("Critical %+v", errors.New("Critical"))
}

func Test_SLogger_Base(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          WARN,
		LogName:           "dummy1",
		LogDirectory:      "dumym2",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}
	r := Slogger{}

	if reflect.DeepEqual(*r.Settings(), DATA) {
		t.Errorf("It is equal to DATA. This is incorrect")
	}

	r.Initialize(DATA)

	if !reflect.DeepEqual(*r.Settings(), DATA) {
		t.Errorf("Settings must be DATA.")
	}

}

func Test_SLogger_Debug(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          DEBUG,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	//Debugなら、全て記録される。
	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    1,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Info(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          INFO,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     1,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_WARN(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          WARN,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     1,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Error(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          ERROR,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    1,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_Critical(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          CRITICAL,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 0,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	_writeLog(&r)

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 1,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Output count does not match.")
	}

}

func Test_SLogger_Cycle_1sec(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          DEBUG,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 1000, //Cycle 1sec.
	}
	WRITE_COUNT := (int64)(10)
	TEST_CHECK := SloggerOutputCount{
		Critical: WRITE_COUNT,
		Error:    WRITE_COUNT,
		Warn:     WRITE_COUNT,
		Info:     WRITE_COUNT,
		Debug:    WRITE_COUNT,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	for i := 0; i < (int)(WRITE_COUNT); i++ {
		_writeLog(&r)
	}

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 0,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record is NG.")
	}

	//1秒待ってから、再度書き込み。
	time.Sleep(1 * time.Second)

	//Debugの記録時にcycleを見るので、Debugだけが書き込んだ回数は1多い
	r.Debug("first")
	TEST_CHECK.Debug++

	r.Info("Second")

	if !reflect.DeepEqual(
		*r.Counters(),
		TEST_CHECK,
	) {
		t.Errorf("Output count does not match. case 1.")
	}

	TEST_CHECK.Info++

	//1秒待ってから、再度書き込み。
	time.Sleep(1 * time.Second)

	r.Warn("Third")
	TEST_CHECK.Warn++

	if !reflect.DeepEqual(
		*r.Counters(),
		TEST_CHECK,
	) {
		t.Errorf("Output count does not match. case 2.")
	}
}

func Test_SLogger_Cycle_100sec(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          DEBUG,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 100000, //Cycle 100sec.
	}
	WRITE_COUNT := (int64)(1234)
	TEST_CHECK := SloggerOutputCount{
		Critical: WRITE_COUNT,
		Error:    WRITE_COUNT,
		Warn:     WRITE_COUNT,
		Info:     WRITE_COUNT,
		Debug:    WRITE_COUNT,
	}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)
	for i := 0; i < (int)(WRITE_COUNT); i++ {
		_writeLog(&r)
	}

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 0,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record is NG.")
	}

	//Flush
	r.Close()

	if !reflect.DeepEqual(
		*r.Counters(),
		TEST_CHECK,
	) {
		t.Errorf("Output count does not match.")
	}
}

func Test_SLogger_MT(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:          DEBUG,
		LogName:           "TEST",
		LogDirectory:      "./",
		LogExtension:      "log",
		RecordCycleMillis: 1000000, //Cycle 1000sec.
	}
	WORKER_COUNT := (int64)(4)
	genRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	var TEST_CHECK = SloggerOutputCount{}

	r := Slogger{}

	defer func() {
		r.Close()
		os.Remove(*r.GetLogPath())
	}()

	r.Initialize(DATA)

	var waiter sync.WaitGroup
	for i := 0; i < (int)(WORKER_COUNT); i++ {
		waiter.Add(1)

		//書き込む数は最大1万回。
		count := genRand.Int63n(10000)
		TEST_CHECK.Debug = TEST_CHECK.Debug + count
		TEST_CHECK.Info = TEST_CHECK.Info + count
		TEST_CHECK.Warn = TEST_CHECK.Warn + count
		TEST_CHECK.Error = TEST_CHECK.Error + count
		TEST_CHECK.Critical = TEST_CHECK.Critical + count

		go func(writeCount int) {
			for ii := 0; ii < writeCount; ii++ {
				_writeLog(&r)
			}
			waiter.Done()
		}((int)(count))
	}

	waiter.Wait()

	if !reflect.DeepEqual(
		*r.Counters(),
		SloggerOutputCount{
			Critical: 0,
			Error:    0,
			Warn:     0,
			Info:     0,
			Debug:    0,
		},
	) {
		t.Errorf("Record is NG.")
	}

	//Flush
	r.Close()

	if !reflect.DeepEqual(
		*r.Counters(),
		TEST_CHECK,
	) {
		t.Errorf("Output count does not match. Actual: %+v,  Expected: %+v", *r.Counters(), TEST_CHECK)
	}
}
