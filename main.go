package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

var shouldPaint = make(chan struct{}, 1)

func repaint() {
	select {
	case shouldPaint <- struct{}{}:
	default:
	}
}

var titlePerm []int
var HUD interface {
	Paint()
	Key(termbox.Event) bool
}

func paint() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)

	w, h := termbox.Size()

	camX := int(ThePlayer.TileX)
	camY := int(ThePlayer.TileY)

	termbox.SetCursor(-1, -1)

	CurrentZone.Lock()
	defer CurrentZone.Unlock()

	for x := 0; x < w; x++ {
		xCoord := x - w/2 + camX
		if xCoord < 0 || xCoord > 255 {
			continue
		}
		x8 := uint8(xCoord)
		for y := 0; y < h; y++ {
			yCoord := y - h/2 + camY
			if yCoord < 0 || yCoord > 255 {
				continue
			}
			y8 := uint8(yCoord)
			if tile := CurrentZone.Tile(x8, y8); tile != nil {
				r, fg := tile.Paint()
				bg := termbox.ColorBlack
				termbox.SetCell(x, y, r, fg, bg)
			}
		}
	}

	if titlePerm == nil {
		titlePerm = rand.Perm(4)
	}
	termbox.SetCell(w/2-3, h/4, 'R', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	termbox.SetCell(w/2+2, h/4, 'm', termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	for i, j := range titlePerm {
		termbox.SetCell(w/2-2+i, h/4, rune("ando"[j]), termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	}

	if HUD != nil {
		HUD.Paint()
	}

	termbox.Flush()
}

var CurrentZone *Zone
var ThePlayer *Player

var Seed int64

func main() {
	flag.Int64Var(&Seed, "seed", time.Now().UnixNano(), "the world seed (default: the number of nanoseconds since midnight UTC on 1970-01-01)")

	flag.Parse()

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go pollEvents(events)
	repaint()
	if z, err := LoadZone(0, 0); err != nil {
		CurrentZone = &Zone{X: 0, Y: 0}
		CurrentZone.Generate()
	} else {
		CurrentZone = z
	}
	if p, err := LoadPlayer(0); err != nil {
		ThePlayer = &Player{TileX: 127, TileY: 127}
	} else {
		ThePlayer = p
	}
	CurrentZone.Lock()
	CurrentZone.Tile(ThePlayer.TileX, ThePlayer.TileY).Add(ThePlayer)
	CurrentZone.Unlock()
	defer func() {
		err := CurrentZone.Save()
		if err != nil {
			panic(err)
		}
		err = ThePlayer.Save()
		if err != nil {
			panic(err)
		}
	}()

	ticker := time.Tick(time.Second)

	for {
		select {
		case event := <-events:
			switch event.Type {
			case termbox.EventError:
				panic(event.Err)

			case termbox.EventResize:
				repaint()

			case termbox.EventKey:
				if HUD != nil && HUD.Key(event) {
					break
				}

				if event.Ch != 0 {
					switch event.Ch {
					case 'w':
						ThePlayer.Move(0, -1)
					case 'a':
						ThePlayer.Move(-1, 0)
					case 's':
						ThePlayer.Move(0, 1)
					case 'd':
						ThePlayer.Move(1, 0)

					case 'e':
						HUD = &InteractHUD{Player: ThePlayer}
						repaint()

					default:
						// TODO: handle more keys
						return
					}
					break
				}

				switch event.Key {
				case termbox.KeyArrowLeft:
					ThePlayer.Move(-1, 0)
				case termbox.KeyArrowRight:
					ThePlayer.Move(1, 0)
				case termbox.KeyArrowUp:
					ThePlayer.Move(0, -1)
				case termbox.KeyArrowDown:
					ThePlayer.Move(0, 1)

				default:
					// TODO: handle more keys
					return
				}
			}

		case <-ticker:
			CurrentZone.Think()

		case <-shouldPaint:
			paint()
		}
	}
}

func pollEvents(ch chan<- termbox.Event) {
	for {
		ch <- termbox.PollEvent()
	}
}

type InteractHUD struct {
	Player       *Player
	TileX, TileY uint8
	Objects      []Object
	Offset       int
}

func (h *InteractHUD) Paint() {
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
		var objects []Object
		for x := minX; x >= minX && x <= maxX; x++ {
			for y := minY; y >= minY && y <= maxY; y++ {
				objects = append(objects, CurrentZone.Tile(x, y).Objects...)
			}
		}
		h.Objects = objects
		h.Offset = 0
	}
	const keys = "12345678"
	for i, o := range h.Objects[h.Offset:] {
		if i >= len(keys) {
			break
		}
		termbox.SetCell(0, i, rune(keys[i]), termbox.AttrBold|termbox.AttrReverse, 0)
		termbox.SetCell(1, i, ' ', termbox.AttrReverse, 0)
		j := 1
		for _, r := range o.Name() {
			j++
			termbox.SetCell(j, i, r, termbox.AttrReverse, 0)
		}
	}
	if h.Offset > 0 {
		termbox.SetCell(0, 8, '9', termbox.AttrBold|termbox.AttrReverse, 0)
		termbox.SetCell(1, 8, ' ', termbox.AttrReverse, 0)
		j := 1
		for _, r := range "previous" {
			j++
			termbox.SetCell(j, 8, r, termbox.AttrReverse, 0)
		}
	}
	if len(h.Objects) > h.Offset+len(keys) {
		termbox.SetCell(0, 9, '0', termbox.AttrBold|termbox.AttrReverse, 0)
		termbox.SetCell(1, 9, ' ', termbox.AttrReverse, 0)
		j := 1
		for _, r := range "next" {
			j++
			termbox.SetCell(j, 9, r, termbox.AttrReverse, 0)
		}
	}
}

func (h *InteractHUD) Key(e termbox.Event) bool {
	switch e.Ch {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		// TODO
		return true
	case '9':
		if h.Offset > 0 {
			h.Offset--
			repaint()
		}
		return true
	case '0':
		if h.Offset+8 < len(h.Objects) {
			h.Offset++
			repaint()
		}
		return true

	case 0:
		switch e.Key {
		case termbox.KeyEsc:
			HUD = nil
			repaint()
			return true
		}
	}
	return false
}
