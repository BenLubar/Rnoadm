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

func (l *Liquid) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Sprite: "ui_fill", // TODO
		Color:  "#44f",
	})
}

func (l *Liquid) ZIndex() int {
	return 0
}
