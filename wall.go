package main

type WallStone struct {
	Type RockType
	Uninteractable
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

func (w *WallStone) ZIndex() int {
	return 100
}

type WallMetal struct {
	Type MetalType
	Uninteractable
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

func (w *WallMetal) ZIndex() int {
	return 100
}

type WallWood struct {
	Type WoodType
	Uninteractable
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

func (w *WallWood) ZIndex() int {
	return 100
}

type FloorStone struct {
	Type RockType
	Uninteractable
}

func (f *FloorStone) Name() string {
	return rockTypeInfo[f.Type].Name + " floor"
}

func (f *FloorStone) Examine() string {
	return "a floor made of " + rockTypeInfo[f.Type].Name + "."
}

func (f *FloorStone) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "floor_stone", rockTypeInfo[f.Type].Color)
}

func (f *FloorStone) Blocking() bool {
	return false
}

func (f *FloorStone) ZIndex() int {
	return -50
}

type FloorMetal struct {
	Type MetalType
	Uninteractable
}

func (f *FloorMetal) Name() string {
	return metalTypeInfo[f.Type].Name + " floor"
}

func (f *FloorMetal) Examine() string {
	return "a floor made of " + metalTypeInfo[f.Type].Name + "."
}

func (f *FloorMetal) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "floor_metal", metalTypeInfo[f.Type].Color)
}

func (f *FloorMetal) Blocking() bool {
	return false
}

func (f *FloorMetal) ZIndex() int {
	return -50
}

type FloorWood struct {
	Type WoodType
	Uninteractable
}

func (f *FloorWood) Name() string {
	return woodTypeInfo[f.Type].Name + " floor"
}

func (f *FloorWood) Examine() string {
	return "a floor made of " + woodTypeInfo[f.Type].Name + "."
}

func (f *FloorWood) Paint(x, y int, setcell func(int, int, string, string, Color)) {
	setcell(x, y, "", "floor_wood", woodTypeInfo[f.Type].Color)
}

func (f *FloorWood) Blocking() bool {
	return false
}

func (f *FloorWood) ZIndex() int {
	return -50
}
