package main

type element struct {
	Name  string
	Links []*element
}

var (
	_ element
	air, earth, fire, ice
	dust, lava, water, steam, mist, smoke, mud
	time, gravity, electric, light, dark
	void, spiritual, chaotic, illusion
)

type Element uint8

const (
	Air Element = iota
	Earth
	Fire
	Ice
	Dust
	Lava
	Water
	Steam
	Mist
	Smoke
	Mud
	Time
	Gravity
	Electric
	Light
	Dark
	Void
	Spiritual
	Chaotic
	Illusion

	elementCount
)

var elements [elementCount]*Element

func init() {
	air = element{
		Name:  "Air",
		Links: []*element{&smoke, &dust, &mist},
	}
	elements[Air] = &air
	// TODO: other elements
}
