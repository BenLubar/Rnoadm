package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
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

	srv := &http.Server{
		Addr: ":2064",
		TLSConfig: &tls.Config{
			// BUG: https://code.google.com/p/go/issues/detail?id=6121
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}

	go func() {
		for {
			err := srv.ListenAndServeTLS("rnoadm-cert.pem", "rnoadm-key.pem")
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
