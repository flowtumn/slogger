package slogger

type _SloggerData struct {
	logLevel          LogLevel
	currentTimeMillis int64
	logMessage        string
}

func (p *_SloggerData) _toLogMessage() string {
	return ConvertTimeStamp(p.currentTimeMillis, Full) + " " + p.logLevel.toStr() + " " + p.logMessage
}
