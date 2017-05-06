package slogger

import "strings"

type SloggerSettings struct {
	LogLevel     LogLevel
	LogName      string
	LogDirectory string
	LogExtension string
}

func (self *SloggerSettings) Trim() {
	for {
		length := len(self.LogDirectory)
		self.LogDirectory = strings.TrimRight(self.LogDirectory, "/")
		self.LogDirectory = strings.TrimRight(self.LogDirectory, "\\")
		if length == len(self.LogDirectory) {
			return
		}
	}
}
