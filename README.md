# slogger(Simple logger)

[![Build Status](https://travis-ci.org/flowtumn/slogger.svg?branch=master)](https://travis-ci.org/flowtumn/slogger)
[![Coverage Status](https://coveralls.io/repos/github/flowtumn/slogger/badge.svg?branch=master)](https://coveralls.io/github/flowtumn/slogger?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

 Go言語から使用出来る、簡易的な Logger Library. (Thread-safe)

LogLevelでのフィルターと、日付毎に出力先を変えてログを記録する機能を有しています。

記録するファイル名のフォーマットは ```{LogName}-yyyy-MM-dd.{LogExtension} ``` となります。

* * *

## Install

```text
go get github.com/flowtumn/slogger
```

## Settings

　SloggerにはSettingsが必要です。

```Settings
slogger.SloggerSettings{
    LogLevel:          slogger.DEBUG,       // LogLevel。指定レベル以上のログが記録されます。
    LogName:           "EXAMPLE-SLOGGER",   // Loggerの名前。
    LogDirectory:      "/tmp",              // Log出力先。
    LogExtension:      "log",               // Logの拡張子。
}
```

## Example

　使い方のサンプル。

Loggerには EXAMPLE-SLOGGER という名称を設定し、/tmpをlogDirectoryとしています。

```
package main

import (
    "github.com/flowtumn/slogger"
)

func main() {
    logger, err := slogger.CreateSlogger(
        slogger.SloggerSettings{
            LogLevel:     slogger.DEBUG,
            LogName:      "EXAMPLE-SLOGGER",
            LogDirectory: "/tmp",
            LogExtension: "log",
        },
        slogger.CreateSloggerProcessorFile(), //Logを処理するProcessor。Fileに記録する。
    )

    if nil != err {
        panic(err)
    }

    defer func() {
        //Closeを呼びlogをflush
        logger.Close()
    }()

    logger.Debug("Debug Message. %+v", map[int]int{100: 200, 300: 400, 500: 600})
    logger.Info("Info Message.")
    logger.Warn("Warn Message.")
    logger.Error("Error Message.")
    logger.Critical("Critical Message.")
}
```

## Log file.

```
$ cat /tmp/EXAMPLE-SLOGGER-2017-05-06.log
2017-05-06 23:41:46 [DEBUG] main.go(22): Debug Message. map[100:200 300:400 500:600]
2017-05-06 23:41:46 [INFO] main.go(23): Info Message.
2017-05-06 23:41:46 [WARN] main.go(24): Warn Message.
2017-05-06 23:41:46 [ERROR] main.go(25): Error Message.
2017-05-06 23:41:46 [CRITICAL] main.go(26): Critical Message.
```
