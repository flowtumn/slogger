package slogger

import "testing"

func Test_Atomic_Bool_Base(t *testing.T) {
	v := AtomicBool{}

	if false != v.Get() {
		t.Errorf("dose not match. case1.")
	}

	v.Set(true)
	if true != v.Get() {
		t.Errorf("dose not match. case1")
	}

	v.Set(false)
	if false != v.Get() {
		t.Errorf("dose not match. case3")
	}
}
