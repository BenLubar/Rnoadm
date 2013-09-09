package critter

import (
	"math/big"
	"math/rand"
)

const hex = "0123456789abcdef"

type CritterType uint64

const (
	SlimeMage CritterType = iota
	SlimeBrute
	StickySlime
	Cow

	critterTypeCount
)

var critterInfo = [critterTypeCount]struct {
	name      string
	examine   string
	maxHealth *big.Int
	sprite    string
	genColors func() []string
}{
	SlimeMage: {
		name:      "slime mage",
		examine:   "a slime with some magical abilities.",
		maxHealth: big.NewInt(100),
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
			})}
		},
	},
	SlimeBrute: {
		name:      "slime brute",
		examine:   "a slime that can bench press a whole hero.",
		maxHealth: big.NewInt(250),
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
				hex[rand.Intn(6)],
			})}
		},
	},
	StickySlime: {
		name:      "sticky slime",
		examine:   "a slime that's a little stickier than most.",
		maxHealth: big.NewInt(50),
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
			})}
		},
	},
	Cow: {
		name:      "cow",
		examine:   "it's a cow, i guess.",
		maxHealth: big.NewInt(10),
		sprite:    "critter_slime",
		genColors: func() []string {
			return []string{string([]byte{
				'#',
				hex[rand.Intn(6)],
				hex[rand.Intn(6)+10],
				hex[rand.Intn(6)],
			})}
		},
	},
}

func (t CritterType) Name() string             { return critterInfo[t].name }
func (t CritterType) Examine() string          { return critterInfo[t].examine }
func (t CritterType) MaxHealth() *big.Int      { return critterInfo[t].maxHealth }
func (t CritterType) Sprite() string           { return critterInfo[t].sprite }
func (t CritterType) GenerateColors() []string { return critterInfo[t].genColors() }
