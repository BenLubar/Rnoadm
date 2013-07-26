package main

import (
	"github.com/nsf/termbox-go"
)

type Wall struct {
}

func (w *Wall) Name() string {
	return "wall"
}

func (w *Wall) Examine() string {
	return "a wall"
}

func (w *Wall) Paint() (rune, termbox.Attribute) {
	return '\u2588', termbox.ColorWhite
}

func (w *Wall) Blocking() bool {
	return true
}
