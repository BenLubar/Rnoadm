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

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go pollEvents(events)
	repaint()

	for {
		select {
		case event := <-events:
			switch event.Type {
			case termbox.EventError:
				panic(event.Err)

			case termbox.EventResize:
				repaint()

			case termbox.EventKey:
				// TODO: handle events
				return
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
