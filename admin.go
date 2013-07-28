package main

import (
	"log"
	"sort"
)

var AdminLog *log.Logger

var adminCommands = map[string]func(*Player){
	"TP": func(p *Player) {
		p.hud = &AdminTeleportHUD{Player: p}
		p.Repaint()
	},
}

type AdminHUD struct {
	Player *Player
	Input  []rune
}

func (h *AdminHUD) Paint(setcell func(int, int, rune, Color)) {
	if !h.Player.Admin {
		h.Player.hud = nil
		return
	}

	setcell(0, 0, '>', "#00f")
	setcell(1, 0, ' ', "#00f")
	for i, r := range h.Input {
		setcell(i+2, 0, r, "#00f")
	}
}

func (h *AdminHUD) Key(code int) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
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
			f(h.Player)
		}
		return true
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	h.Input = append(h.Input, rune(code))
	h.Player.Repaint()
	return true
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

func (h *AdminTeleportHUD) Paint(setcell func(int, int, rune, Color)) {
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
	for i, r := range "TELEPORT TO PLAYER" {
		setcell(i, 0, r, "#fff")
	}
	for i, p := range h.List[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, '1'+rune(i), "#fff")
		setcell(1, i+1, ' ', "#fff")
		id := p.ID
		for k := 0; k < 16; k++ {
			setcell(17-k, i+1, rune("0123456789ABCDEF"[id&15]), "#44f")
			id >>= 4
		}
		setcell(18, i+1, ' ', "#00f")
		j := 19
		p.lock.Lock()
		name := p.Name()
		p.lock.Unlock()
		for _, r := range name {
			setcell(j, i+1, r, "#00f")
			j++
		}
	}
	if h.Offset > 0 {
		setcell(0, 9, '9', "#fff")
		setcell(1, 9, ' ', "#fff")
		j := 1
		for _, r := range "previous" {
			j++
			setcell(j, 9, r, "#fff")
		}
	}
	if len(h.List) > h.Offset+8 {
		setcell(0, 10, '0', "#fff")
		setcell(1, 10, ' ', "#fff")
		j := 2
		for _, r := range "next" {
			setcell(j, 10, r, "#fff")
			j++
		}
	}
}

func (h *AdminTeleportHUD) Key(code int) bool {
	if !h.Player.Admin {
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.List) {
			p := h.List[i]
			p.lock.Lock()
			zx, zy := p.ZoneX, p.ZoneY
			tx, ty := p.TileX, p.TileY
			name := p.Name()
			p.lock.Unlock()

			h.Player.lock.Lock()
			az := h.Player.zone
			azx, azy := h.Player.ZoneX, h.Player.ZoneY
			atx, aty := h.Player.TileX, h.Player.TileY
			aname := h.Player.Name()
			h.Player.lock.Unlock()

			AdminLog.Printf("TELEPORT [%d:%q] (%d:%d, %d:%d) => [%d:%q] (%d:%d, %d:%d)", h.Player.ID, aname, azx, atx, azy, aty, p.ID, name, zx, tx, zy, ty)

			az.Lock()
			az.Tile(atx, aty).Remove(h.Player)
			az.Repaint()
			az.Unlock()

			ReleaseZone(az)
			z := GrabZone(zx, zy)

			h.Player.lock.Lock()
			h.Player.zone = z
			h.Player.ZoneX, h.Player.ZoneY = zx, zy
			h.Player.TileX, h.Player.TileY = tx, ty
			h.Player.lock.Unlock()
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
