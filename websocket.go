package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/BenLubar/Rnoadm/hero"
	"io"
	"log"
	"net"
)

type clientPacket struct {
	Admin *string
	Auth  *hero.LoginPacket
}

type packetKick struct {
	Kick string
}

var packetClientHash struct {
	ClientHash uint64 `json:",string"`
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
		player *hero.Player
		kick   <-chan string
	)

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
				kick = player.InitPlayer()
				defer hero.PlayerDisconnected(player)
			}
			log.Printf("%s: %+v", addr, packet)

		case message := <-kick:
			websocket.JSON.Send(ws, packetKick{message})
			return
		}
	}
}
