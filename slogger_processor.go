package slogger

type SloggerProcessor interface {
	GetLogPath() *string
	Record(*SloggerData) error
	Shutdown()
}
