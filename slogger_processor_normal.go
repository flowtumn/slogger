package slogger

import "os"

type SloggerProcessorNormal struct {
	currentTimeStamp string
	logPath          string
	logFp            *os.File
}

func (self *SloggerProcessorNormal) GetLogPath() *string {
	r := "qwer.txt"
	return &r
}

func (self *SloggerProcessorNormal) Record(data *SloggerData) error {
	if err := self._UpdateSink(data.currentTimeMillis); nil != err {
		return err
	}

	//Log write.
	if _, err := self.logFp.WriteString(data.ToLogMessage()); nil != err {
		return err
	}

	return nil
}

func (self *SloggerProcessorNormal) Shutdown() {
	if nil != self.logFp {
		self.logFp.Close()
		self.logFp = nil
	}
}

func (self *SloggerProcessorNormal) _UpdateSink(currentTimeMillis int64) error {
	tm := ConvertTimeStamp(currentTimeMillis, Normal)

	if self.currentTimeStamp == tm {
		return nil
	}

	//Update a currentTimeStamp.
	self.currentTimeStamp = tm
	if nil != self.logFp {
		self.logFp.Close()
	}

	// self.logPath = _CreateLogFileName(p.settings.LogName, p.settings.LogExtension)
	self.logPath = _CreateLogFileName("TEST", "log")
	if fp, err := os.OpenFile(self.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); nil == err {
		self.logFp = fp
		return nil
	} else {
		return err
	}

	return nil
}

func CreateSloggerProcessorNormal() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorNormal{}
	return &r
}
