package hero

type Race uint64
type Gender uint64
type Occupation uint64

const (
	RaceHuman Race = iota

	raceCount
)

const (
	GenderMale Gender = iota
	GenderFemale

	genderCount
)

const (
	OccupationAdventurer Occupation = iota
	OccupationCitizen
	OccupationKnight

	occupationCount
)

var raceInfo = [raceCount]struct {
	sprite string
}{
	RaceHuman: {
		sprite: "body_human",
	},
}

var genderInfo = [genderCount]struct {
}{
	GenderMale:   {},
	GenderFemale: {},
}

var occupationInfo = [occupationCount]struct {
}{
	OccupationAdventurer: {},
	OccupationCitizen:    {},
	OccupationKnight:     {},
}

func (r Race) Sprite() string { return raceInfo[r].sprite }
