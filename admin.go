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

	case "health":
		p.SendMessage(strconv.FormatUint(p.Health(), 10) + "/" + strconv.FormatUint(p.MaxHealth(), 10))

	case "clear":
		if len(parts) < 2 {
			return
		}

		switch parts[1] {
		case "inventory":
			p.Backpack = nil
		case "hat", "helmet", "headwear", "headgear":
			p.Worn[Headwear] = Cosmetic{}
		case "shirt":
			p.Worn[Shirt] = Cosmetic{}
		case "pants":
			p.Worn[Pants] = Cosmetic{}
		case "shoes", "boots":
			p.Worn[Shoes] = Cosmetic{}
		case "pauldrons":
			p.Worn[Pauldrons] = Cosmetic{}
		case "breastplate":
			p.Worn[Breastplate] = Cosmetic{}
		case "vambraces":
			p.Worn[Vambraces] = Cosmetic{}
		case "gauntlets":
			p.Worn[Gauntlets] = Cosmetic{}
		case "tassets":
			p.Worn[Tassets] = Cosmetic{}
		case "greaves":
			p.Worn[Greaves] = Cosmetic{}
		case "armor":
			p.Worn[Headwear] = Cosmetic{}
			p.Worn[Shoes] = Cosmetic{}
			p.Worn[Pauldrons] = Cosmetic{}
			p.Worn[Breastplate] = Cosmetic{}
			p.Worn[Vambraces] = Cosmetic{}
			p.Worn[Gauntlets] = Cosmetic{}
			p.Worn[Tassets] = Cosmetic{}
			p.Worn[Greaves] = Cosmetic{}
		}
		select {
		case p.backpackDirty <- struct{}{}:
		default:
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
			if !woodOk {
				wood, woodOk = ParseWoodType(c)
				if woodOk {
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
		case "rich rock":
			if metalOk && stoneOk && !woodOk {
				item = &Rock{
					Type: stone,
					Ore:  metal,
					Big:  true,
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
		case "pickaxe":
			if metalOk && !stoneOk && woodOk {
				item = &Pickaxe{
					Handle: wood,
					Head:   metal,
				}
			}
		case "hatchet":
			if metalOk && !stoneOk && woodOk {
				item = &Hatchet{
					Handle: wood,
					Head:   metal,
				}
			}
		case "armor set":
			if metalOk && !stoneOk && !woodOk {
				cmd = cmd[:len(cmd)-2]
				AdminCommandSpawn(addr, p, append(cmd, "helmet"))
				AdminCommandSpawn(addr, p, append(cmd, "breastplate"))
				AdminCommandSpawn(addr, p, append(cmd, "pauldrons"))
				AdminCommandSpawn(addr, p, append(cmd, "vambraces"))
				AdminCommandSpawn(addr, p, append(cmd, "gauntlets"))
				AdminCommandSpawn(addr, p, append(cmd, "tassets"))
				AdminCommandSpawn(addr, p, append(cmd, "greaves"))
				AdminCommandSpawn(addr, p, append(cmd, "boots"))
				return
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
