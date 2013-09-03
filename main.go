package main

import (
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
	_ "github.com/BenLubar/Rnoadm/critter"
	"github.com/BenLubar/Rnoadm/hero"
	_ "github.com/BenLubar/Rnoadm/material"
	"github.com/BenLubar/Rnoadm/maybetls"
	"github.com/BenLubar/Rnoadm/world"
	"log"
	"math/rand"
	"net"
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
			NextProtos: []string{"http/1.1"},

			// BUG: https://code.google.com/p/go/issues/detail?id=6121
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
	}
	cert, err := tls.LoadX509KeyPair("rnoadm-cert.pem", "rnoadm-key.pem")
	if err != nil {
		panic(err)
	}
	srv.TLSConfig.Certificates = []tls.Certificate{cert}

	go func() {
		for {
			ln, err := net.Listen("tcp", ":2064")
			if err != nil {
				log.Print(err)
				time.Sleep(time.Second)
				continue
			}
			ln = maybetls.Listener(ln, srv.TLSConfig)
			err = srv.Serve(ln)
			if err != nil {
				ln.Close()
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
