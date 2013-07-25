package main

import (
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

	camX := int(CameraX)
	camY := int(CameraY)

	termbox.SetCursor(w/2, h/2)

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

var CurrentZone *Zone
var CameraX, CameraY uint8 = 127, 127

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go pollEvents(events)
	repaint()
	CurrentZone = &Zone{X: 0, Y: 0, Element: Earth}
	CurrentZone.Tile(120, 120).Add(&Rock{})
	CurrentZone.Tile(121, 120).Add(&Rock{})
	CurrentZone.Tile(120, 121).Add(&Rock{})
	CurrentZone.Tile(121, 122).Add(&Rock{})

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
					if CameraX != 0 && !CurrentZone.Blocked(CameraX-1, CameraY) {
						CameraX--
						repaint()
					}
				case termbox.KeyArrowRight:
					if CameraX != 255 && !CurrentZone.Blocked(CameraX+1, CameraY) {
						CameraX++
						repaint()
					}
				case termbox.KeyArrowUp:
					if CameraY != 0 && !CurrentZone.Blocked(CameraX, CameraY-1) {
						CameraY--
						repaint()
					}
				case termbox.KeyArrowDown:
					if CameraY != 255 && !CurrentZone.Blocked(CameraX, CameraY+1) {
						CameraY++
						repaint()
					}

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
