package main

import (
	"github.com/nsf/termbox-go"
)

type Rock struct {
}

func (r *Rock) Name() string {
	return "rock"
}

func (r *Rock) Examine() string {
	return "a rock"
}

func (r *Rock) Paint() (rune, termbox.Attribute) {
	return '\u25B2', termbox.ColorWhite
}

func (r *Rock) Blocking() bool {
	return true
}
