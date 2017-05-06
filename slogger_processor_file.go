package slogger

import "os"

type SloggerProcessorFile struct {
	currentTimeStamp string
	logPath          string
	logFp            *os.File
}

func _CreateLogFileName(prefix string, suffix string) string {
	return prefix + "-" + GetTimeStamp(Normal) + "." + suffix
}

func (self *SloggerProcessorFile) GetLogPath() *string {
	r := self.logPath
	return &r
}

func (self *SloggerProcessorFile) Record(setting SloggerSettings, data *SloggerData) error {
	if err := self._UpdateSink(&setting, data.CurrentTimeMillis); nil != err {
		return err
	}

	//Log write.
	if _, err := self.logFp.WriteString(data.ToLogMessage() + "\n"); nil != err {
		return err
	}

	return nil
}

func (self *SloggerProcessorFile) Shutdown() {
	if nil != self.logFp {
		self.logFp.Close()
		self.logFp = nil
	}
}

func (self *SloggerProcessorFile) _UpdateSink(setting *SloggerSettings, currentTimeMillis int64) error {
	tm := ConvertTimeStamp(currentTimeMillis, Normal)

	if self.currentTimeStamp == tm {
		return nil
	}

	//Update a currentTimeStamp.
	self.currentTimeStamp = tm
	if nil != self.logFp {
		self.logFp.Close()
	}

	self.logPath = _CreateLogFileName(setting.LogDirectory+"/"+setting.LogName, setting.LogExtension)
	if fp, err := os.OpenFile(self.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); nil == err {
		self.logFp = fp
		return nil
	} else {
		return err
	}

	return nil
}

func CreateSloggerProcessorFile() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorFile{}
	return &r
}
