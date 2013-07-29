package main

import (
	"strings"
	"unicode"
)

type ExamineHUD struct {
	Player  *Player
	Object  Object
	Name    string
	Examine string
}

func (h *ExamineHUD) Paint(setcell func(int, int, rune, Color)) {
	if h.Object != nil {
		h.Name = h.Object.Name()

		h.Player.lock.Lock()
		z := h.Player.zone
		tx, ty := h.Player.TileX, h.Player.TileY
		h.Player.lock.Unlock()

		minX := tx - 1
		if minX > tx {
			minX = 0
		}
		maxX := tx + 1
		if maxX < tx {
			maxX = 255
		}
		minY := ty - 1
		if minY > ty {
			minY = 0
		}
		maxY := ty + 1
		if maxY < ty {
			maxY = 255
		}

		z.Lock()
		for x := minX; x >= minX && x <= maxX && h.Object != nil; x++ {
			for y := minY; y >= minY && y <= maxY && h.Object != nil; y++ {
				t := z.Tile(x, y)
				if t != nil {
					for _, o := range t.Objects {
						if o == h.Object {
							h.Object = nil
							h.Examine = o.Examine()
							break
						}
					}
				}
			}
		}
		z.Unlock()

		if h.Object != nil {
			h.Examine = "that is too far away!"
			h.Object = nil
		}
	}

	i := 0
	for _, r := range h.Name {
		setcell(i, 0, unicode.ToUpper(r), "#fff")
		i++
	}
	i = 0
	for _, r := range h.Examine {
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
	h.Player.lock.Lock()
	tx, ty := h.Player.TileX, h.Player.TileY
	h.Player.lock.Unlock()

	if tx != h.TileX || ty != h.TileY || h.Objects == nil {
		h.TileX, h.TileY = tx, ty
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
		h.Player.lock.Lock()
		z := h.Player.zone
		h.Player.lock.Unlock()
		z.Lock()
		var objects []Object
		for x := minX; x >= minX && x <= maxX; x++ {
			for y := minY; y >= minY && y <= maxY; y++ {
				objects = append(objects, z.Tile(x, y).Objects...)
			}
		}
		z.Unlock()
		h.Objects = objects
		h.Offset = 0
	}
	for i, r := range "INTERACT" {
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
			h.Player.hud = &InteractMenuHUD{
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

type InteractMenuHUD struct {
	Player  *Player
	Object  Object
	Options []string
	Offset  int

	Examine   int
	AdminTake int
}

func (h *InteractMenuHUD) Paint(setcell func(int, int, rune, Color)) {
	if h.Options == nil {
		h.Options = append(h.Options, h.Object.InteractOptions()...)
		h.Examine = len(h.Options)
		h.Options = append(h.Options, "examine")
		h.AdminTake = -1
		if _, ok := h.Object.(*Player); h.Player.Admin && !ok {
			h.AdminTake = len(h.Options)
			h.Options = append(h.Options, "take [ADMIN]")
		}
	}

	i := 0
	for _, r := range strings.ToUpper(h.Object.Name()) {
		setcell(i, 0, r, "#fff")
		i++
	}
	for i, o := range h.Options[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, rune(i)+'1', "#fff")
		setcell(1, i+1, ' ', "#fff")
		j := 2
		for _, r := range o {
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
	if len(h.Options) > h.Offset+8 {
		setcell(0, 10, '0', "#fff")
		setcell(1, 10, ' ', "#fff")
		j := 2
		for _, r := range "next" {
			setcell(j, 10, r, "#fff")
			j++
		}
	}

}

func (h *InteractMenuHUD) Key(code int) bool {
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.Options) {
			if i == h.Examine {
				h.Player.hud = &ExamineHUD{
					Player: h.Player,
					Object: h.Object,
				}
				h.Player.Repaint()
			} else if i == h.AdminTake {
				h.Player.lock.Lock()
				z := h.Player.zone
				tx, ty := h.Player.TileX, h.Player.TileY
				h.Player.lock.Unlock()

				minX := tx - 1
				if minX > tx {
					minX = 0
				}
				maxX := tx + 1
				if maxX < tx {
					maxX = 255
				}
				minY := ty - 1
				if minY > ty {
					minY = 0
				}
				maxY := ty + 1
				if maxY < ty {
					maxY = 255
				}

				z.Lock()
				found := false
				for x := minX; !found && x >= minX && x <= maxX; x++ {
					for y := minY; !found && y >= minY && y <= maxY; y++ {
						t := z.Tile(x, y)
						if t != nil {
							for _, o := range t.Objects {
								if o == h.Object {
									t.Remove(o)
									z.Unlock()
									h.Player.lock.Lock()
									h.Player.GiveItem(o)
									AdminLog.Printf("TAKE [%d:%q] (%d:%d %d:%d) %q %q", h.Player.ID, h.Player.Name(), h.Player.ZoneX, x, h.Player.ZoneY, y, o.Name(), o.Examine())
									h.Player.lock.Unlock()
									found = true
									break
								}
							}
						}
					}
				}
				if !found {
					z.Unlock()
				}

				h.Player.hud = nil
				h.Player.Repaint()
			} else {
				// TODO
			}
		}
		return true
	case '9':
		if h.Offset > 0 {
			h.Offset--
			h.Player.Repaint()
		}
		return true
	case '0':
		if h.Offset+8 < len(h.Options) {
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
