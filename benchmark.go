package main

import (
	"fmt"
	"os"
	"time"

	"github.com/flowtumn/slogger"
)

const (
	TEST_REPEAT = 3
	TEST_COUNT  = 0x20000
)

func benchmark(level slogger.LogLevel, processor *slogger.SloggerProcessor) {
	logger, err := slogger.CreateSlogger(
		slogger.SloggerSettings{
			LogLevel:     level,
			LogName:      "EXAMPLE-SLOGGER",
			LogDirectory: "./",
			LogExtension: "log",
		},
		processor,
	)

	if nil != err {
		panic(err)
	}

	defer func() {
		logger.Close()
		if v := logger.GetLogPath(); nil != v {
			os.Remove(*v)
		}
	}()

	for i := 0; i < TEST_COUNT; i++ {
		logger.Debug("Debug Message.")
		logger.Info("Info Message.")
		logger.Warn("Warn Message.")
		logger.Error("Error Message.")
		logger.Critical("Critical Message.")
	}
}

func Measurement(name string, repeat int, f func()) {
	tickCount := int64(0)
	for i := 0; i < repeat; i++ {
		start := time.Now().UnixNano()
		f()
		tickCount = tickCount + time.Now().UnixNano() - start
	}

	fmt.Printf("%s: avg. %d ms\n", name, (tickCount/(int64)(repeat))/(int64)(time.Millisecond))
}

func main() {
	for _, level := range []slogger.LogLevel{slogger.DEBUG, slogger.INFO, slogger.WARN, slogger.ERROR, slogger.CRITICAL} {
		fmt.Printf("Measurement loglevel -> %s\n", level.ToString())
		Measurement(" Null Sink", TEST_REPEAT, func() { benchmark(level, slogger.CreateSloggerProcessorNullSink()) })
		Measurement("Cache File", TEST_REPEAT, func() { benchmark(level, slogger.CreateSloggerProcessorCacheFile()) })
		Measurement("      File", TEST_REPEAT, func() { benchmark(level, slogger.CreateSloggerProcessorFile()) })
	}
}
