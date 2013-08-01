package main

type Liquid struct {
	Uninteractable
}

func (l *Liquid) Name() string {
	return "water"
}

func (l *Liquid) Examine() string {
	return "a pool of water."
}

func (l *Liquid) Blocking() bool {
	return true
}

func (l *Liquid) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "ui_fill", "#00f") // TODO
}
