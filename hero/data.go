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
	sprite      string
	baseHealth  uint64
	genders     []Gender
	occupations []Occupation
}{
	RaceHuman: {
		sprite:      "body_human",
		baseHealth:  10000,
		genders:     []Gender{GenderMale, GenderFemale},
		occupations: []Occupation{OccupationCitizen, OccupationKnight},
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

func (r Race) Sprite() string            { return raceInfo[r].sprite }
func (r Race) BaseHealth() uint64        { return raceInfo[r].baseHealth }
func (r Race) Genders() []Gender         { return raceInfo[r].genders }
func (r Race) Occupations() []Occupation { return raceInfo[r].occupations }
