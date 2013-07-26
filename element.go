package main

type element struct {
	Name  string
	Links []*element
}

var (
	air, earth, fire, ice, nature              element
	dust, lava, water, steam, mist, smoke, mud element
	time_, gravity, electric, light, dark      element
	void, spiritual, chaotic, illusion         element
)

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

var elements [elementCount]*element

func init() {
	air = element{
		Name:  "Air",
		Links: []*element{&smoke, &dust, &mist, &nature},
	}
	elements[Air] = &air
	earth = element{
		Name:  "Earth",
		Links: []*element{&dust, &mud, &lava, &nature},
	}
	elements[Earth] = &earth
	fire = element{
		Name:  "Fire",
		Links: []*element{&smoke, &lava, &steam, &nature},
	}
	elements[Fire] = &fire
	ice = element{
		Name:  "Ice",
		Links: []*element{&water, &mist, &nature},
	}
	elements[Ice] = &ice
	nature = element{
		Name:  "Nature",
		Links: []*element{}, //special: all elements contain nature
	}
	elements[Nature] = &nature
	dust = element{
		Name:  "Dust",
		Links: []*element{&earth, &air, &nature}, 
	}
	elements[Dust] = &dust
	lava = element{
		Name:  "Lava",
		Links: []*element{&earth, &fire, &nature},
	}
	elements[Lava] = &lava
	water = element{
		Name:  "Water",
		Links: []*element{&mist, &steam, &mud, &ice, &nature},
	}
	elements[Water] = &water
	steam = element{
		Name:  "Steam",
		Links: []*element{&water, &fire, &nature},
	}
	elements[Steam] = &steam
	mist = element{
		Name:  "Mist",
		Links: []*element{&air, &water, &nature},
	}
	elements[Mist] = &mist
	smoke = element{
		Name:  "Smoke",
		Links: []*element{&air, &fire, &nature},
	}
	elements[Smoke] = &smoke
	mud = element{
		Name:  "Mud",
		Links: []*element{&water, &earth, &nature},
	}
	elements[Mud] = &mud
	time_ = element{
		Name:  "Time",
		Links: []*element{&earth, &gravity, &void, &nature},
	}
	elements[Time] = &time_
	gravity = element{
		Name:  "Gravity",
		Links: []*element{&earth, &water, &time_, &nature},
	}
	elements[Gravity] = &gravity
	electric = element{
		Name:  "Electric",
		Links: []*element{&air, &light, &void, &nature},
	}
	elements[Electric] = &electric
	light = element{
		Name:  "Light",
		Links: []*element{&air, &water, &electric, &spiritual, &nature},
	}
	elements[Light] = &light
	dark = element{
		Name:  "Dark",
		Links: []*element{&fire, &smoke, &void, &spiritual, &nature},
	}
	elements[Dark] = &dark
	void = element{
		Name:  "Void",
		Links: []*element{&dark, &time_, &illusion, &nature},
	}
	elements[Void] = &void
	spiritual = element{
		Name:  "Spiritual",
		Links: []*element{&air, &mist, &time_, &nature},
	}
	elements[Spiritual] = &spiritual
	chaotic = element{
		Name:  "Chaotic",
		// Special: Chaotic zones contain all element types.
	}
	elements[Chaotic] = &chaotic
	illusion = element{
		Name:  "Illusion",
		Links: []*element{&void, &time_, &gravity, &nature},
	}
	elements[Illusion] = &illusion
}
