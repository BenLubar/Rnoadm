package main

import (
	"log"
)

var AdminLog *log.Logger

func AdminCommand(addr string, p *Player, cmd string) {
	if !p.Admin {
		p.Kick("I'm sorry, Dave. I'm afraid you can't do that.")
		return
	}
	p.Lock()
	defer p.Unlock()

	AdminLog.Printf("[%s:%d] %q (%d:%d, %d:%d) COMMAND: %q", addr, p.ID, p.Name(), p.ZoneX, p.TileX, p.ZoneY, p.TileY, cmd)
}
