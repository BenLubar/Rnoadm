package main

import (
	"strings"
	"unicode"
)

type ExamineHUD struct {
	Player  *Player
	Name    string
	Examine string
}

func (h *ExamineHUD) Paint(setcell func(int, int, rune, Color)) {
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

	Drop      int
	Examine   int
	AdminTake int

	Inventory bool
	Slot      int
}

func (h *InteractMenuHUD) Paint(setcell func(int, int, rune, Color)) {
	if h.Options == nil {
		h.Options = append(h.Options, h.Object.InteractOptions()...)
		h.Drop = -1
		if h.Inventory {
			h.Drop = len(h.Options)
			h.Options = append(h.Options, "drop")
		}
		h.Examine = len(h.Options)
		h.Options = append(h.Options, "examine")
		h.AdminTake = -1
		if _, ok := h.Object.(*Player); h.Player.Admin && !ok {
			h.AdminTake = len(h.Options)
			if h.Inventory {
				h.Options = append(h.Options, "destroy [ADMIN]")
			} else {
				h.Options = append(h.Options, "take [ADMIN]")
			}
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
					Player:  h.Player,
					Name:    h.Object.Name(),
					Examine: h.Object.Examine(),
				}
				h.Player.Repaint()
			} else if i == h.Drop {
				h.Player.lock.Lock()
				if h.Slot < len(h.Player.Backpack) && h.Player.Backpack[h.Slot] == h.Object {
					zone := h.Player.zone
					tx, ty := h.Player.TileX, h.Player.TileY
					h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					h.Player.lock.Unlock()

					zone.Lock()
					zone.Tile(tx, ty).Add(h.Object)
					zone.Unlock()
					zone.Repaint()
				} else {
					h.Player.lock.Unlock()
				}
				h.Player.hud = nil
				h.Player.Repaint()
			} else if i == h.AdminTake {
				if h.Inventory {
					h.Player.lock.Lock()
					if h.Slot < len(h.Player.Backpack) && h.Player.Backpack[h.Slot] == h.Object {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.lock.Unlock()
					h.Player.hud = nil
					h.Player.Repaint()
					return true
				}

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

type InventoryHUD struct {
	Player *Player
	Offset int
}

func (h *InventoryHUD) Paint(setcell func(int, int, rune, Color)) {
	h.Player.lock.Lock()
	defer h.Player.lock.Unlock()

	if h.Offset > len(h.Player.Backpack) {
		h.Offset = 0
	}

	for i, r := range "INVENTORY" {
		setcell(i, 0, r, "#fff")
	}
	for i, o := range h.Player.Backpack[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, rune(i)+'1', "#fff")
		setcell(1, i+1, ' ', "#fff")
		r, color := o.Paint()
		setcell(2, i+1, r, color)
		setcell(3, i+1, ' ', "#fff")
		j := 4
		for _, r = range o.Name() {
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
	if len(h.Player.Backpack) > h.Offset+8 {
		setcell(0, 10, '0', "#fff")
		setcell(1, 10, ' ', "#fff")
		j := 2
		for _, r := range "next" {
			setcell(j, 10, r, "#fff")
			j++
		}
	}
}

func (h *InventoryHUD) Key(code int) bool {
	h.Player.lock.Lock()
	defer h.Player.lock.Unlock()

	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		i := code - '1' + h.Offset
		if i < len(h.Player.Backpack) {
			h.Player.hud = &InteractMenuHUD{
				Player:    h.Player,
				Object:    h.Player.Backpack[i],
				Inventory: true,
				Slot:      i,
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
		if h.Offset+8 < len(h.Player.Backpack) {
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
