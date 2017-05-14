package slogger

type SloggerProcessorNullSink struct {
}

func (self *SloggerProcessorNullSink) GetLogPath() *string {
	return nil
}

func (self *SloggerProcessorNullSink) Record(setting SloggerSettings, data *SloggerData) error {
	return nil
}

func (self *SloggerProcessorNullSink) Shutdown() {
}

func CreateSloggerProcessorNullSink() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorNullSink{}
	return &r
}
