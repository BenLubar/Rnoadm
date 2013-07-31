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
		p.hud = &AdminChangeExamineHUD{Player: p, Input: []rune(p.Examine())}
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
	adminCommands["SPAWN HERO"] = func(p *Player) {
		p.Lock()
		p.GiveItem(&Hero{Name_: GenerateName(rand.New(rand.NewSource(rand.Int63())), NameHero)})
		p.Unlock()
	}
}

type AdminHUD struct {
	Player *Player
	Input  []rune
	Key_   bool
}

func (h *AdminHUD) Paint(setcell func(int, int, string, string, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	setcell(0, 0, ">", "", "#00f")
	setcell(1, 0, string(h.Input), "", "#00f")
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

func (h *AdminTeleportHUD) Paint(setcell func(int, int, string, string, Color)) {
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
		setcell(0, 0, "SUMMON PLAYER", "", "#fff")
	} else {
		setcell(0, 0, "TELEPORT TO PLAYER", "", "#fff")
	}
	for i, p := range h.List[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#fff")
		id := p.ID
		var idBuf [16]byte
		for k := range idBuf {
			idBuf[k] = "0123456789ABCDEF"[id&15]
			id >>= 4
		}
		setcell(1, i+1, string(idBuf[:]), "", "#44f")
		p.Lock()
		setcell(7, i+1, p.Name(), "", "#00f")
		p.Unlock()
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(1, 9, "previous", "", "#fff")
	}
	if len(h.List) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(1, 10, "next", "", "#fff")
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

func (h *AdminTeleportZoneHUD) Paint(setcell func(int, int, string, string, Color)) {
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

	setcell(0, 0, "TELEPORT TO ZONE", "", "#00f")
	setcell(0, 1, "X", "", "#00f")
	setcell(1, 1, strconv.FormatInt(h.X, 10), "", "#0ff")
	setcell(0, 2, "Y", "", "#00f")
	setcell(1, 2, strconv.FormatInt(h.Y, 10), "", "#0ff")
	setcell(1, 3, h.Name, "", "#0ff")
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

func (h *AdminNoclipHUD) Paint(setcell func(int, int, string, string, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return
	}

	setcell(0, 0, "NOCLIP", "", "#00f")
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

type AdminChangeExamineHUD struct {
	Player *Player
	Input  []rune
}

func (h *AdminChangeExamineHUD) Paint(setcell func(int, int, string, string, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, strings.ToUpper(h.Player.Name()), "", "#fff")
	h.Player.Unlock()

	setcell(0, 1, ">", "", "#00f")
	setcell(1, 1, string(h.Input), "", "#00f")
}

func (h *AdminChangeExamineHUD) Key(code int, special bool) bool {
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
		AdminLog.Printf("CHANGEEXAMINE:%q [%d:%q] (%d:%d, %d:%d)", string(h.Input), h.Player.ID, h.Player.Name(), h.Player.ZoneX, h.Player.TileX, h.Player.ZoneY, h.Player.TileY)
		h.Player.Examine_ = strings.TrimSpace(strings.ToLower(string(h.Input)))
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

func (h *AdminChangeExamineHUD) Click(x, y int) bool {
	return false
}

type AdminChangeColorHUD struct {
	Player *Player
	Var    *Color
	Input  []rune
}

func (h *AdminChangeColorHUD) Paint(setcell func(int, int, string, string, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, "CHANGE COLOR", "", "#00f")
	h.Player.Unlock()

	setcell(0, 1, ">", "", "#00f")
	setcell(1, 1, string(h.Input), "", "#0ff")
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

func (h *AdminMenuHUD) Paint(setcell func(int, int, string, string, Color)) {
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

	setcell(0, 0, "ADMIN MENU-O-MATIC", "", "#00f")
	for i, c := range h.Commands[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#0ff")
		setcell(1, i+1, c, "", "#0ff")
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(1, 9, "previous", "", "#fff")
	}
	if len(h.Commands) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(1, 10, "next", "", "#fff")
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
