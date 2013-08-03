package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"code.google.com/p/go.net/websocket"
	"compress/gzip"
	"encoding/gob"
	"github.com/BenLubar/Rnoadm/resource"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var OnlinePlayers = make(map[uint64]*Player)
var onlinePlayersLock sync.Mutex

func init() {
	http.HandleFunc("/", httpHandler)
	http.Handle("/ws", websocket.Handler(websocketHandler))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if b, ok := resource.Resource[r.URL.Path[1:]]; ok {
			if strings.HasSuffix(r.URL.Path, ".png") {
				w.Header().Set("Content-Type", "image/png")
			} else if strings.HasSuffix(r.URL.Path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			}
			w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
			w.Header().Set("Cache-Control", "public")
			w.Write(b)
			return
		}
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Rnoadm</title>
<style>
html {
	background: #000;
	text-align: center;
}
</style>
</head>
<body>
<canvas></canvas>
<script src="client.js"></script>
</body>
</html>`))
}

type userLogin struct {
	ID       uint64
	Username string
	Password []byte

	Registered     time.Time
	RegisteredAddr string
}

type packetIn struct {
	Auth *struct {
		Login    string
		Password string
	}
}

type packetKick struct {
	Kick string
}

func websocketHandler(conn *websocket.Conn) {
	defer conn.Close()

	addr := conn.Request().RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			addr = addr[:i]
			break
		}
	}

	packets := make(chan packetIn)
	go func() {
		for {
			var p packetIn
			err := websocket.JSON.Receive(conn, &p)
			if err != nil {
				log.Printf("[%s] %v", addr, err)
				close(packets)
				return
			}
			select {
			case packets <- p:
			case <-time.After(time.Second):
				log.Printf("[%s] dropped a packet (server)", addr)
			}
		}
	}()

	for {
		select {
		case p, ok := <-packets:
			if !ok {
				return
			}

			if p.Auth != nil {
				// TODO: throttling
				filename := filepath.Join(seedFilename(), "login"+Base32Encode([]byte(strings.ToLower(p.Auth.Login)))+".gz")
				var login userLogin

				f, err := os.Open(filename)
				if err != nil {
					hash, err := bcrypt.GenerateFromPassword([]byte(p.Auth.Password), bcrypt.DefaultCost)
					if err != nil {
						log.Printf("[%s] registration for %q failed: %v", addr, p.Auth.Login, err)
						return
					}
					login = userLogin{
						ID:             newUserID(),
						Username:       p.Auth.Login,
						Password:       hash,
						Registered:     time.Now().UTC(),
						RegisteredAddr: addr,
					}
					f, err = os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
					if err != nil {
						log.Printf("[%s] registration for %q failed: %v", addr, p.Auth.Login, err)
						return
					}
					g, _ := gzip.NewWriterLevel(f, gzip.BestCompression) // only possible error is invalid level
					err = gob.NewEncoder(g).Encode(&login)
					g.Close()
					f.Close()
				} else {
					g, err := gzip.NewReader(f)
					if err != nil {
						f.Close()
						log.Printf("[%s] login for %q failed: %v", addr, p.Auth.Login, err)
						return
					}
					err = gob.NewDecoder(g).Decode(&login)
					g.Close()
					f.Close()
					if err != nil {
						log.Printf("[%s] login for %q failed: %v", addr, p.Auth.Login, err)
						return
					}
					err = bcrypt.CompareHashAndPassword(login.Password, []byte(p.Auth.Password))
					if err != nil {
						websocket.JSON.Send(conn, packetKick{
							Kick: "password incorrect",
						})
						return
					}
				}
				websocket.JSON.Send(conn, packetKick{
					Kick: "logging in has temporarily been disabled",
				})
				return
			}
		}
	}
}
