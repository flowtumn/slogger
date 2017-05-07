package slogger

import "fmt"

type SloggerProcessorStdout struct {
}

func (self *SloggerProcessorStdout) GetLogPath() *string {
	return nil
}

func (self *SloggerProcessorStdout) Record(setting SloggerSettings, data *SloggerData) error {
	fmt.Println(data.ToLogMessage())
	return nil
}

func (self *SloggerProcessorStdout) Shutdown() {

}

func CreateSloggerProcessorStdout() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorStdout{}
	return &r
}
