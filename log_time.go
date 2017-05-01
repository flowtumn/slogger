package slogger

import "time"

type LogLayout int

const (
	Normal = iota
	Full
)

var logLayoutTable = map[LogLayout]string{
	Normal: "2006-01-02",
	Full:   "2006-01-02 15:04:05",
}

func GetTimeStamp(layout LogLayout) string {
	day := time.Now()
	return day.Format(
		func() string {
			return logLayoutTable[layout]
		}(),
	)
}
