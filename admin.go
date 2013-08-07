package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
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

	parts := strings.Split(cmd, " ")
	if len(parts) < 1 {
		return
	}
	switch strings.ToLower(parts[0]) {
	case "spawn":
		if len(parts) < 2 {
			return
		}
		switch strings.ToLower(parts[1]) {
		case "human":
			p.GiveItem(GenerateHero(Human, rand.New(&p.Seed)))
		case "headwear", "shirt", "pants", "shoes", "breastplate", "pauldrons", "vambraces", "gauntlets", "tassets", "greaves":
			if len(parts) < 3 {
				return
			}
			t := map[string]CosmeticType{
				"headwear":    Headwear,
				"shirt":       Shirt,
				"pants":       Pants,
				"shoes":       Shoes,
				"breastplate": Breastplate,
				"pauldrons":   Pauldrons,
				"vambraces":   Vambraces,
				"gauntlets":   Gauntlets,
				"tassets":     Tassets,
				"greaves":     Greaves,
			}[strings.ToLower(parts[1])]
			i, err := strconv.ParseUint(parts[2], 0, 64)
			if err != nil {
				p.SendMessage(err.Error())
				return
			}
			if len(parts) > 3 {
				// TODO: materials
				return
			}
			c := &Cosmetic{
				Type: t,
				ID:   i,
			}
			if !c.Exists() {
				p.SendMessage("nonexistent item.")
				return
			}
			p.GiveItem(c)
			p.SendMessage("spawned " + c.Name())
		}
	}
}
