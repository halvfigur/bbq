package main

import (
	"log"
	"os"
)

var _d = os.Getenv("DBUS") != ""

var _dl = log.New(os.Stderr, "[dbus] ", log.Lmicroseconds)

func debug(format string, v ...interface{}) {
	if _d {
		_dl.Printf(format, v...)
	}
}
