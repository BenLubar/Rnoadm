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
	flag.Int64Var(&Seed, "seed", 0, "the world seed")

	flag.Parse()

	os.MkdirAll(seedFilename(), 755)
	f, err := os.OpenFile(filepath.Join(seedFilename(), "admin.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	AdminLog = log.New(f, "[ADMIN] ", log.Ldate|log.Ltime|log.Lshortfile)

	go func() {
		for _ = range time.Tick(10 * time.Minute) {
			EachLoadedZone(func(z *Zone) {
				err := z.Save()
				if err != nil {
					log.Printf("ZONE %d %d: %v", z.X, z.Y, err)
				}
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
		err := z.Save()
		if err != nil {
			log.Printf("ZONE %d %d: %v", z.X, z.Y, err)
		}
	})
}
