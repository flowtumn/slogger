package slogger

type SloggerProcessor interface {
	GetLogPath() *string
	Record(SloggerSettings, *SloggerData) error
	Shutdown()
}
