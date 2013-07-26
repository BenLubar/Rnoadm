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
	air = element{
		Name:  "Air",
		Links: []*element{&smoke, &dust, &mist},
	}
	elements[Air] = &air
	earth = element{
		Name:  "Earth",
		Links: []*element{&dust, &mud, &lava},
	}
	elements[Earth] = &earth
	fire = element{
		Name:  "Fire",
		Links: []*element{&smoke, &lava, &steam},
	}
	elements[Fire] = &fire
	ice = element{
		Name:  "Ice",
		Links: []*element{&water, &mist},
	}
	elements[Ice] = &ice
	nature = element{
		Name:  "Nature",
		Links: []*element{&air, &water, &earth},
	}
	elements[Nature] = &nature
	dust = element{
		Name:  "Dust",
		Links: []*element{&earth, &air},
	}
	elements[Dust] = &dust
	lava = element{
		Name:  "Lava",
		Links: []*element{&earth, &fire},
	}
	elements[Lava] = &lava
	water = element{
		Name:  "Water",
		Links: []*element{&mist, &steam, &mud, &ice},
	}
	elements[Water] = &water
	steam = element{
		Name:  "Steam",
		Links: []*element{&water, &fire},
	}
	elements[Steam] = &steam
	mist = element{
		Name:  "Mist",
		Links: []*element{&air, &water},
	}
	elements[Mist] = &mist
	smoke = element{
		Name:  "Smoke",
		Links: []*element{&air, &fire},
	}
	elements[Smoke] = &smoke
	mud = element{
		Name:  "Mud",
		Links: []*element{&water, &earth},
	}
	elements[Mud] = &mud
	time_ = element{
		Name:  "Time",
		Links: []*element{&earth, &gravity, &void},
	}
	elements[Time] = &time_
	gravity = element{
		Name:  "Gravity",
		Links: []*element{&earth, &water, &time_},
	}
	elements[Gravity] = &gravity
	electric = element{
		Name:  "Electric",
		Links: []*element{&air, &light, &void},
	}
	elements[Electric] = &electric
	light = element{
		Name:  "Light",
		Links: []*element{&air, &water, &electric, &spiritual},
	}
	elements[Light] = &light
	dark = element{
		Name:  "Dark",
		Links: []*element{&fire, &smoke, &void, &spiritual},
	}
	elements[Dark] = &dark
	void = element{
		Name:  "Void",
		Links: []*element{&dark, &time_, &illusion},
	}
	elements[Void] = &void
	spiritual = element{
		Name:  "Spiritual",
		Links: []*element{&air, &mist, &time_},
	}
	elements[Spiritual] = &spiritual
	chaotic = element{
		Name: "Chaotic",
		// Special: Chaotic zones contain all element types.
	}
	elements[Chaotic] = &chaotic
	illusion = element{
		Name:  "Illusion",
		Links: []*element{&void, &time_, &gravity},
	}
	elements[Illusion] = &illusion
}
