package slogger

import "fmt"

type SloggerData struct {
	CurrentTimeMillis int64
	LogLevel          LogLevel
	LogMessage        string
	SourceName        string
	SourceLine        int
}

func (self *SloggerData) ToLogMessage() string {
	return fmt.Sprintf(
		"%s %s %s(%d): %s",
		ConvertTimeStamp(self.CurrentTimeMillis, Full),
		self.LogLevel.ToString(),
		self.SourceName,
		self.SourceLine,
		self.LogMessage,
	)
}
