package slogger

import "os"

type SloggerProcessorFile struct {
	currentTimeStamp string
	logPath          string
	logFp            *os.File
}

type UpdateSinkResult int

const (
	Update UpdateSinkResult = iota + 1
	NoChange
	Fail
)

func _CreateLogFileName(timeStamp int64, prefix string, suffix string) string {
	return prefix + "-" + ConvertTimeStamp(timeStamp, Normal) + "." + suffix
}

func (self *SloggerProcessorFile) GetLogPath() *string {
	r := self.logPath
	return &r
}

func (self *SloggerProcessorFile) Record(setting SloggerSettings, data *SloggerData) error {
	if _, err := self.UpdateSink(&setting, data.CurrentTimeMillis); nil != err {
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

func (self *SloggerProcessorFile) UpdateSink(setting *SloggerSettings, currentTimeMillis int64) (UpdateSinkResult, error) {
	tm := ConvertTimeStamp(currentTimeMillis, Normal)

	if self.currentTimeStamp == tm {
		return NoChange, nil
	}

	if nil != self.logFp {
		self.logFp.Close()
	}

	self.logPath = _CreateLogFileName(currentTimeMillis, setting.LogDirectory+"/"+setting.LogName, setting.LogExtension)
	if fp, err := os.OpenFile(self.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); nil == err {
		//Update a currentTimeStamp.
		self.currentTimeStamp = tm
		self.logFp = fp
		return Update, nil
	} else {
		return Fail, err
	}
}

func CreateSloggerProcessorFile() *SloggerProcessor {
	var r SloggerProcessor
	r = &SloggerProcessorFile{}
	return &r
}
