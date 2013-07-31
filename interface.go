package main

import (
	"strings"
)

type ExamineHUD struct {
	Player  *Player
	Name    string
	Examine string
}

func (h *ExamineHUD) Paint(setcell func(int, int, string, string, Color)) {
	setcell(0, 0, strings.ToUpper(h.Name), "", "#fff")
	setcell(0, 1, h.Examine, "", "#fff")
}

func (h *ExamineHUD) Key(code int, special bool) bool {
	if !special {
		return false
	}
	switch code {
	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}

func (h *ExamineHUD) Click(x, y int) bool {
	return false
}

type InteractHUD struct {
	Player       *Player
	TileX, TileY uint8
	Objects      []Object
	Offset       int
}

func (h *InteractHUD) Paint(setcell func(int, int, string, string, Color)) {
	h.Player.Lock()
	tx, ty := h.Player.TileX, h.Player.TileY
	h.Player.Unlock()

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
		h.Player.Lock()
		z := h.Player.zone
		h.Player.Unlock()
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
	setcell(0, 0, "INTERACT", "", "#fff")
	for i, o := range h.Objects[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#fff")
		setcell(2, i+1, o.Name(), "", "#fff")
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(2, 9, "previous", "", "#fff")
	}
	if len(h.Objects) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(2, 10, "next", "", "#fff")
	}
}

func (h *InteractHUD) Key(code int, special bool) bool {
	if !special {
		return false
	}
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

func (h *InteractHUD) Click(x, y int) bool {
	return false
}

type InteractMenuHUD struct {
	Player  *Player
	Object  Object
	Options []string
	Offset  int

	Wear      int
	Take      int
	Drop      int
	Examine   int
	AdminTake int

	Inventory bool
	Slot      int
}

func (h *InteractMenuHUD) Paint(setcell func(int, int, string, string, Color)) {
	if h.Options == nil {
		h.Options = append(h.Options, h.Object.InteractOptions()...)
		h.Wear = -1
		if _, ok := h.Object.(*Hat); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "wear")
		}
		if _, ok := h.Object.(*Shirt); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "wear")
		}
		if _, ok := h.Object.(*Pants); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "wear")
		}
		if _, ok := h.Object.(*Shoes); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "wear")
		}
		if _, ok := h.Object.(*Pickaxe); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "put on toolbelt")
		}
		if _, ok := h.Object.(*Hatchet); h.Inventory && ok {
			h.Wear = len(h.Options)
			h.Options = append(h.Options, "put on toolbelt")
		}
		h.Take = -1
		if _, ok := h.Object.(Item); !h.Inventory && ok {
			h.Take = len(h.Options)
			h.Options = append(h.Options, "take")
		}
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

	setcell(0, 0, strings.ToUpper(h.Object.Name()), "", "#fff")
	for i, o := range h.Options[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#fff")
		setcell(2, i+1, o, "", "#fff")
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(2, 9, "previous", "", "#fff")
	}
	if len(h.Options) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(2, 10, "next", "", "#fff")
	}
}

func (h *InteractMenuHUD) Key(code int, special bool) bool {
	if !special {
		return false
	}
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
				h.Player.Lock()
				if h.Slot < len(h.Player.Backpack) && h.Player.Backpack[h.Slot] == h.Object {
					zone := h.Player.zone
					tx, ty := h.Player.TileX, h.Player.TileY
					h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					h.Player.Unlock()

					zone.Lock()
					zone.Tile(tx, ty).Add(h.Object)
					zone.Unlock()
					zone.Repaint()
				} else {
					h.Player.Unlock()
				}
				h.Player.hud = nil
				h.Player.Repaint()
			} else if i == h.Wear {
				h.Player.Lock()
				switch o := h.Object.(type) {
				case *Hat:
					if h.Player.Head != nil {
						h.Player.Backpack[h.Slot] = h.Player.Head
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Head = o
				case *Shirt:
					if h.Player.Top != nil {
						h.Player.Backpack[h.Slot] = h.Player.Top
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Top = o
				case *Pants:
					if h.Player.Legs != nil {
						h.Player.Backpack[h.Slot] = h.Player.Legs
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Legs = o
				case *Shoes:
					if h.Player.Feet != nil {
						h.Player.Backpack[h.Slot] = h.Player.Feet
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Feet = o
				case *Pickaxe:
					if h.Player.Toolbelt.Pickaxe != nil {
						h.Player.Backpack[h.Slot] = h.Player.Toolbelt.Pickaxe
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Toolbelt.Pickaxe = o
				case *Hatchet:
					if h.Player.Toolbelt.Hatchet != nil {
						h.Player.Backpack[h.Slot] = h.Player.Toolbelt.Hatchet
					} else {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Toolbelt.Hatchet = o
				}

				h.Player.zone.Repaint()
				h.Player.Unlock()
				h.Player.hud = nil
				h.Player.Repaint()
			} else if i == h.Take || i == h.AdminTake {
				if h.Inventory {
					h.Player.Lock()
					if h.Slot < len(h.Player.Backpack) && h.Player.Backpack[h.Slot] == h.Object {
						h.Player.Backpack = append(h.Player.Backpack[:h.Slot], h.Player.Backpack[h.Slot+1:]...)
					}
					h.Player.Unlock()
					h.Player.hud = nil
					h.Player.Repaint()
					return true
				}

				h.Player.Lock()
				z := h.Player.zone
				tx, ty := h.Player.TileX, h.Player.TileY
				h.Player.Unlock()

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
									z.Repaint()
									h.Player.Lock()
									h.Player.GiveItem(o)
									if i == h.AdminTake {
										AdminLog.Printf("TAKE [%d:%q] (%d:%d %d:%d) %q %q", h.Player.ID, h.Player.Name(), h.Player.ZoneX, x, h.Player.ZoneY, y, o.Name(), o.Examine())
									}
									h.Player.Unlock()
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

func (h *InteractMenuHUD) Click(x, y int) bool {
	return false
}

type InventoryHUD struct {
	Player *Player
	Offset int
}

func (h *InventoryHUD) Paint(setcell func(int, int, string, string, Color)) {
	h.Player.Lock()
	defer h.Player.Unlock()

	if h.Offset > len(h.Player.Backpack) {
		h.Offset = 0
	}

	setcell(0, 0, "INVENTORY", "", "#fff")
	for i, o := range h.Player.Backpack[h.Offset:] {
		if i == 8 {
			break
		}
		setcell(0, i+1, string(rune(i)+'1'), "", "#fff")
		o.Paint(2, i+1, setcell)
		setcell(4, i+1, o.Name(), "", "#fff")
	}
	if h.Offset > 0 {
		setcell(0, 9, "9", "", "#fff")
		setcell(2, 9, "previous", "", "#fff")
	}
	if len(h.Player.Backpack) > h.Offset+8 {
		setcell(0, 10, "0", "", "#fff")
		setcell(2, 10, "next", "", "#fff")
	}
}

func (h *InventoryHUD) Key(code int, special bool) bool {
	if !special {
		return false
	}

	h.Player.Lock()
	defer h.Player.Unlock()

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

func (h *InventoryHUD) Click(x, y int) bool {
	return false
}

type clickHUDOption struct {
	Object Object
	Text   string
	Exec   func()
}

type ClickHUD struct {
	X, Y         int
	W, H         int
	Player       *Player
	Options      []clickHUDOption
	Blocked      bool
	TileX, TileY uint8
}

func (h *ClickHUD) Paint(setcell func(int, int, string, string, Color)) {
	if h.Options == nil {
		h.Player.Lock()
		zone := h.Player.zone
		px, py := h.Player.TileX, h.Player.TileY
		tx, ty := int(px)+h.X-h.W/2, int(py)+h.Y-h.H/2
		h.TileX, h.TileY = uint8(tx), uint8(ty)
		h.Player.Unlock()

		tile := zone.Tile(h.TileX, h.TileY)

		if tile == nil || int(h.TileX) != tx || int(h.TileY) != ty {
			h.Player.hud = nil
			h.Player.Repaint()
			return
		}

		zone.Lock()
		h.Blocked = tile.Blocked()
		h.Options = []clickHUDOption{}
		for _, o := range tile.Objects {
			for _, i := range o.InteractOptions() {
				h.Options = append(h.Options, clickHUDOption{
					Object: o,
					Text:   i + " " + o.Name(),
					Exec:   func() {}, // TODO
				})
			}
			if item, ok := o.(Item); ok {
				h.Options = append(h.Options, clickHUDOption{
					Object: o,
					Text:   "take " + o.Name(),
					Exec: func() {
						var schedule Schedule = &TakeSchedule{
							Item: item.(Object),
						}
						if px != h.TileX || py != h.TileY {
							moveSchedule := MoveSchedule(FindPath(zone, px, py, h.TileX, h.TileY, true))
							schedule = &ScheduleSchedule{&moveSchedule, schedule}
						}
						h.Player.Lock()
						h.Player.schedule = schedule
						h.Player.Unlock()
					},
				})
			}
			h.Options = append(h.Options, clickHUDOption{
				Object: o,
				Text:   "examine " + o.Name(),
				Exec: func(o Object) func() {
					return func() {
						h.Player.hud = &ExamineHUD{
							Player:  h.Player,
							Name:    o.Name(),
							Examine: o.Examine(),
						}
					}
				}(o),
			})
		}
		zone.Unlock()

		if len(h.Options) == 0 {
			h.Click(h.X+1, h.Y)
			return
		}
	}
	for i := 1; i < 8; i++ {
		setcell(h.X+i, h.Y, "", "ui_fill", "#111")
	}
	setcell(h.X+8, h.Y, "", "ui_largecorner_tr", "#111")

	if !h.Blocked {
		setcell(h.X+1, h.Y, "walk here", "", "#fff")
	}

	row := 1
	for _, option := range h.Options {
		for i := 0; i <= 8; i++ {
			setcell(h.X+i, h.Y+row, "", "ui_fill", "#333")
		}
		option.Object.Paint(h.X, h.Y+row, setcell)
		setcell(h.X+1, h.Y+row, option.Text, "", "#fff")
		row++
	}

	setcell(h.X, h.Y+row, "", "ui_largecorner_bl", "#333")
	for i := 1; i < 8; i++ {
		setcell(h.X+i, h.Y+row, "", "ui_fill", "#333")
	}
	setcell(h.X+8, h.Y+row, "", "ui_largecorner_br", "#333")
	setcell(h.X+1, h.Y+row, "cancel", "", "#fff")
}

func (h *ClickHUD) Key(code int, special bool) bool {
	h.Player.hud = nil
	h.Player.Repaint()
	return false
}

func (h *ClickHUD) Click(x, y int) bool {
	if x < h.X || x > h.X+8 || y < h.Y || y > h.Y+1+len(h.Options) {
		h.Player.hud = nil
		h.Player.Repaint()
		return false
	}
	y -= h.Y
	switch y {
	case 0:
		if x == h.X {
			h.Player.hud = nil
			h.Player.Repaint()
			return true
		}
		h.Player.Lock()
		px, py := h.Player.TileX, h.Player.TileY
		zone := h.Player.zone
		h.Player.Unlock()
		schedule := MoveSchedule(FindPath(zone, px, py, h.TileX, h.TileY, true))
		h.Player.Lock()
		h.Player.schedule = &schedule
		h.Player.Unlock()

		h.Player.hud = nil
		h.Player.Repaint()
		return true
	case 1 + len(h.Options):
		h.Player.hud = nil
		h.Player.Repaint()
	default:
		h.Player.hud = nil
		h.Options[y-1].Exec()
		h.Player.Repaint()
	}
	return true
}
