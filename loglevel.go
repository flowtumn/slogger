package slogger

type LogLevel int32

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRITICAL
)

func (v LogLevel) toStr() string {
	switch v {
	default:
		return "UNKNOWN"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	}
}
