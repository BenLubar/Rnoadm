package main

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"code.google.com/p/go.net/websocket"
	"compress/gzip"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"github.com/BenLubar/Rnoadm/resource"
	"io"
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
var clientGzip []byte

func init() {
	http.HandleFunc("/", httpHandler)
	http.Handle("/ws", websocket.Handler(websocketHandler))

	h := sha1.New()
	h.Write(resource.Resource["client.js"])
	clientHash = hex.EncodeToString(h.Sum(nil))

	var buf bytes.Buffer
	g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	g.Write(resource.Resource["client.js"])
	g.Close()
	clientGzip = buf.Bytes()
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if b, ok := resource.Resource[r.URL.Path[1:]]; ok {
			if strings.HasSuffix(r.URL.Path, ".png") {
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Expires", time.Now().UTC().AddDate(1, 0, 0).Format(http.TimeFormat))
			} else if strings.HasSuffix(r.URL.Path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
				if r.URL.Path == "/client.js" && strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
					w.Header().Set("Content-Encoding", "gzip")
					b = clientGzip
				}
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
<link href="http://fonts.googleapis.com/css?family=Jolly+Lodger|Open+Sans+Condensed:300&subset=latin,latin-ext,cyrillic,cyrillic-ext,greek-ext,greek,vietnamese" rel="stylesheet">
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
	Admin *string
	Auth  *struct {
		Login    string
		Password string
	}
	CharacterCreation *struct {
		Command string
	}
	Walk *struct {
		X, Y uint8
	}
	Interact *struct {
		ID     uint64 `json:"I,string"`
		Option int    `json:"O"`
		X, Y   uint8
	}
}

type packetClientHash struct {
	ClientHash string
}

type packetKick struct {
	Kick string
}

type packetMessage struct {
	Message string
}

type _SetHUD struct {
	Name string
	Data map[string]interface{}
}

type packetSetHUD struct {
	SetHUD _SetHUD
}

type NetworkedObject struct {
	Name     string             `json:"N",omitempty`
	Options  []string           `json:"O",omitempty`
	Sprite   string             `json:"I",omitempty`
	Colors   []Color            `json:"C"`
	Scale    uint8              `json:"S,omitempty"`
	Height   uint16             `json:"H,omitempty"`
	Attached []*NetworkedObject `json:"A,omitempty"`
	Moves    bool               `json:"M,omitempty"`
}

type TileChange struct {
	ID      uint64 `json:",string"`
	X, Y    uint8
	Removed bool             `json:"R,omitempty"`
	Obj     *NetworkedObject `json:"O,omitempty"`
}

type packetTileChange struct {
	TileChange       []TileChange
	PlayerX, PlayerY uint8
	ResetZone        bool
}

type InventoryItem struct {
	ID  uint64           `json:"I,string"`
	Obj *NetworkedObject `json:"O"`
}

type packetInventory struct {
	Inventory []InventoryItem
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
				if err != io.EOF {
					log.Printf("[%s] %v", addr, err)
				}
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
	inventory := make(chan packetInventory, 1)
	messages := make(chan string, 8)

	var updates <-chan TileChange
	var updateQueue packetTileChange
	updateTimer := time.NewTicker(5 * time.Millisecond)
	defer updateTimer.Stop()

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
				player.inventory = inventory
				player.hud = hud
				player.messages = messages
				player.kick = kick

				onlinePlayersLock.Lock()
				if otherSession, ok := OnlinePlayers[player.ID]; ok {
					otherSession.Kick("Logged in from another location.")
				}
				OnlinePlayers[player.ID] = player
				onlinePlayersLock.Unlock()

				player.Lock()
				zone, zoneUpdates := GrabZone(player.ZoneX, player.ZoneY, player)
				updates = zoneUpdates
				player.zone = zone
				if player.Hero != nil {
					tx, ty := player.TileX, player.TileY
					tile := zone.Tile(player.TileX, player.TileY)
					if tile == nil {
						tx, ty = 127, 127
						player.TileX, player.TileY = 127, 127
						tile = zone.Tile(127, 127)
					}
					player.Unlock()

					zone.Lock()
					tile.Add(player)
					SendZoneTileChange(zone.X, zone.Y, TileChange{
						ID:  player.NetworkID(),
						X:   tx,
						Y:   ty,
						Obj: player.Serialize(),
					})
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
						SendZoneTileChange(zone.X, zone.Y, TileChange{
							ID:      player.NetworkID(),
							Removed: true,
						})
					}
					if zone != nil {
						ReleaseZone(zone, player)
					}
				}()

				updateQueue.ResetZone = true
				updateQueue.TileChange = zone.AllTileChange()

				if player.Hero == nil {
					player.CharacterCreationCommand("")
					continue
				}

				player.backpackDirty = make(chan struct{}, 1)
				select {
				case player.backpackDirty <- struct{}{}:
				default:
				}

				player.SetHUD("", nil)
			}

			if player == nil {
				continue
			}

			if p.Admin != nil {
				AdminCommand(addr, player, *p.Admin)
			}

			if p.CharacterCreation != nil {
				player.CharacterCreationCommand(p.CharacterCreation.Command)
			}

			if player.Hero == nil {
				continue
			}

			if p.Walk != nil {
				player.Lock()
				zone := player.zone
				tx, ty := player.TileX, player.TileY
				player.Unlock()

				zone.Lock()
				schedule := MoveSchedule(FindPath(zone, tx, ty, p.Walk.X, p.Walk.Y, true))
				zone.Unlock()

				player.Lock()
				player.schedule = &schedule
				player.Unlock()
			}

			if p.Interact != nil {
				player.Lock()
				zone := player.zone
				player.Unlock()

				tile := zone.Tile(p.Interact.X, p.Interact.Y)
				if tile == nil {
					continue
				}
				zone.Lock()
				for _, o := range tile.Objects {
					if o.NetworkID() == p.Interact.ID {
						if p.Interact.Option < 0 {
							switch p.Interact.Option {
							case -1:
								player.SendMessage(o.Examine())
							}
							break
						}
						zone.Unlock()
						o.Interact(p.Interact.X, p.Interact.Y, player, zone, p.Interact.Option)
						zone.Lock()
						break
					}
				}
				zone.Unlock()
			}

		case p := <-inventory:
			websocket.JSON.Send(conn, p)

		case p := <-hud:
			websocket.JSON.Send(conn, p)

		case message := <-messages:
			websocket.JSON.Send(conn, packetMessage{
				Message: message,
			})

		case message := <-kick:
			websocket.JSON.Send(conn, packetKick{
				Kick: message,
			})
			return

		case update, ok := <-updates:
			if !ok {
				updates = nil
				updateQueue.ResetZone = true
				updateQueue.TileChange = updateQueue.TileChange[:0]
				continue
			}
			updateQueue.TileChange = append(updateQueue.TileChange, update)

		case <-updateTimer.C:
			if len(updateQueue.TileChange) > 0 {
				player.Lock()
				updateQueue.PlayerX = player.TileX
				updateQueue.PlayerY = player.TileY
				player.Unlock()

				leftover := updateQueue.TileChange[:0]
				if len(updateQueue.TileChange) > 100 {
					leftover = updateQueue.TileChange[100:]
					updateQueue.TileChange = updateQueue.TileChange[:100]
				}

				websocket.JSON.Send(conn, updateQueue)

				updateQueue.ResetZone = false
				updateQueue.TileChange = leftover
			}
		}
	}
}
