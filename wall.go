package main

type WallStone struct {
	Type RockType
}

func (w *WallStone) Name() string {
	return rockTypeInfo[w.Type].Name + " wall"
}

func (w *WallStone) Examine() string {
	return "a wall made of " + rockTypeInfo[w.Type].Name + "."
}

func (w *WallStone) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "wall_stone", rockTypeInfo[w.Type].Color)
}

func (w *WallStone) Blocking() bool {
	return true
}

func (w *WallStone) InteractOptions() []string {
	return nil
}

func (w *WallStone) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
}

type WallMetal struct {
	Type MetalType
}

func (w *WallMetal) Name() string {
	return metalTypeInfo[w.Type].Name + " wall"
}

func (w *WallMetal) Examine() string {
	return "a wall made of " + metalTypeInfo[w.Type].Name + "."
}

func (w *WallMetal) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "wall_metal", metalTypeInfo[w.Type].Color)
}

func (w *WallMetal) Blocking() bool {
	return true
}

func (w *WallMetal) InteractOptions() []string {
	return nil
}

func (w *WallMetal) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
}

type WallWood struct {
	Type WoodType
}

func (w *WallWood) Name() string {
	return woodTypeInfo[w.Type].Name + " wall"
}

func (w *WallWood) Examine() string {
	return "a wall made of " + woodTypeInfo[w.Type].Name + "."
}

func (w *WallWood) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "wall_wood", woodTypeInfo[w.Type].Color)
}

func (w *WallWood) Blocking() bool {
	return true
}

func (w *WallWood) InteractOptions() []string {
	return nil
}

func (w *WallWood) Interact(x uint8, y uint8, player *Player, zone *Zone, opt int) {
}
