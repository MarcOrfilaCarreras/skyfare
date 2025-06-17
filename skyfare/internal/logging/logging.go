package logging

import (
    "log"
    "os"
)

var quiet bool = false

func SetQuiet(q bool) {
    quiet = q
}

func Println(v ...interface{}) {
    if !quiet {
        log.Println(v...)
    }
}

func Printf(format string, v ...interface{}) {
    if !quiet {
        log.Printf(format, v...)
    }
}

func Fatalf(format string, v ...interface{}) {
    if !quiet {
        log.Fatalf(format, v...)
    } else {
        os.Exit(1)
    }
}
