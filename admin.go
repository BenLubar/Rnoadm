package main

import (
	"log"
	"math/rand"
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

	parts := strings.Split(strings.ToLower(cmd), " ")
	if len(parts) < 1 {
		return
	}
	switch parts[0] {
	case "equip":
		objects := make([]Object, len(p.Backpack))
		copy(objects, p.Backpack)
		for _, o := range objects {
			p.Equip(o, true)
		}

	case "clear":
		if len(parts) < 2 {
			return
		}

		switch parts[1] {
		case "inventory":
			p.Backpack = nil
			select {
			case p.backpackDirty <- struct{}{}:
			default:
			}
		}
	case "spawn":
		AdminCommandSpawn(addr, p, parts[1:])
	}
}

func ParseMetalType(metal string) (MetalType, bool) {
	for t, m := range metalTypeInfo {
		if t == 0 {
			continue
		}
		if m.Name == metal {
			return MetalType(t), true
		}
	}
	return 0, false
}

func ParseRockType(rock string) (RockType, bool) {
	for t, r := range rockTypeInfo {
		if r.Name == rock {
			return RockType(t), true
		}
	}
	return 0, false
}

func ParseWoodType(wood string) (WoodType, bool) {
	for t, w := range woodTypeInfo {
		if w.Name == wood {
			return WoodType(t), true
		}
	}
	return 0, false
}

func ParseCosmeticName(name string) (CosmeticType, uint64, bool) {
	for t := range cosmetics {
		for i := range cosmetics[t] {
			if cosmetics[t][i].Name == name {
				return CosmeticType(t), uint64(i), true
			}
		}
	}
	return 0, 0, false
}

func AdminCommandSpawn(addr string, p *Player, cmd []string) {
	if len(cmd) < 1 {
		return
	}
	switch cmd[0] {
	case "human":
		h := GenerateHero(Human, rand.New(&p.Seed))
		p.GiveItem(h)
		p.SendMessage("spawned " + h.Name())
	default:
		var (
			metal   MetalType
			metalOk bool

			stone   RockType
			stoneOk bool

			wood   WoodType
			woodOk bool
		)
		var leftover string
		for i, c := range cmd {
			switch c {
			case "", "and", "an", "a", "with":
				continue
			}
			if !metalOk {
				metal, metalOk = ParseMetalType(c)
				if metalOk {
					continue
				}
			}
			if !stoneOk {
				stone, stoneOk = ParseRockType(c)
				if stoneOk {
					continue
				}
			}

			leftover = strings.Join(cmd[i:], " ")
			break
		}
		var item Object
		switch leftover {
		case "ore":
			if metalOk && !stoneOk && !woodOk {
				item = &Ore{
					Type: metal,
				}
			}
		case "stone":
			if !metalOk && stoneOk && !woodOk {
				item = &Stone{
					Type: stone,
				}
			}
		case "rock":
			if stoneOk && !woodOk {
				item = &Rock{
					Type: stone,
					Ore:  metal,
				}
			}
		case "log", "logs":
			if !metalOk && !stoneOk && woodOk {
				item = &Logs{
					Type: wood,
				}
			}
		case "tree":
			if !metalOk && !stoneOk && woodOk {
				item = &Tree{
					Type: wood,
				}
			}
		}
		if item == nil && !stoneOk && !woodOk {
			t, i, ok := ParseCosmeticName(leftover)
			if ok {
				c := &Cosmetic{
					Type:  t,
					ID:    i,
					Metal: metal,
				}
				metalExpected := false
				for _, color := range cosmetics[t][i].Colors {
					if color == metalColor {
						metalExpected = true
						break
					}
				}
				if metalOk != metalExpected {
					return
				}
				item = c
			}
		}
		if item != nil {
			p.GiveItem(item)
			p.SendMessage("spawned " + item.Name())
		}
	}
}
