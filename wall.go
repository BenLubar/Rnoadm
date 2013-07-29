package main

type Wall struct {
}

func (w *Wall) Name() string {
	return "wall"
}

func (w *Wall) Examine() string {
	return "a wall."
}

func (w *Wall) Paint() (rune, Color) {
	return '\u2588', "#fff"
}

func (w *Wall) Blocking() bool {
	return true
}

func (w *Wall) InteractOptions() []string {
	return nil
}
