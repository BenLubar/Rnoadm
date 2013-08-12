package main

type WallStone struct {
	networkID
	Type    RockType
	Quality uint64
}

func (w *WallStone) Name() string {
	return rockTypeInfo[w.Type].Name + " wall"
}

func (w *WallStone) Examine() string {
	return "a wall made of " + rockTypeInfo[w.Type].Name + "."
}

func (w *WallStone) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   w.Name(),
		Sprite: "wall_stone",
		Colors: []Color{rockTypeInfo[w.Type].Color},
	}
}

func (w *WallStone) Blocking() bool {
	return true
}

func (w *WallStone) ZIndex() int {
	return 100
}

type WallMetal struct {
	networkID
	Type    MetalType
	Quality uint64
}

func (w *WallMetal) Name() string {
	return metalTypeInfo[w.Type].Name + " wall"
}

func (w *WallMetal) Examine() string {
	return "a wall made of " + metalTypeInfo[w.Type].Name + "."
}

func (w *WallMetal) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   w.Name(),
		Sprite: "wall_metal",
		Colors: []Color{metalTypeInfo[w.Type].Color},
	}
}

func (w *WallMetal) Blocking() bool {
	return true
}

func (w *WallMetal) ZIndex() int {
	return 100
}

type WallWood struct {
	networkID
	Type    WoodType
	Quality uint64
}

func (w *WallWood) Name() string {
	return woodTypeInfo[w.Type].Name + " wall"
}

func (w *WallWood) Examine() string {
	return "a wall made of " + woodTypeInfo[w.Type].Name + "."
}

func (w *WallWood) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   w.Name(),
		Sprite: "wall_wood",
		Colors: []Color{woodTypeInfo[w.Type].Color},
	}
}

func (w *WallWood) Blocking() bool {
	return true
}

func (w *WallWood) ZIndex() int {
	return 100
}

type FloorStone struct {
	networkID
	Type    RockType
	Quality uint64
}

func (f *FloorStone) Name() string {
	return rockTypeInfo[f.Type].Name + " floor"
}

func (f *FloorStone) Examine() string {
	return "a floor made of " + rockTypeInfo[f.Type].Name + "."
}

func (f *FloorStone) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   f.Name(),
		Sprite: "floor_stone",
		Colors: []Color{rockTypeInfo[f.Type].Color},
	}
}

func (f *FloorStone) Blocking() bool {
	return false
}

func (f *FloorStone) ZIndex() int {
	return -50
}

type FloorMetal struct {
	networkID
	Type    MetalType
	Quality uint64
}

func (f *FloorMetal) Name() string {
	return metalTypeInfo[f.Type].Name + " floor"
}

func (f *FloorMetal) Examine() string {
	return "a floor made of " + metalTypeInfo[f.Type].Name + "."
}

func (f *FloorMetal) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   f.Name(),
		Sprite: "floor_metal",
		Colors: []Color{metalTypeInfo[f.Type].Color},
	}
}

func (f *FloorMetal) Blocking() bool {
	return false
}

func (f *FloorMetal) ZIndex() int {
	return -50
}

type FloorWood struct {
	networkID
	Type    WoodType
	Quality uint64
}

func (f *FloorWood) Name() string {
	return woodTypeInfo[f.Type].Name + " floor"
}

func (f *FloorWood) Examine() string {
	return "a floor made of " + woodTypeInfo[f.Type].Name + "."
}

func (f *FloorWood) Serialize() *NetworkedObject {
	return &NetworkedObject{
		Name:   f.Name(),
		Sprite: "floor_wood",
		Colors: []Color{woodTypeInfo[f.Type].Color},
	}
}

func (f *FloorWood) Blocking() bool {
	return false
}

func (f *FloorWood) ZIndex() int {
	return -50
}
