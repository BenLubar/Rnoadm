package main

import (
	"log"
	"math/rand"
	"sort"
	"strings"
	"unicode"
)

var AdminLog *log.Logger

var adminCommands = map[string]func(*Player){
	"MENU": func(p *Player) {
		p.hud = &AdminMenuHUD{Player: p}
	},
	"TP": func(p *Player) {
		p.hud = &AdminTeleportHUD{Player: p}
	},
	"TELEPORT": func(p *Player) {
		p.hud = &AdminTeleportHUD{Player: p}
	},
	"CHANGE EXAMINE": func(p *Player) {
		p.hud = &AdminChangeExamineHUD{Player: p, Input: []rune(p.Examine())}
	},
	"CHANGE SKIN COLOR": func(p *Player) {
		p.Lock()
		p.hud = &AdminChangeColorHUD{Player: p, Input: []rune(string(p.BaseColor))}
		p.Unlock()
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
	setcell(2, 0, string(h.Input), "", "#00f")
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
	setcell(0, 0, "TELEPORT TO PLAYER", "", "#fff")
	for i, p := range h.List[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#fff")
		id := p.ID
		for k := 0; k < 16; k++ {
			setcell(17-k, i+1, string(rune("0123456789ABCDEF"[id&15])), "", "#44f")
			id >>= 4
		}
		p.Lock()
		setcell(19, i+1, p.Name(), "", "#00f")
		p.Unlock()
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(2, 9, "previous", "", "#fff")
	}
	if len(h.List) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(2, 10, "next", "", "#fff")
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
			p := h.List[i]
			p.Lock()
			zx, zy := p.ZoneX, p.ZoneY
			tx, ty := p.TileX, p.TileY
			name := p.Name()
			p.Unlock()

			h.Player.Lock()
			az := h.Player.zone
			azx, azy := h.Player.ZoneX, h.Player.ZoneY
			atx, aty := h.Player.TileX, h.Player.TileY
			aname := h.Player.Name()
			h.Player.Unlock()

			AdminLog.Printf("TELEPORT [%d:%q] (%d:%d, %d:%d) => [%d:%q] (%d:%d, %d:%d)", h.Player.ID, aname, azx, atx, azy, aty, p.ID, name, zx, tx, zy, ty)

			az.Lock()
			az.Tile(atx, aty).Remove(h.Player)
			az.Repaint()
			az.Unlock()

			ReleaseZone(az)
			z := GrabZone(zx, zy)

			h.Player.Lock()
			h.Player.zone = z
			h.Player.ZoneX, h.Player.ZoneY = zx, zy
			h.Player.TileX, h.Player.TileY = tx, ty
			h.Player.Unlock()
			h.Player.Save()

			z.Lock()
			z.Tile(tx, ty).Add(h.Player)
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
	setcell(2, 1, string(h.Input), "", "#00f")
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
		h.Player.Examine_ = string(h.Input)
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
	Input  []rune
}

func (h *AdminChangeColorHUD) Paint(setcell func(int, int, string, string, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	h.Player.Lock()
	setcell(0, 0, "SKIN COLOR", "", "#fff")
	h.Player.Unlock()

	setcell(0, 1, ">", "", "#00f")
	setcell(2, 1, string(h.Input), "", "#fff")
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
		h.Player.BaseColor = Color(h.Input)
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
		setcell(2, i+1, c, "", "#0ff")
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(2, 9, "previous", "", "#fff")
	}
	if len(h.Commands) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(2, 10, "next", "", "#fff")
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
