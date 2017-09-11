package slogger

import (
	"os"
	"testing"
	"time"
)

func Test_SLogger_Cache_File_Processor(t *testing.T) {
	SETTINGS := SloggerSettings{
		LogLevel:     WARN,
		LogName:      "TEST_CACHE_FILE_PROCESSOR",
		LogDirectory: "./",
		LogExtension: "log",
	}
	TEST_LOG_DATA := SloggerData{
		CurrentTimeMillis: 1494756991, //2017-5-14
		LogLevel:          INFO,
		LogMessage:        "TEST MESSAGE",
		SourceName:        "slogger_processor_cache_file_test.go",
	}
	TEST_RECORD_COUNT := 5
	CHECK_FILE_LIST := []string{}

	defer func() {
		for _, path := range CHECK_FILE_LIST {
			os.Remove(path)
		}
	}()

	processor, ok := (*CreateSloggerProcessorCacheFile()).(*SloggerProcessorCacheFile)

	if !ok {
		t.Fatalf("Down cast error.")
	}

	if p := (*processor)._Poll(); 0 != len(*p) {
		t.Fatalf("Poll data count not zero. case1.")
	}

	(*processor).Record(SETTINGS, &TEST_LOG_DATA)
	time.Sleep(100 * time.Millisecond)
	CHECK_FILE_LIST = append(CHECK_FILE_LIST, (*processor.GetLogPath()))

	//記録時間を変更。(2017-5-10)
	TEST_LOG_DATA.CurrentTimeMillis = 1494412437
	(*processor).Record(SETTINGS, &TEST_LOG_DATA)
	time.Sleep(100 * time.Millisecond)
	CHECK_FILE_LIST = append(CHECK_FILE_LIST, (*processor.GetLogPath()))

	//ファイルは２つ出来ていないと駄目。
	if CHECK_FILE_LIST[0] == CHECK_FILE_LIST[1] {
		t.Fatalf("It is recorded in the same file")
	}

	(*processor).Shutdown()

	for _, v := range CHECK_FILE_LIST {
		info, err := os.Stat(v)
		//logの出力を確認。
		if nil != err {
			t.Fatalf("Not found log file.")
		}
		//Shutdownしているので書き込まれていなければならない。
		if 0 == info.Size() {
			t.Fatalf("It is not saved in the log file")
		}
	}

	if p := (*processor)._Poll(); 0 != len(*p) {
		t.Fatalf("Poll data count not zero. case2.")
	}

	for i := 0; i < TEST_RECORD_COUNT; i++ {
		(*processor).Record(SETTINGS, &TEST_LOG_DATA)
	}

	time.Sleep(100 * time.Millisecond)

	//Shutdownしているので、詰んだデータは処理されていない。
	if v := (*processor)._Poll(); nil != v {
		if TEST_RECORD_COUNT != len(*v) {
			t.Errorf("Poll data count does not match.")
		}
	} else {
		t.Fatalf("Poll data is nil.")
	}
}
