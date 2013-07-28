package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var Seed int64

func main() {
	flag.Int64Var(&Seed, "seed", time.Now().UnixNano(), "the world seed (default: the number of nanoseconds since midnight UTC on 1970-01-01)")

	flag.Parse()

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

	for {
		log.Print(http.ListenAndServe(":2064", nil))
		time.Sleep(time.Second)
	}
}
