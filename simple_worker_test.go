package slogger

import (
	"testing"
	"time"
)

func Test_Simple_Worker_Base(t *testing.T) {
	TEST_COUNT := 100
	counter := 0
	task := CreateWork(
		func() bool {
			counter = counter + 1
			//falseを返したら、Workerは終了。
			return counter < TEST_COUNT
		},
		10, //10msのwaitを挟みながら、関数を呼び出す。
	)

	for task.IsRunning() {
		time.Sleep(10 * time.Millisecond)
	}

	if TEST_COUNT != counter {
		t.Errorf("Call count does not match.")
	}
}
