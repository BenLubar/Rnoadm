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
	earth = element{
		Name:  "Earth",
		Links: []*element{&dust, &mud, &lava},
	}
	fire = element{
		Name:  "Fire",
		Links: []*element{&smoke, &lava, &steam},
	}
	ice = element{
		Name:  "Ice",
		Links: []*element{&water, &mist},
	}
	dust = element{
		Name:  "Dust",
		Links: []*element{&earth, &air},
	}
	lava = element{
		Name:  "Lava",
		Links: []*element{&earth, &fire},
	}
	water = element{
		Name:  "Water",
		Links: []*element{&mist, &steam, &mud, &ice},
	}
	steam = element{
		Name:  "Steam",
		Links: []*element{&water, &fire},
	}
	mist = element{
		Name:  "Mist",
		Links: []*element{&air, &water},
	}
	smoke = element{
		Name:  "Smoke",
		Links: []*element{&air, &fire},
	}
	mud = element{
		Name:  "Mud",
		Links: []*element{&water, &earth},
	}
	time = element{
		Name:  "Time",
		Links: []*element{&earth, &gravity, &void},
	}
	gravity = element{
		Name:  "Gravity",
		Links: []*element{&earth, &water, &time},
	}
	electric = element{
		Name:  "Electric",
		Links: []*element{&air, &light, &void},
	}
	light = element{
		Name:  "Light",
		Links: []*element{&air, &water, &electric},
	}
	dark = element{
		Name:  "Dark",
		Links: []*element{&fire, &smoke, &void},
	}
	void = element{
		Name:  "Void",
		Links: []*element{&dark, &time, &illusion},
	}
	spiritual = element{
		Name:  "Spiritual",
		Links: []*element{&air, &mist, &time},
	}
	chaotic = element{
		Name:  "Chaotic",
		Links: []*element{&light, &dark, &void},
	}
	illusion = element{
		Name:  "Illusion",
		Links: []*element{&void, &time, &gravity},
	}
	elements[Air] = &air
	// TODO: other elements
}
