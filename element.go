package main

import (
	"math/rand"
)

type element struct {
	Name  string
	Links []Element
	Rocks []RockType
	Flora []FloraType
}

type Element uint8

const (
	Air Element = iota
	Earth
	Fire
	Ice
	Nature
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

var elements [elementCount]element
var weakness = map[Element]Element{
	Air:       Mud,
	Water:     Electric,
	Ice:       Fire,
	Earth:     Gravity,
	Fire:      Water,
	Dust:      Time,
	Lava:      Mist,
	Steam:     Lava,
	Mist:      Smoke,
	Smoke:     Light,
	Mud:       Steam,
	Time:      Ice,
	Gravity:   Air,
	Electric:  Earth,
	Light:     Void,
	Dark:      Dust,
	Void:      Spiritual,
	Spiritual: Illusion,
	Illusion:  Dark,
}

func init() {
	elements[Air] = element{
		Name:  "Air",
		Links: []Element{Smoke, Dust, Mist},
	}
	elements[Earth] = element{
		Name:  "Earth",
		Links: []Element{Dust, Mud, Lava},
	}
	elements[Fire] = element{
		Name:  "Fire",
		Links: []Element{Smoke, Lava, Steam},
	}
	elements[Ice] = element{
		Name:  "Ice",
		Links: []Element{Water, Mist},
	}
	elements[Nature] = element{
		Name:  "Nature",
		Links: []Element{Air, Water, Earth},
		Rocks: []RockType{Granite},
		Flora: []FloraType{0},
	}
	elements[Dust] = element{
		Name:  "Dust",
		Links: []Element{Earth, Air},
	}
	elements[Lava] = element{
		Name:  "Lava",
		Links: []Element{Earth, Fire},
	}
	elements[Water] = element{
		Name:  "Water",
		Links: []Element{Mist, Steam, Mud, Ice},
	}
	elements[Steam] = element{
		Name:  "Steam",
		Links: []Element{Water, Fire},
	}
	elements[Mist] = element{
		Name:  "Mist",
		Links: []Element{Air, Water},
	}
	elements[Smoke] = element{
		Name:  "Smoke",
		Links: []Element{Air, Fire},
	}
	elements[Mud] = element{
		Name:  "Mud",
		Links: []Element{Water, Earth},
	}
	elements[Time] = element{
		Name:  "Time",
		Links: []Element{Earth, Gravity, Void},
	}
	elements[Gravity] = element{
		Name:  "Gravity",
		Links: []Element{Earth, Water, Time},
	}
	elements[Electric] = element{
		Name:  "Electric",
		Links: []Element{Air, Light, Void},
	}
	elements[Light] = element{
		Name:  "Light",
		Links: []Element{Air, Water, Electric, Spiritual},
	}
	elements[Dark] = element{
		Name:  "Dark",
		Links: []Element{Fire, Smoke, Void, Spiritual},
	}
	elements[Void] = element{
		Name:  "Void",
		Links: []Element{Dark, Time, Illusion},
	}
	elements[Spiritual] = element{
		Name:  "Spiritual",
		Links: []Element{Air, Mist, Time},
	}
	elements[Chaotic] = element{
		Name: "Chaotic",
		// Special: Chaotic zones contain all element types.
	}
	elements[Illusion] = element{
		Name:  "Illusion",
		Links: []Element{Void, Time, Gravity},
	}
}

func (e Element) Linked(r *rand.Rand) Element {
	if e == Chaotic {
		return Element(r.Intn(int(elementCount)))
	}
	el := elements[e].Links
	i := r.Intn(len(el) + 3)
	if i == len(el) {
		return Nature
	}
	if i > len(el) {
		return e
	}
	return el[i]

}

func (e Element) Rock(r *rand.Rand) (RockType, bool) {
	rocks := elements[e].Rocks
	if len(rocks) == 0 {
		return 0, false
	}
	return rocks[r.Intn(len(rocks))], true
}

func (e Element) Flora(r *rand.Rand) (FloraType, bool) {
	flora := elements[e].Flora
	if len(flora) == 0 {
		return 0, false
	}
	return flora[r.Intn(len(flora))], true
}
