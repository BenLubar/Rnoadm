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

	termbox.Flush()
}

func move(dx, dy int) {
	for dx+int(ThePlayer.TileX) > 255 {
		dx--
	}
	for dx+int(ThePlayer.TileX) < 0 {
		dx++
	}
	for dy+int(ThePlayer.TileY) > 255 {
		dy--
	}
	for dy+int(ThePlayer.TileY) < 0 {
		dy++
	}
	CurrentZone.Lock()
	defer CurrentZone.Unlock()
	if CurrentZone.Blocked(ThePlayer.TileX+uint8(dx), ThePlayer.TileY+uint8(dy)) {
		return
	}
	CurrentZone.Tile(ThePlayer.TileX, ThePlayer.TileY).Remove(ThePlayer)
	ThePlayer.TileX += uint8(dx)
	ThePlayer.TileY += uint8(dy)
	CurrentZone.Tile(ThePlayer.TileX, ThePlayer.TileY).Add(ThePlayer)
	repaint()
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
	ThePlayer = &Player{TileX: 127, TileY: 127}
	CurrentZone.Lock()
	CurrentZone.Tile(ThePlayer.TileX, ThePlayer.TileY).Add(ThePlayer)
	CurrentZone.Unlock()
	defer func() {
		err := CurrentZone.Save()
		if err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case event := <-events:
			switch event.Type {
			case termbox.EventError:
				panic(event.Err)

			case termbox.EventResize:
				repaint()

			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyArrowLeft:
					move(-1, 0)
				case termbox.KeyArrowRight:
					move(1, 0)
				case termbox.KeyArrowUp:
					move(0, -1)
				case termbox.KeyArrowDown:
					move(0, 1)

				default:
					// TODO: handle more keys
					return
				}
			}

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
