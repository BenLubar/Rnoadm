package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"code.google.com/p/go.net/websocket"
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"github.com/BenLubar/Rnoadm/resource"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var OnlinePlayers = make(map[uint64]*Player)
var onlinePlayersLock sync.Mutex

var clientHash string

func init() {
	http.HandleFunc("/", httpHandler)
	http.Handle("/ws", websocket.Handler(websocketHandler))

	h := sha1.New()
	h.Write(resource.Resource["client.js"])
	clientHash = hex.EncodeToString(h.Sum(nil))
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
	CharacterCreation *struct {
		Command string
	}
}

type packetClientHash struct {
	ClientHash string
}

type packetKick struct {
	Kick string
}

type _SetHUD struct {
	Name string
	Data map[string]interface{}
}

type packetSetHUD struct {
	SetHUD _SetHUD
}

var bruteThrottle = make(map[string]uint8)
var bruteThrottleLock sync.Mutex

func init() {
	go func() {
		for {
			bruteThrottleLock.Lock()
			for addr, n := range bruteThrottle {
				if n <= 1 {
					delete(bruteThrottle, addr)
				} else {
					bruteThrottle[addr]--
				}
			}
			bruteThrottleLock.Unlock()

			time.Sleep(time.Minute)
		}
	}()
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

	bruteThrottleLock.Lock()
	if bruteThrottle[addr] > 5 {
		bruteThrottleLock.Unlock()
		websocket.JSON.Send(conn, packetKick{
			Kick: "Too many login attempts. Come back later.",
		})
		return
	}
	bruteThrottleLock.Unlock()

	websocket.JSON.Send(conn, packetClientHash{
		ClientHash: clientHash,
	})

	packets := make(chan packetIn)
	go func() {
		for {
			var p packetIn
			err := websocket.JSON.Receive(conn, &p)
			if err != nil {
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

	var player *Player
	kick := make(chan string, 1)
	hud := make(chan packetSetHUD, 16)

	for {
		select {
		case p, ok := <-packets:
			if !ok {
				return
			}

			if p.Auth != nil {
				if player != nil {
					return
				}
				if p.Auth.Login == "" || p.Auth.Password == "" {
					return
				}
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
						bruteThrottleLock.Lock()
						bruteThrottle[addr]++
						if bruteThrottle[addr] >= 5 {
							bruteThrottle[addr] = 20
						}
						bruteThrottleLock.Unlock()
						websocket.JSON.Send(conn, packetKick{
							Kick: "password incorrect",
						})
						return
					}
				}
				player, err = LoadPlayer(login.ID)
				if err != nil {
					player = &Player{
						ID: login.ID,
					}
					player.Seed.Seed(int64(login.ID))
				}
				player.LastLogin = time.Now().UTC()
				player.LastLoginAddr = addr
				player.Save()
				player.kick = kick
				player.hud = hud

				onlinePlayersLock.Lock()
				if otherSession, ok := OnlinePlayers[player.ID]; ok {
					otherSession.Kick("Logged in from another location.")
				}
				OnlinePlayers[player.ID] = player
				onlinePlayersLock.Unlock()

				player.Lock()
				if player.Hero != nil {
					zone := GrabZone(player.ZoneX, player.ZoneY)
					player.zone = zone
					tile := zone.Tile(player.TileX, player.TileY)
					if tile == nil {
						player.TileX, player.TileY = 127, 127
						tile = zone.Tile(127, 127)
					}
					player.Unlock()

					zone.Lock()
					tile.Add(player)
					zone.Unlock()
				} else {
					player.Unlock()
				}

				defer func() {
					onlinePlayersLock.Lock()
					if session := OnlinePlayers[player.ID]; session == player {
						delete(OnlinePlayers, player.ID)
						defer player.Save()
					}
					onlinePlayersLock.Unlock()

					player.Lock()
					zone := player.zone
					var tile *Tile
					if zone != nil {
						tile = zone.Tile(player.TileX, player.TileY)
					}
					player.Unlock()

					if tile != nil {
						zone.Lock()
						tile.Remove(player)
						zone.Unlock()
					}
				}()

				if player.Hero == nil {
					player.CharacterCreation("")
					continue
				}
				player.Kick("logging in has temporarily been disabled")
			}
			if player == nil {
				continue
			}

			if p.CharacterCreation != nil {
				player.CharacterCreation(p.CharacterCreation.Command)
			}

			if player.Hero == nil {
				continue
			}

		case p := <-hud:
			websocket.JSON.Send(conn, p)

		case message := <-kick:
			websocket.JSON.Send(conn, packetKick{
				Kick: message,
			})
			return
		}
	}
}
