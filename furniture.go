package main

type Bed struct {
	Frame WoodType
	Uninteractable
}

func (b *Bed) Name() string {
	return woodTypeInfo[b.Frame].Name + " bed"
}

func (b *Bed) Examine() string {
	return "a bed made of " + woodTypeInfo[b.Frame].Name + " wood."
}

func (b *Bed) Blocking() bool {
	return false
}

func (b *Bed) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Text:   "BED",
		Color:  "#f00",
		ZIndex: 75,
	})
}

func (b *Bed) ZIndex() int {
	return 0
}

type Chest struct {
	Type WoodType
	Uninteractable
}

func (c *Chest) Name() string {
	return woodTypeInfo[c.Type].Name + " chest"
}

func (c *Chest) Examine() string {
	return "a chest made of " + woodTypeInfo[c.Type].Name + " wood."
}

func (c *Chest) Blocking() bool {
	return true
}

func (c *Chest) Paint(x, y int, setcell func(int, int, PaintCell)) {
	setcell(x, y, PaintCell{
		Text:   "CHEST",
		Color:  "#f00",
		ZIndex: 75,
	})
}

func (c *Chest) ZIndex() int {
	return 0
}
