package slogger

type LogLevel int32

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRITICAL
)

var logLevelTable = map[LogLevel]string{
	DEBUG:    "[DEBUG]",
	INFO:     "[INFO]",
	WARN:     "[WARN]",
	ERROR:    "[ERROR]",
	CRITICAL: "[CRITICAL]",
}

func (v LogLevel) toStr() string {
	if v, ok := logLevelTable[v]; ok {
		return v
	}
	return "UNKNOWN"
}
