package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/BenLubar/Rnoadm/hero"
	"github.com/BenLubar/Rnoadm/world"
	"io"
	"log"
	"net"
	"time"
)

type clientPacket struct {
	Admin []string
	Auth  *hero.LoginPacket
	Walk  *struct {
		X, Y uint8
	}
	HUD *struct {
		Name string      `json:"N"`
		Data interface{} `json:"D"`
	}
	Chat *string
}

type packetKick struct {
	Kick string
}

var packetClientHash struct {
	ClientHash uint64 `json:",string"`
}

type packetUpdate struct {
	Update           []packetUpdateUpdate
	PlayerX, PlayerY uint8
}

type packetUpdateUpdate struct {
	ID     uint64              `json:"I,string"`
	X      uint8               `json:"X"`
	Y      uint8               `json:"Y"`
	FromX  *uint8              `json:"Fx,omitempty"`
	FromY  *uint8              `json:"Fy,omitempty"`
	Remove bool                `json:"R,omitempty"`
	Object *packetUpdateObject `json:"O,omitempty"`
}

type packetUpdateObject struct {
	Name    string               `json:"N"`
	Sprites []packetUpdateSprite `json:"S"`
}

type packetUpdateSprite struct {
	Sheet string                 `json:"S"`
	Color string                 `json:"C"`
	Extra map[string]interface{} `json:"E,omitempty"`
}

type packetInventory struct {
	Inventory []*packetUpdateObject
}

type packetMessage struct {
	Message []hero.Message `json:"Msg"`
}

func addSprites(u *packetUpdateObject, obj world.Visible) *packetUpdateObject {
	sheet := obj.Sprite()
	x, y := obj.SpritePos()
	width, height := obj.SpriteSize()
	scale := obj.Scale()
	animation := obj.AnimationType()
	for i, c := range obj.Colors() {
		if c == "" {
			continue
		}
		u.Sprites = append(u.Sprites, packetUpdateSprite{
			Sheet: sheet,
			Color: c,
			Extra: map[string]interface{}{
				"w": width,
				"h": height,
				"x": x,
				"y": y + uint(i),
				"s": scale,
				"a": animation,
			},
		})
	}
	for _, a := range obj.Attached() {
		addSprites(u, a)
	}
	return u
}

type packetHUD struct {
	HUD hero.HUD
}

func socketHandler(ws *websocket.Conn) {
	defer ws.Close()

	addr, _, err := net.SplitHostPort(ws.Request().RemoteAddr)
	if err != nil {
		panic(err)
	}

	packets := make(chan clientPacket)
	go func() {
		defer close(packets)
		for {
			var packet clientPacket
			err = websocket.JSON.Receive(ws, &packet)
			if err != nil {
				if err == io.EOF {
					return
				}
				if _, ok := err.(*net.OpError); ok {
					return
				}
				log.Printf("%s: %T %v", addr, err, err)
				return
			}
			packets <- packet
		}
	}()

	websocket.JSON.Send(ws, packetClientHash)

	var (
		player      *hero.Player
		zone        *world.Zone
		kick        <-chan string
		hud         <-chan hero.HUD
		inventory   <-chan []world.Visible
		messages    <-chan []hero.Message
		updateQueue packetUpdate
	)
	updateTick := time.NewTicker(time.Second / 10)
	defer updateTick.Stop()
	updates := make(chan []packetUpdateUpdate, 1)

	sendUpdates := func(newUpdates ...packetUpdateUpdate) {
		for {
			select {
			case updates <- newUpdates:
				return
			case other := <-updates:
				newUpdates = append(other, newUpdates...)
			}
		}
	}
	listener := &world.ZoneListener{
		Add: func(t *world.Tile, obj world.ObjectLike) {
			o, ok := obj.(world.Visible)
			if !ok {
				return
			}
			x, y := t.Position()
			sendUpdates(packetUpdateUpdate{
				ID:     o.NetworkID(),
				X:      x,
				Y:      y,
				Object: addSprites(&packetUpdateObject{Name: o.Name()}, o),
			})
		},
		Remove: func(t *world.Tile, obj world.ObjectLike) {
			o, ok := obj.(world.Visible)
			if !ok {
				return
			}
			x, y := t.Position()
			sendUpdates(packetUpdateUpdate{
				ID:     o.NetworkID(),
				X:      x,
				Y:      y,
				Remove: true,
			})
		},
		Move: func(from, to *world.Tile, obj world.ObjectLike) {
			o, ok := obj.(world.Visible)
			if !ok {
				return
			}
			fx, fy := from.Position()
			tx, ty := to.Position()
			sendUpdates(packetUpdateUpdate{
				ID:    o.NetworkID(),
				X:     tx,
				Y:     ty,
				FromX: &fx,
				FromY: &fy,
			})
		},
		Update: func(t *world.Tile, obj world.ObjectLike) {
			o, ok := obj.(world.Visible)
			if !ok {
				return
			}
			x, y := t.Position()
			sendUpdates(packetUpdateUpdate{
				ID:     o.NetworkID(),
				X:      x,
				Y:      y,
				Object: addSprites(&packetUpdateObject{Name: o.Name()}, o),
			})
		},
	}

	for {
		select {
		case packet, ok := <-packets:
			if !ok {
				return
			}
			if packet.Auth != nil {
				if player != nil {
					return
				}

				var err string
				player, err = hero.Login(addr, packet.Auth)
				if err != "" {
					websocket.JSON.Send(ws, packetKick{err})
					return
				}
				world.InitObject(player)
				kick, hud, inventory, messages = player.InitPlayer()
				zx, tx, zy, ty := player.LoginPosition()
				defer hero.PlayerDisconnected(player)

				zone = world.GetZone(zx, zy)
				zone.AddListener(listener)
				if player.CanSpawn() {
					zone.Tile(tx, ty).Add(player)
				} else {
					player.CharacterCreation("")
				}
				defer func() {
					if t := player.Position(); t != nil {
						t.Remove(player)
					}
					zone.RemoveListener(listener)
					world.ReleaseZone(zone)
				}()

				inventoryObjects := []*packetUpdateObject{}
				for _, item := range player.Inventory() {
					inventoryObjects = append(inventoryObjects, addSprites(&packetUpdateObject{Name: item.Name()}, item))
				}
				websocket.JSON.Send(ws, packetInventory{inventoryObjects})
			}
			if player == nil {
				continue
			}
			if packet.Walk != nil {
				if t := player.Position(); t != nil {
					x, y := t.Position()
					if packet.Walk.X == x && packet.Walk.Y == y {
						player.ClearSchedule()
					} else {
						player.SetSchedule(world.NewWalkSchedule(x, y, packet.Walk.X, packet.Walk.Y))
					}
				}
			}
			if packet.HUD != nil {
				switch packet.HUD.Name {
				default:
					return

				case "cc":
					if cmd, ok := packet.HUD.Data.(string); ok {
						player.CharacterCreation(cmd)
					} else {
						return
					}
				}
			}
			if packet.Chat != nil {
				player.Chat(*packet.Chat)
			}
			if len(packet.Admin) > 0 {
				player.AdminCommand(addr, packet.Admin...)
			}

		case message := <-kick:
			websocket.JSON.Send(ws, packetKick{message})
			return

		case h := <-hud:
			websocket.JSON.Send(ws, packetHUD{h})

		case items := <-inventory:
			objects := make([]*packetUpdateObject, len(items))
			for i, item := range items {
				objects[i] = addSprites(&packetUpdateObject{Name: item.Name()}, item)
			}
			websocket.JSON.Send(ws, packetInventory{objects})

		case update := <-updates:
			updateQueue.Update = append(updateQueue.Update, update...)

		case <-updateTick.C:
			if len(updateQueue.Update) == 0 {
				continue
			}

			if player == nil {
				continue
			}

			leftover := updateQueue.Update[:0]
			if len(updateQueue.Update) > 100 {
				updateQueue.Update, leftover = updateQueue.Update[:100], updateQueue.Update[100:]
			}

			if t := player.Position(); t != nil {
				updateQueue.PlayerX, updateQueue.PlayerY = t.Position()
			} else {
				updateQueue.PlayerX, updateQueue.PlayerY = 127, 127
			}

			websocket.JSON.Send(ws, updateQueue)

			updateQueue.Update = leftover

		case msg := <-messages:
			websocket.JSON.Send(ws, packetMessage{msg})
		}
	}
}