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

func ConvertTimeStamp(currentTime int64, layout LogLayout) string {
	return time.Unix(currentTime, 0).Format(
		func() string {
			return logLayoutTable[layout]
		}(),
	)
}

func GetCurrentTimeMillis() int64 {
	return time.Now().Unix()
}

func GetCurrentTimeNanos() int64 {
	return time.Now().UnixNano()
}

func GetTimeStamp(layout LogLayout) string {
	return ConvertTimeStamp(time.Now().Unix(), layout)
}
