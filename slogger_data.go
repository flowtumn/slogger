package slogger

type SloggerData struct {
	currentTimeMillis int64
	logLevel          LogLevel
	logMessage        string
}

func (p *SloggerData) _toLogMessage() string {
	return ConvertTimeStamp(p.currentTimeMillis, Full) + " " + p.logLevel.toStr() + " " + p.logMessage
}
