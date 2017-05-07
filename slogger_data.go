package slogger

import "strconv"

type SloggerData struct {
	CurrentTimeMillis int64
	LogLevel          LogLevel
	LogMessage        string
	SourceName        string
	SourceLine        int
}

func (self *SloggerData) ToLogMessage() string {
	return ConvertTimeStamp(self.CurrentTimeMillis, Full) + " " +
		self.LogLevel.ToString() + " " +
		self.SourceName + "(" + strconv.Itoa(self.SourceLine) + "): " +
		self.LogMessage
}
