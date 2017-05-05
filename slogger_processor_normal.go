package slogger

import "os"

type SloggerProcessorNormal struct {
	logFp *os.File
}

func (p *SloggerProcessorNormal) GetLogPath() *string {
	r := "qwer.txt"
	return &r
}

func (p *SloggerProcessorNormal) Record(data *SloggerData) {

}

func (p *SloggerProcessorNormal) Shutdown() {
	if nil != p.logFp {
		p.logFp.Close()
		p.logFp = nil
	}
}

func CreateSloggerProcessorNormal() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorNormal{}
	return &r
}
