package main

type ZoneHUD string

/*func (h ZoneHUD) Paint(setcell func(int, int, PaintCell)) {
	for i := 0; i < 20; i++ {
		setcell(i, 0, PaintCell{
			Sprite: "ui_fill_small",
			Color:  "rgba(0,0,0,0.7)",
			Y:      -16,
			ZIndex: 9999,
		})
		setcell(i, 0, PaintCell{
			Sprite: "ui_fill_small",
			Color:  "rgba(0,0,0,0.7)",
			X:      16,
			Y:      -16,
			ZIndex: 9999,
		})
	}
	setcell(0, 0, PaintCell{
		Text:   string(h),
		Color:  "#fff",
		Y:      -16,
		ZIndex: 10000,
	})
}*/

func (h ZoneHUD) Key(code int, special bool) bool {
	return false
}

func (h ZoneHUD) Click(x, y int) bool {
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
	Inventory    bool
}

/*func (h *ClickHUD) Paint(setcell func(int, int, PaintCell)) {
	if h.Options == nil {
		h.Player.Lock()
		zone := h.Player.zone
		px, py := h.Player.TileX, h.Player.TileY
		h.Player.Unlock()

		var objects []Object

		if h.Inventory {
			x, y := h.X-h.W-1, h.Y-3
			i := x + y*10
			h.Player.Lock()
			if i < 0 || len(h.Player.Backpack) <= i {
				h.Player.Unlock()
				h.Player.hud = nil
				h.Player.Repaint()
				return
			}
			objects = []Object{h.Player.Backpack[i]}
		} else {
			tx, ty := int(px)+h.X-h.W/2, int(py)+h.Y-h.H/2
			h.TileX, h.TileY = uint8(tx), uint8(ty)
			tile := zone.Tile(h.TileX, h.TileY)

			if tile == nil || int(h.TileX) != tx || int(h.TileY) != ty {
				h.Player.hud = nil
				h.Player.Repaint()
				return
			}
			objects = tile.Objects
			zone.Lock()
			h.Blocked = tile.Blocked()
		}

		h.Options = []clickHUDOption{}
		for _, o := range objects {
			if item, ok := o.(Item); ok {
				if h.Inventory {
					for i, s := range o.InteractOptions() {
						h.Options = append(h.Options, clickHUDOption{
							Object: o,
							Text:   s + " " + o.Name(),
							Exec: func(o Object, i int) func() {
								return func() {
									o.Interact(0, 0, h.Player, zone, i)
								}
							}(o, i),
						})
					}
				} else {
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
			} else if !h.Inventory {
				for i, s := range o.InteractOptions() {
					h.Options = append(h.Options, clickHUDOption{
						Object: o,
						Text:   s + " " + o.Name(),
						Exec: func(o Object, i int) func() {
							return func() {
								o.Interact(h.TileX, h.TileY, h.Player, zone, i)
							}
						}(o, i),
					})
				}
			}
			h.Options = append(h.Options, clickHUDOption{
				Object: o,
				Text:   "examine " + o.Name(),
				Exec: func(e string) func() {
					return func() {
						h.Player.SendMessage(e)
					}
				}(o.Examine()),
			})
			if _, ok := o.(*Player); h.Player.Admin && !ok && !h.Inventory {
				h.Options = append(h.Options, clickHUDOption{
					Object: o,
					Text:   "take [ADMIN] " + o.Name(),
					Exec: func(o Object) func() {
						return func() {
							if !h.Player.Admin {
								return
							}

							h.Player.Lock()
							zone := h.Player.zone
							h.Player.Unlock()

							zone.Lock()
							if !zone.Tile(h.TileX, h.TileY).Remove(o) {
								zone.Unlock()
								return
							}
							zone.Unlock()

							h.Player.Lock()
							h.Player.GiveItem(o)
							AdminLog.Printf("TAKE [%d:%q] (%d:%d %d:%d) %q %q", h.Player.ID, h.Player.Name(), h.Player.ZoneX, h.TileX, h.Player.ZoneY, h.TileY, o.Name(), o.Examine())
							h.Player.Unlock()
						}
					}(o),
				})
			}
		}
		if h.Inventory {
			h.Player.Unlock()
		} else {
			zone.Unlock()
		}

		if len(h.Options) == 0 {
			h.Click(h.X+1, h.Y)
			return
		}
	}
	x := h.X
	if h.Inventory {
		setcell(h.X-8, h.Y, PaintCell{
			Sprite: "ui_largecorner_tl",
			Color:  "#111",
			ZIndex: 10001,
		})
		x -= 8
	} else {
		setcell(h.X+8, h.Y, PaintCell{
			Sprite: "ui_largecorner_tr",
			Color:  "#111",
			ZIndex: 10001,
		})
	}
	for i := 1; i < 8; i++ {
		setcell(x+i, h.Y, PaintCell{
			Sprite: "ui_fill",
			Color:  "#111",
			ZIndex: 10001,
		})
	}

	if !h.Inventory && !h.Blocked {
		setcell(h.X+1, h.Y, PaintCell{
			Text:   "walk here",
			Color:  "#fff",
			ZIndex: 10002,
		})
	} else if h.Inventory {
		setcell(h.X-7, h.Y, PaintCell{
			Text:   "drop",
			Color:  "#fff",
			ZIndex: 10002,
		})
	}

	row := 1
	for _, option := range h.Options {
		for i := 0; i <= 8; i++ {
			if h.Inventory {
				setcell(h.X-i, h.Y+row, PaintCell{
					Sprite: "ui_fill",
					Color:  "#333",
					ZIndex: 10001,
				})
			} else {
				setcell(h.X+i, h.Y+row, PaintCell{
					Sprite: "ui_fill",
					Color:  "#333",
					ZIndex: 10001,
				})
			}
		}
		option.Object.Paint(h.X, h.Y+row, setcell)
		if h.Inventory {
			setcell(h.X-7, h.Y+row, PaintCell{
				Text:   option.Text,
				Color:  "#fff",
				ZIndex: 10002,
			})
		} else {
			setcell(h.X+1, h.Y+row, PaintCell{
				Text:   option.Text,
				Color:  "#fff",
				ZIndex: 10002,
			})
		}
		row++
	}

	x = h.X
	if h.Inventory {
		x -= 8
	}
	setcell(x, h.Y+row, PaintCell{
		Sprite: "ui_largecorner_bl",
		Color:  "#333",
		ZIndex: 10001,
	})
	for i := 1; i < 8; i++ {
		setcell(x+i, h.Y+row, PaintCell{
			Sprite: "ui_fill",
			Color:  "#333",
			ZIndex: 10001,
		})
	}
	setcell(x+8, h.Y+row, PaintCell{
		Sprite: "ui_largecorner_br",
		Color:  "#333",
		ZIndex: 10001,
	})
	setcell(x+1, h.Y+row, PaintCell{
		Text:   "cancel",
		Color:  "#fff",
		ZIndex: 10002,
	})
}*/

func (h *ClickHUD) Key(code int, special bool) bool {
	h.Player.hud = nil
	h.Player.Repaint()
	return false
}

func (h *ClickHUD) Click(x, y int) bool {
	if (h.Inventory && (x < h.X-8 || x > h.X)) || (!h.Inventory && (x < h.X || x > h.X+8)) || y < h.Y || y > h.Y+1+len(h.Options) {
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
		if h.Inventory {
			x, y := h.X-h.W-1, h.Y-3
			i := x + y*10
			h.Player.Lock()
			if i < 0 || len(h.Player.Backpack) <= i {
				h.Player.Unlock()
				h.Player.hud = nil
				h.Player.Repaint()
				return true
			}
			o := h.Player.Backpack[i]
			h.Player.Backpack = append(h.Player.Backpack[:i], h.Player.Backpack[i+1:]...)
			zone := h.Player.zone
			tile := zone.Tile(h.Player.TileX, h.Player.TileY)
			h.Player.Unlock()

			zone.Lock()
			tile.Add(o)
			zone.Unlock()
			zone.Repaint()

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
