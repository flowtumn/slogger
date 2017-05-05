package slogger

type SloggerRecordCount struct {
	Debug    int64
	Info     int64
	Warn     int64
	Error    int64
	Critical int64
}

func (p *SloggerRecordCount) _CountUpOnLogLevel(logLevel LogLevel) {
	switch logLevel {
	default:
	case DEBUG:
		p.Debug = p.Debug + 1
	case INFO:
		p.Info = p.Info + 1
	case WARN:
		p.Warn = p.Warn + 1
	case ERROR:
		p.Error = p.Error + 1
	case CRITICAL:
		p.Critical = p.Critical + 1
	}
}
