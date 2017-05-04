# slogger(Simple logger)

 Go言語から使用出来る、簡易的な Logger Library. (Thread-safe)

LogLevelでのフィルターと、日付毎に出力先を変えてログを記録する機能を有しています。

* * *

## Install

```text
go get github.com/flowtumn/slogger
```

## Settings

　SloggerにはSettingsが必要です。

```Settings
slogger.SloggerSettings{
    LogLevel:          slogger.DEBUG,       // LogLevel。指定レベルより以上のログが記録されます。
    LogName:           "EXAMPLE-SLOGGER",   // Loggerの名前。
    LogDirectory:      "/tmp/",             // Log出力先。
    LogExtension:      "log",               // Logの拡張子。
}
```

## Example

　使い方のサンプル。

Loggerには EXAMPLE-SLOGGER という名称を設定し、/tmp/をlogDirectoryとしています。

```
package main

import (
    "github.com/flowtumn/slogger"
)

func main() {
    r := slogger.Slogger{}
    r.Initialize(
        slogger.SloggerSettings{
            LogLevel:          slogger.DEBUG,
            LogName:           "EXAMPLE-SLOGGER",
            LogDirectory:      "/tmp/",
            LogExtension:      "log",
        },
    )

    defer func() {
        //Closeを呼びlogをflush
        r.Close()
    }()

    r.Debug("Debug Message. %+v", map[int]int{100: 200, 300: 400, 500: 600});
    r.Info("Info Message.");
    r.Warn("Warn Message.");
    r.Error("Error Message.");
    r.Critical("Critical Message.");
}
```

## Log file.

```
$ cat /tmp/EXAMPLE-SLOGGER-2017-05-03.log
2017-05-03 00:29:25 [DEBUG] main.go(24): Debug Message. map[100:200 300:400 500:600]
2017-05-03 00:29:25 [INFO] main.go(25): Info Message.
2017-05-03 00:29:25 [WARN] main.go(26): Warn Message.
2017-05-03 00:29:25 [ERROR] main.go(27): Error Message.
2017-05-03 00:29:25 [CRITICAL] main.go(28): Critical Message.
```
