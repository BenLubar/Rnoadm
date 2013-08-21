// +build !windows

package main

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

func init() {
	w, err := syslog.New(syslog.LOG_INFO, "rnoadm")
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(w, os.Stderr))
}
