package main

import (
	"code.google.com/p/go.net/websocket"
	_ "github.com/BenLubar/Rnoadm/critter"
	"github.com/BenLubar/Rnoadm/hero"
	_ "github.com/BenLubar/Rnoadm/material"
	"github.com/BenLubar/Rnoadm/world"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.Handle("/ws", websocket.Handler(socketHandler))
	http.HandleFunc("/", staticHandler)
	go func() {
		for {
			err := http.ListenAndServe(":2064", nil)
			if err != nil {
				log.Print(err)
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		for _ = range time.Tick(time.Second / 5) {
			world.Think()
		}
	}()

	defer world.SaveAllZones()
	defer hero.SaveAllPlayers()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	<-sigCh
}
