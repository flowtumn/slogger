package slogger

import (
	"reflect"
	"testing"
)

func Test_SLogger_Settings_Base(t *testing.T) {
	DATA := SloggerSettings{
		LogLevel:     DEBUG,
		LogName:      "NAME",
		LogDirectory: "/tmp/tmp",
		LogExtension: "log",
	}

	DATA_1 := DATA
	DATA_1.LogDirectory = DATA.LogDirectory + "\\"
	DATA_1.Trim()

	if !reflect.DeepEqual(
		DATA,
		DATA_1,
	) {
		t.Errorf("Setting data dose not match. case1.")
	}

	DATA_1.LogDirectory = DATA.LogDirectory + "/"
	DATA_1.Trim()
	if !reflect.DeepEqual(
		DATA,
		DATA_1,
	) {
		t.Errorf("Setting data dose not match. case2.")
	}

	DATA_1.LogDirectory = DATA.LogDirectory + "\\/\\/////\\////\\/"
	DATA_1.Trim()
	if !reflect.DeepEqual(
		DATA,
		DATA_1,
	) {
		t.Errorf("Setting data dose not match. case3.")
	}
}
