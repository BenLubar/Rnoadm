package main

import (
	"github.com/nsf/termbox-go"
)

type Player struct {
	Hero
	ZoneX, ZoneY int64
	TileX, TileY uint8
}

type Hero struct{}

func (h *Hero) Name() string {
	return "hero"
}

func (h *Hero) Examine() string {
	return "a hero."
}

func (h *Hero) Blocking() bool {
	return false
}

func (h *Hero) Paint() (rune, termbox.Attribute) {
	return '☻', termbox.ColorWhite
}
