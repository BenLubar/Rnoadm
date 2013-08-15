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
	skinTones   []string
}{
	RaceHuman: {
		sprite:      "body_human",
		baseHealth:  10000,
		genders:     []Gender{GenderMale, GenderFemale},
		occupations: []Occupation{OccupationCitizen, OccupationKnight},
		skinTones:   []string{"#ffe3cc", "#ffdbbd", "#ffcda3", "#f7e9dc", "#edd0b7", "#e8d1be", "#e5c298", "#e3c3a8", "#c9a281", "#c2a38a", "#ba9c82", "#ad8f76", "#a17a5a", "#876d58", "#6e5948", "#635041", "#4f3f33"},
	},
}

var genderInfo = [genderCount]struct {
	name string
}{
	GenderMale: {
		name: "male",
	},
	GenderFemale: {
		name: "female",
	},
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
func (r Race) SkinTones() []string       { return raceInfo[r].skinTones }
func (g Gender) Name() string            { return genderInfo[g].name }
