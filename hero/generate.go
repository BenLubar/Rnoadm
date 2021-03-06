package hero

import (
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"time"
)

func GenerateHero(r *rand.Rand) *Hero {
	return GenerateHeroRace(r, Race(r.Intn(int(raceCount))))
}

func GenerateHeroRace(r *rand.Rand, race Race) *Hero {
	return GenerateHeroOccupation(r, race, race.Occupations()[r.Intn(len(race.Occupations()))])
}

func GenerateHeroOccupation(r *rand.Rand, race Race, occupation Occupation) *Hero {
	hero := &Hero{}
	world.InitObject(hero)
	hero.mtx.Lock()
	defer hero.mtx.Unlock()
	hero.birth = time.Now().UTC()
	hero.race = race
	hero.occupation = occupation
	hero.gender = race.Genders()[r.Intn(len(race.Genders()))]
	hero.skinTone = uint(r.Intn(len(race.SkinTones())))
	switch race {
	case RaceHuman:
		hero.name = GenerateHumanName(r, hero.gender)
	}
	hero.equip(&Equip{
		slot:         SlotShirt,
		kind:         0,
		customColors: []string{randomColor(r)},
	})
	hero.equip(&Equip{
		slot:         SlotPants,
		kind:         0,
		customColors: []string{randomColor(r)},
	})
	hero.equip(&Equip{
		slot: SlotFeet,
		kind: 0,
	})
	return hero
}

func randomColor(r *rand.Rand) string {
	choices := []string{"23456", "6789a", "abcde"}
	choice := choices[r.Intn(len(choices))]
	return string([]byte{
		'#',
		choice[r.Intn(len(choice))],
		choice[r.Intn(len(choice))],
		choice[r.Intn(len(choice))],
	})
}
