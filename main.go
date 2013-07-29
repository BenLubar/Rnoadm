package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var Seed int64

func main() {
	flag.Int64Var(&Seed, "seed", time.Now().UnixNano(), "the world seed (default: the number of nanoseconds since midnight UTC on 1970-01-01)")

	flag.Parse()

	f, err := os.OpenFile(filepath.Join(seedFilename(), "admin.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	AdminLog = log.New(f, "[ADMIN] ", log.Ldate|log.Ltime|log.Lshortfile)

	go func() {
		for _ = range time.Tick(time.Minute) {
			EachLoadedZone(func(z *Zone) {
				z.Save()
			})
		}
	}()

	go func() {
		for _ = range time.Tick(200 * time.Millisecond) {
			EachLoadedZone(func(z *Zone) {
				z.Think()
			})
		}
	}()

	go func() {
		for {
			log.Print(http.ListenAndServe(":2064", nil))
			time.Sleep(time.Second)
		}
	}()

	sigkill := make(chan os.Signal, 1)
	signal.Notify(sigkill, os.Kill)
	<-sigkill
	EachLoadedZone(func(z *Zone) {
		z.Save()
	})
}
