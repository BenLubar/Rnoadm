package main

import (
	"unicode"
)

type ExamineHUD struct {
	Player *Player
	Object Object
}

func (h *ExamineHUD) Paint(setcell func(int, int, rune, Color)) {
	i := 0
	for _, r := range h.Object.Name() {
		setcell(i, 0, unicode.ToUpper(r), "#fff")
		i++
	}
	i = 0
	for _, r := range h.Object.Examine() {
		setcell(i, 1, r, "#fff")
		i++
	}
}

func (h *ExamineHUD) Key(code int) bool {
	switch code {
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

type InteractHUD struct {
	Player       *Player
	TileX, TileY uint8
	Objects      []Object
	Offset       int
}

func (h *InteractHUD) Paint(setcell func(int, int, rune, Color)) {
	if h.Player.TileX != h.TileX || h.Player.TileY != h.TileY || h.Objects == nil {
		h.TileX, h.TileY = h.Player.TileX, h.Player.TileY
		minX := h.TileX - 1
		if minX == 255 {
			minX = 0
		}
		maxX := h.TileX + 1
		if maxX == 0 {
			maxX = 255
		}
		minY := h.TileY - 1
		if minY == 255 {
			minY = 0
		}
		maxY := h.TileY + 1
		if maxY == 0 {
			maxY = 255
		}
		z := GrabZone(h.Player.ZoneX, h.Player.ZoneY)
		z.Lock()
		var objects []Object
		for x := minX; x >= minX && x <= maxX; x++ {
			for y := minY; y >= minY && y <= maxY; y++ {
				objects = append(objects, z.Tile(x, y).Objects...)
			}
		}
		z.Unlock()
		ReleaseZone(z)
		h.Objects = objects
		h.Offset = 0
	}
	for i, r := range "EXAMINE" {
		setcell(i, 0, r, "#fff")
	}
	const keys = "12345678"
	for i, o := range h.Objects[h.Offset:] {
		if i >= len(keys) {
			break
		}
		setcell(0, i+1, rune(keys[i]), "#fff")
		setcell(1, i+1, ' ', "#fff")
		j := 2
		for _, r := range o.Name() {
			setcell(j, i+1, r, "#fff")
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
	if len(h.Objects) > h.Offset+len(keys) {
		setcell(0, 10, '0', "#fff")
		setcell(1, 10, ' ', "#fff")
		j := 2
		for _, r := range "next" {
			setcell(j, 10, r, "#fff")
			j++
		}
	}
}

func (h *InteractHUD) Key(code int) bool {
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.Objects) {
			h.Player.hud = &ExamineHUD{
				Player: h.Player,
				Object: h.Objects[i],
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
		if h.Offset+8 < len(h.Objects) {
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
