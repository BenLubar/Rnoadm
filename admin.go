package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var AdminLog *log.Logger

var adminCommands = map[string]func(*Player){
	"MENU": func(p *Player) {
		p.hud = &AdminMenuHUD{Player: p}
	},
	"TP": func(p *Player) {
		p.hud = &AdminTeleportHUD{Player: p, Summon: false}
	},
	"SUMMON": func(p *Player) {
		p.hud = &AdminTeleportHUD{Player: p, Summon: true}
	},
	"TZ": func(p *Player) {
		p.Lock()
		p.hud = &AdminTeleportZoneHUD{Player: p, X: p.ZoneX, Y: p.ZoneY}
		p.Unlock()
	},
	"NOCLIP": func(p *Player) {
		p.hud = &AdminNoclipHUD{Player: p}
	},
	"CHANGE EXAMINE": func(p *Player) {
		p.hud = &AdminChangeTextHUD{Player: p, Var: &p.Examine_, Input: []rune(p.Examine_)}
	},
	"CHANGE NICKNAME": func(p *Player) {
		p.hud = &AdminChangeTextHUD{Player: p, Var: &p.HeroName.Nickname, Input: []rune(p.HeroName.Nickname)}
	},
	"CHANGE NAME": func(p *Player) {
		p.Lock()
		name := *p.HeroName
		p.Unlock()
		p.hud = &AdminChangeNameHUD{Player: p, Name: &name}
	},
	"CHANGE SKIN COLOR": func(p *Player) {
		p.Lock()
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.CustomColor, Input: []rune(string(p.CustomColor))}
		p.Unlock()
	},
	"CHANGE HAT BASE COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Head == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Head.CustomColor[0], Input: []rune(string(p.Head.CustomColor[0]))}
	},
	"CHANGE HAT FIRST COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Head == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Head.CustomColor[1], Input: []rune(string(p.Head.CustomColor[1]))}
	},
	"CHANGE HAT SECOND COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Head == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Head.CustomColor[2], Input: []rune(string(p.Head.CustomColor[2]))}
	},
	"CHANGE HAT THIRD COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Head == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Head.CustomColor[3], Input: []rune(string(p.Head.CustomColor[3]))}
	},
	"CHANGE HAT FOURTH COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Head == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Head.CustomColor[4], Input: []rune(string(p.Head.CustomColor[4]))}
	},
	"CHANGE SHIRT BASE COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Top == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Top.CustomColor[0], Input: []rune(string(p.Top.CustomColor[0]))}
	},
	"CHANGE SHIRT FIRST COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Top == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Top.CustomColor[1], Input: []rune(string(p.Top.CustomColor[1]))}
	},
	"CHANGE SHIRT SECOND COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Top == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Top.CustomColor[2], Input: []rune(string(p.Top.CustomColor[2]))}
	},
	"CHANGE SHIRT THIRD COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Top == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Top.CustomColor[3], Input: []rune(string(p.Top.CustomColor[3]))}
	},
	"CHANGE SHIRT FOURTH COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Top == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Top.CustomColor[4], Input: []rune(string(p.Top.CustomColor[4]))}
	},
	"CHANGE PANTS BASE COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Legs == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Legs.CustomColor[0], Input: []rune(string(p.Legs.CustomColor[0]))}
	},
	"CHANGE PANTS FIRST COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Legs == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Legs.CustomColor[1], Input: []rune(string(p.Legs.CustomColor[1]))}
	},
	"CHANGE PANTS SECOND COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Legs == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Legs.CustomColor[2], Input: []rune(string(p.Legs.CustomColor[2]))}
	},
	"CHANGE PANTS THIRD COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Legs == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Legs.CustomColor[3], Input: []rune(string(p.Legs.CustomColor[3]))}
	},
	"CHANGE PANTS FOURTH COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Legs == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Legs.CustomColor[4], Input: []rune(string(p.Legs.CustomColor[4]))}
	},
	"CHANGE SHOES BASE COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Feet == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Feet.CustomColor[0], Input: []rune(string(p.Feet.CustomColor[0]))}
	},
	"CHANGE SHOES FIRST COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Feet == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Feet.CustomColor[1], Input: []rune(string(p.Feet.CustomColor[1]))}
	},
	"CHANGE SHOES SECOND COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Feet == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Feet.CustomColor[2], Input: []rune(string(p.Feet.CustomColor[2]))}
	},
	"CHANGE SHOES THIRD COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Feet == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Feet.CustomColor[3], Input: []rune(string(p.Feet.CustomColor[3]))}
	},
	"CHANGE SHOES FOURTH COLOR": func(p *Player) {
		p.Lock()
		defer p.Unlock()
		if p.Feet == nil {
			return
		}
		p.hud = &AdminChangeColorHUD{Player: p, Var: &p.Feet.CustomColor[4], Input: []rune(string(p.Feet.CustomColor[4]))}
	},
	"DELETE THE ENTIRE ZONE": func(p *Player) {
		p.Lock()
		z := p.zone
		p.Unlock()

		z.Lock()
		for i := range z.Tiles {
			t := &z.Tiles[i]
			for j := 0; j < len(t.Objects); j++ {
				if _, ok := t.Objects[j].(*Player); !ok {
					t.Objects = append(t.Objects[:j], t.Objects[j+1:]...)
					j--
				}
			}
			if len(t.Objects) == 0 {
				t.Objects = nil
			}
		}
		z.Unlock()
		z.Save()
		z.Repaint()
	},
	"REGENERATE THE ENTIRE ZONE": func(p *Player) {
		p.Lock()
		z := p.zone
		p.Unlock()

		z.Lock()
		for i := range z.Tiles {
			t := &z.Tiles[i]
			for j := 0; j < len(t.Objects); j++ {
				if _, ok := t.Objects[j].(*Player); !ok {
					t.Objects = append(t.Objects[:j], t.Objects[j+1:]...)
					j--
				}
			}
			if len(t.Objects) == 0 {
				t.Objects = nil
			}
		}
		z.Unlock()
		z.Generate()
		z.Save()
		z.Repaint()
	},
}

func init() {
	for t := range rockTypeInfo {
		rt := RockType(t)
		adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" ROCK"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Rock{Type: rt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" STONE"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Stone{Type: rt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" WALL"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&WallStone{Type: rt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" FLOOR"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&FloorStone{Type: rt})
			p.Unlock()
		}
		for m := range metalTypeInfo {
			if m == 0 {
				continue
			}
			mt := MetalType(m)
			adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" ROCK WITH "+strings.ToUpper(metalTypeInfo[mt].Name)+" ORE"] = func(p *Player) {
				p.Lock()
				p.GiveItem(&Rock{Type: rt, Ore: mt})
				p.Unlock()
			}
			adminCommands["SPAWN "+strings.ToUpper(rockTypeInfo[rt].Name)+" ROCK WITH RICH "+strings.ToUpper(metalTypeInfo[mt].Name)+" ORE"] = func(p *Player) {
				p.Lock()
				p.GiveItem(&Rock{Type: rt, Ore: mt, Big: true})
				p.Unlock()
			}
		}
	}
	for t := range metalTypeInfo {
		if t == 0 {
			continue
		}
		mt := MetalType(t)
		adminCommands["SPAWN "+strings.ToUpper(metalTypeInfo[mt].Name)+" ORE"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Ore{Type: mt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(metalTypeInfo[mt].Name)+" WALL"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&WallMetal{Type: mt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(metalTypeInfo[mt].Name)+" FLOOR"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&FloorMetal{Type: mt})
			p.Unlock()
		}
		for w := range woodTypeInfo {
			wt := WoodType(w)
			adminCommands["SPAWN "+strings.ToUpper(metalTypeInfo[mt].Name)+" AND "+strings.ToUpper(woodTypeInfo[wt].Name)+" PICKAXE"] = func(p *Player) {
				p.Lock()
				p.GiveItem(&Pickaxe{Head: mt, Handle: wt})
				p.Unlock()
			}
			adminCommands["SPAWN "+strings.ToUpper(metalTypeInfo[mt].Name)+" AND "+strings.ToUpper(woodTypeInfo[wt].Name)+" HATCHET"] = func(p *Player) {
				p.Lock()
				p.GiveItem(&Hatchet{Head: mt, Handle: wt})
				p.Unlock()
			}
		}
	}
	for t := range woodTypeInfo {
		wt := WoodType(t)
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" TREE"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Tree{Type: wt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" LOGS"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Logs{Type: wt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" WALL"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&WallWood{Type: wt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" FLOOR"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&FloorWood{Type: wt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" BED"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Bed{Frame: wt})
			p.Unlock()
		}
		adminCommands["SPAWN "+strings.ToUpper(woodTypeInfo[wt].Name)+" CHEST"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Chest{Type: wt})
			p.Unlock()
		}
	}
	for t := range floraTypeInfo {
		ft := FloraType(t)
		adminCommands["SPAWN "+strings.ToUpper(floraTypeInfo[ft].Name)+" PLANT"] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Flora{Type: ft})
			p.Unlock()
		}
	}
	for t := range hatTypeInfo {
		ht := HatType(t)
		adminCommands["SPAWN "+strings.ToUpper(hatTypeInfo[ht].Name)] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Hat{Type: ht})
			p.Unlock()
		}
	}
	for t := range shirtTypeInfo {
		st := ShirtType(t)
		adminCommands["SPAWN "+strings.ToUpper(shirtTypeInfo[st].Name)] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Shirt{Type: st})
			p.Unlock()
		}
	}
	for t := range pantsTypeInfo {
		pt := PantsType(t)
		adminCommands["SPAWN "+strings.ToUpper(pantsTypeInfo[pt].Name)] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Pants{Type: pt})
			p.Unlock()
		}
	}
	for t := range shoeTypeInfo {
		st := ShoeType(t)
		adminCommands["SPAWN "+strings.ToUpper(shoeTypeInfo[st].Name)] = func(p *Player) {
			p.Lock()
			p.GiveItem(&Shoes{Type: st})
			p.Unlock()
		}
	}
	for t := range raceInfo {
		rt := Race(t)
		adminCommands["SPAWN "+strings.ToUpper(raceInfo[rt].Name)] = func(p *Player) {
			p.Lock()
			p.GiveItem(GenerateHero(rt, rand.New(rand.NewSource(rand.Int63()))))
			p.Unlock()
		}
	}
}

type AdminHUD struct {
	Player *Player
	Input  []rune
	Key_   bool
}

func (h *AdminHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	setcell(0, 0, PaintCell{
		Text:   ">",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 0, PaintCell{
		Text:   string(h.Input),
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		if code != 0 && h.Key_ {
			h.Key_ = false
			h.Input = append(h.Input, unicode.ToUpper(rune(code)))
			h.Player.Repaint()
		}
		return true
	}
	switch code {
	case 37, 38, 39, 40: // arrow keys
		return false
	case 9, 16, 17, 18: // tab shift ctrl alt
		return false
	case 8: // backspace
		if len(h.Input) > 0 {
			h.Input = h.Input[:len(h.Input)-1]
			h.Player.Repaint()
		}
		return true
	case 13: // enter
		if f, ok := adminCommands[string(h.Input)]; ok {
			h.Player.Lock()
			AdminLog.Printf("COMMAND:%q [%d:%q] (%d:%d, %d:%d)", string(h.Input), h.Player.ID, h.Player.Name(), h.Player.ZoneX, h.Player.TileX, h.Player.ZoneY, h.Player.TileY)
			h.Player.Unlock()

			h.Player.hud = nil
			f(h.Player)
			h.Player.Repaint()
		}
		return true
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	h.Key_ = true
	return true
}

func (h *AdminHUD) Click(x, y int) bool {
	return false
}

type PlayerList []*Player

func (l PlayerList) Len() int {
	return len(l)
}
func (l PlayerList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l PlayerList) Less(i, j int) bool {
	return l[i].ID < l[j].ID
}

type AdminTeleportHUD struct {
	Player *Player
	List   PlayerList
	Offset int
	Summon bool
}

func (h *AdminTeleportHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return
	}
	if h.List == nil {
		onlinePlayersLock.Lock()
		for _, p := range OnlinePlayers {
			h.List = append(h.List, p)
		}
		onlinePlayersLock.Unlock()
		sort.Sort(h.List)
	}
	if h.Summon {
		setcell(0, 0, PaintCell{
			Text:   "SUMMON PLAYER",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	} else {
		setcell(0, 0, PaintCell{
			Text:   "TELEPORT TO PLAYER",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	}
	for i, p := range h.List[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, PaintCell{
			Text:   string(rune(i) + '1'),
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
		id := p.ID
		var idBuf [16]byte
		for k := range idBuf {
			idBuf[k] = "0123456789ABCDEF"[id&15]
			id >>= 4
		}
		setcell(1, i+1, PaintCell{
			Text:   string(idBuf[:]),
			Color:  "#44f",
			ZIndex: 1<<31 - 1,
		})
		p.Lock()
		setcell(7, i+1, PaintCell{
			Text:   p.Name(),
			Color:  "#00f",
			ZIndex: 1<<31 - 1,
		})
		p.Unlock()
	}
	if h.Offset > 0 {
		setcell(0, 9, PaintCell{
			Text:   "9",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
		setcell(1, 9, PaintCell{
			Text:   "previous",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	}
	if len(h.List) > h.Offset+8 {
		setcell(0, 10, PaintCell{
			Text:   "0",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
		setcell(1, 10, PaintCell{
			Text:   "next",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	}
}

func (h *AdminTeleportHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		return false
	}
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.List) {
			from := h.Player
			to := h.List[i]
			action := "TELEPORT"
			if h.Summon {
				from, to = to, from
				action = "SUMMON"
			}
			to.Lock()
			zx, zy := to.ZoneX, to.ZoneY
			tx, ty := to.TileX, to.TileY
			name := to.Name()
			to.Unlock()

			from.Lock()
			az := from.zone
			azx, azy := from.ZoneX, from.ZoneY
			atx, aty := from.TileX, from.TileY
			aname := from.Name()
			from.Unlock()

			AdminLog.Printf("%s [%d:%q] (%d:%d, %d:%d) => [%d:%q] (%d:%d, %d:%d)", action, from.ID, aname, azx, atx, azy, aty, to.ID, name, zx, tx, zy, ty)

			az.Lock()
			az.Tile(atx, aty).Remove(from)
			az.Unlock()
			az.Repaint()

			ReleaseZone(az)
			z := GrabZone(zx, zy)

			from.Lock()
			from.zone = z
			from.ZoneX, from.ZoneY = zx, zy
			from.TileX, from.TileY = tx, ty
			from.Unlock()

			z.Lock()
			z.Tile(tx, ty).Add(from)
			z.Repaint()
			z.Unlock()

			h.Player.hud = nil
			h.Player.Repaint()
		}
		return true
	case '9':
		if h.Offset > 0 {
			h.Offset--
			h.Player.Repaint()
		}
		return true
	case '0':
		if h.Offset+8 < len(h.List) {
			h.Offset++
			h.Player.Repaint()
		}
		return true

	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *AdminTeleportHUD) Click(x, y int) bool {
	return false
}

type AdminTeleportZoneHUD struct {
	Player *Player
	X, Y   int64
	Name   string
}

func (h *AdminTeleportZoneHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return
	}

	if h.Name == "" {
		loadedZoneLock.Lock()
		h.Name = ZoneInfo[[2]int64{h.X, h.Y}].Name
		loadedZoneLock.Unlock()
		if h.Name == "" {
			h.Name = "???"
		}
	}

	setcell(0, 0, PaintCell{
		Text:   "TELEPORT TO ZONE",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	setcell(0, 1, PaintCell{
		Text:   "X",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 1, PaintCell{
		Text:   strconv.FormatInt(h.X, 10),
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
	setcell(0, 2, PaintCell{
		Text:   "Y",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 2, PaintCell{
		Text:   strconv.FormatInt(h.Y, 10),
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 3, PaintCell{
		Text:   h.Name,
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminTeleportZoneHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		return false
	}
	switch code {
	case 38: // up
		h.Y--
		h.Name = ""
		h.Player.Repaint()
		return true
	case 37: // left
		h.X--
		h.Name = ""
		h.Player.Repaint()
		return true
	case 40: // down
		h.Y++
		h.Name = ""
		h.Player.Repaint()
		return true
	case 39: // right
		h.X++
		h.Name = ""
		h.Player.Repaint()
		return true

	case 13: // enter
		h.Player.Lock()
		tx, ty := h.Player.TileX, h.Player.TileY
		from := h.Player.zone
		AdminLog.Printf("TPZONE [%d:%q] (%d:%d, %d:%d) => (%d, %d)", h.Player.ID, h.Player.Name(), h.Player.ZoneX, tx, h.Player.ZoneY, ty, h.X, h.Y)
		h.Player.Unlock()

		from.Lock()
		if !from.Tile(tx, ty).Remove(h.Player) {
			from.Unlock()
			h.Player.hud = nil
			h.Player.Repaint()
			return true
		}
		from.Unlock()
		ReleaseZone(from)

		to := GrabZone(h.X, h.Y)
		to.Lock()
		to.Tile(tx, ty).Add(h.Player)
		to.Unlock()

		h.Player.Lock()
		h.Player.ZoneX, h.Player.ZoneY = h.X, h.Y
		h.Player.zone = to
		h.Player.Unlock()

		h.Player.hud = nil
		h.Player.Repaint()
		return true

	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *AdminTeleportZoneHUD) Click(x, y int) bool {
	return false
}

type AdminNoclipHUD struct {
	Player *Player
}

func (h *AdminNoclipHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return
	}

	setcell(0, 0, PaintCell{
		Text:   "NOCLIP",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminNoclipHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		return false
	}
	move := func(dx, dy int) {
		h.Player.Lock()
		zone := h.Player.zone
		fx, fy := h.Player.TileX, h.Player.TileY
		tx, ty := fx+uint8(dx), fy+uint8(dy)
		h.Player.Unlock()
		if int(tx) != int(fx)+dx || int(ty) != int(fy)+dy {
			return
		}
		to := zone.Tile(tx, ty)
		if to == nil {
			return
		}

		zone.Lock()
		if zone.Tile(fx, fy).Remove(h.Player) {
			to.Add(h.Player)
		}
		zone.Unlock()
		zone.Repaint()

		h.Player.Lock()
		h.Player.TileX, h.Player.TileY = tx, ty
		h.Player.Unlock()
	}
	switch code {
	case 38: // up
		move(0, -1)
		return true
	case 37: // left
		move(-1, 0)
		return true
	case 40: // down
		move(0, 1)
		return true
	case 39: // right
		move(1, 0)
		return true

	case 13, 27: // enter, esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *AdminNoclipHUD) Click(x, y int) bool {
	return false
}

type AdminChangeNameHUD struct {
	Player *Player
	Name   *HeroName
	Index  int
}

func (h *AdminChangeNameHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, PaintCell{
		Text:   strings.ToUpper(h.Player.Name()),
		Color:  "#fff",
		ZIndex: 1<<31 - 1,
	})
	h.Player.Unlock()

	setcell(0, 1, PaintCell{
		Text:   h.Name.Name(),
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	desc, subtype, index := h.index()
	setcell(0, 2, PaintCell{
		Text:   "< " + desc + " >",
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
	setcell(0, 3, PaintCell{
		Text:   "subtype [ " + strconv.FormatUint(uint64(*subtype), 10) + " ]",
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
	setcell(0, 4, PaintCell{
		Text:   "index - " + strconv.FormatUint(uint64(*index), 10) + " +",
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminChangeNameHUD) index() (string, *NameSubtype, *uint16) {
	switch h.Index {
	default:
		fallthrough
	case 0:
		return "First Name", &h.Name.FirstT, &h.Name.First
	case 1:
		return "Last Name", &h.Name.Last1T, &h.Name.Last1
	case 2:
		return "Last Suffix 1", &h.Name.Last2T, &h.Name.Last2
	case 3:
		return "Last Suffix 2", &h.Name.Last3T, &h.Name.Last3
	}
}

func (h *AdminChangeNameHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		return true
	}
	switch code {
	case 188: // ,
		if h.Index != 0 {
			h.Index--
		} else {
			h.Index = 3
		}
		h.Player.Repaint()
		return true
	case 190: // .
		if h.Index != 3 {
			h.Index++
		} else {
			h.Index = 0
		}
		h.Player.Repaint()
		return true
	case 219: // [
		_, subtype, index := h.index()
		if *subtype != 0 {
			*subtype--
		} else {
			*subtype = NameSubtype(len(names) - 1)
		}
		*index = 0
		h.Player.Repaint()
		return true
	case 221: // ]
		_, subtype, index := h.index()
		if *subtype != NameSubtype(len(names)-1) {
			*subtype++
		} else {
			*subtype = 0
		}
		*index = 0
		h.Player.Repaint()
		return true
	case 189: // -
		_, subtype, index := h.index()
		if *index != 0 {
			*index--
		} else {
			*index = uint16(len(names[*subtype]) - 1)
		}
		h.Player.Repaint()
		return true
	case 187: // +
		_, subtype, index := h.index()
		if *index != uint16(len(names[*subtype])-1) {
			*index++
		} else {
			*index = 0
		}
		h.Player.Repaint()
		return true
	case 13: // enter
		h.Player.Lock()
		h.Player.HeroName = h.Name
		h.Player.Unlock()

		h.Player.hud = nil
		h.Player.Repaint()
		return true
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *AdminChangeNameHUD) Click(x, y int) bool {
	return false
}

type AdminChangeTextHUD struct {
	Player *Player
	Var    *string
	Input  []rune
}

func (h *AdminChangeTextHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, PaintCell{
		Text:   strings.ToUpper(h.Player.Name()),
		Color:  "#fff",
		ZIndex: 1<<31 - 1,
	})
	h.Player.Unlock()

	setcell(0, 1, PaintCell{
		Text:   ">",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 1, PaintCell{
		Text:   string(h.Input),
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminChangeTextHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		if code != 0 {
			h.Input = append(h.Input, rune(code))
			h.Player.Repaint()
		}
		return true
	}
	switch code {
	case 37, 38, 39, 40: // arrow keys
		return false
	case 9, 16, 17, 18: // tab shift ctrl alt
		return false
	case 8: // backspace
		if len(h.Input) > 0 {
			h.Input = h.Input[:len(h.Input)-1]
			h.Player.Repaint()
		}
		return true
	case 13: // enter
		h.Player.Lock()
		*h.Var = strings.TrimSpace(string(h.Input))
		h.Player.Unlock()

		h.Player.hud = nil
		h.Player.Repaint()
		return true
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return true
}

func (h *AdminChangeTextHUD) Click(x, y int) bool {
	return false
}

type AdminChangeColorHUD struct {
	Player *Player
	Var    *Color
	Input  []rune
}

func (h *AdminChangeColorHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, PaintCell{
		Text:   "CHANGE COLOR",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	h.Player.Unlock()

	setcell(0, 1, PaintCell{
		Text:   ">",
		Color:  Color(h.Input),
		ZIndex: 1<<31 - 1,
	})
	setcell(1, 1, PaintCell{
		Text:   string(h.Input),
		Color:  "#0ff",
		ZIndex: 1<<31 - 1,
	})
}

func (h *AdminChangeColorHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	if !special {
		if code != 0 {
			h.Input = append(h.Input, rune(code))
			h.Player.Repaint()
		}
		return true
	}
	switch code {
	case 37, 38, 39, 40: // arrow keys
		return false
	case 9, 16, 17, 18: // tab shift ctrl alt
		return false
	case 8: // backspace
		if len(h.Input) > 0 {
			h.Input = h.Input[:len(h.Input)-1]
			h.Player.Repaint()
		}
		return true
	case 13: // enter
		h.Player.Lock()
		*h.Var = Color(strings.TrimSpace(string(h.Input)))
		h.Player.Unlock()

		h.Player.hud = nil
		h.Player.Repaint()
		return true
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return true
}

func (h *AdminChangeColorHUD) Click(x, y int) bool {
	return false
}

type AdminMenuHUD struct {
	Player   *Player
	Commands []string
	Offset   int
}

func (h *AdminMenuHUD) Paint(setcell func(int, int, PaintCell)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return
	}

	if h.Commands == nil {
		h.Commands = make([]string, 0, len(adminCommands))
		for c := range adminCommands {
			h.Commands = append(h.Commands, c)
		}
		sort.Strings(h.Commands)
	}

	setcell(0, 0, PaintCell{
		Text:   "ADMIN MENU-O-MATIC",
		Color:  "#00f",
		ZIndex: 1<<31 - 1,
	})
	for i, c := range h.Commands[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, PaintCell{
			Text:   string(rune(i) + '1'),
			Color:  "#0ff",
			ZIndex: 1<<31 - 1,
		})
		setcell(1, i+1, PaintCell{
			Text:   c,
			Color:  "#0ff",
			ZIndex: 1<<31 - 1,
		})
	}
	if h.Offset > 0 {
		setcell(0, 9, PaintCell{
			Text:   "9",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
		setcell(1, 9, PaintCell{
			Text:   "previous",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	}
	if len(h.Commands) > h.Offset+8 {
		setcell(0, 10, PaintCell{
			Text:   "0",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
		setcell(1, 10, PaintCell{
			Text:   "next",
			Color:  "#fff",
			ZIndex: 1<<31 - 1,
		})
	}
}

func (h *AdminMenuHUD) Key(code int, special bool) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}

	if !special {
		return false
	}
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.Commands) {
			h.Player.hud = &AdminHUD{
				Player: h.Player,
				Input:  []rune(h.Commands[i]),
			}
			h.Player.Repaint()
		}
		return true
	case '9':
		if h.Offset > 0 {
			h.Offset--
			h.Player.Repaint()
		}
		return true
	case '0':
		if h.Offset+8 < len(h.Commands) {
			h.Offset++
			h.Player.Repaint()
		}
		return true

	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *AdminMenuHUD) Click(x, y int) bool {
	return false
}
