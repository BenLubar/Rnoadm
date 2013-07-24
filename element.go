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
		Name: "Earth",
		Links: []*element{&dust,&mud,&lava),
	}
	fire = element{
		Name: "Fire",
		Links: []*element{&smoke,&lava,&steam),
	}
	// TODO: figure out what the fuck ice is
	dust = element{
		Name: "Dust",
		Links: []*element(&earth, &air),
	}
	lava = element{
		Name: "Lava",
		Links: []*element(&earth, &fire),
	}
	water = element{
		Name: "Water",
		Links: []*element(&mist,&steam,&dust),
	}
	steam = element{
		Name: "Steam",
		Links: []*element(&water,&fire),
	}
	mist = element{
		Name: "Mist",
		Links: []*element(&air,&water),
	}
	smoke = element{
		Name: "Smoke",
		Links: []*element(&air,&fire),
	}
	mud = element{
		Name: "Mud",
		Links: []*element(&water,&earth),
	}
	elements[Air] = &air
	// TODO: other elements
}
