package main

type FloraType uint16

type Flora struct {
	Type FloraType
	R    rune
}

func (f *Flora) Name() string {
	return "plant"
}

func (f *Flora) Examine() string {
	return "a plant."
}

func (f *Flora) Paint() (rune, Color) {
	return '"', "#8f0"
}

func (f *Flora) Blocking() bool {
	return false
}

func (f *Flora) InteractOptions() []string {
	return []string{"pick"}
}
