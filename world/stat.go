package world

import (
	"math/big"
)

type Stat uint16

const (
	// combination
	StatPower        Stat = 0<<8 + iota // melee damage, melee armor
	StatMagic                           // magic damage, magic armor
	StatAgility                         // attack speed, move speed
	StatLuck                            // crit chance, resist chance
	StatIntelligence                    // mana, mana regen
	StatStamina                         // health, health regen
	StatIntegrity                       // gathering, structure health

	// offensive
	StatMeleeDamage Stat = 1<<8 + iota // increases max hit (melee)
	StatMagicDamage                    // increases max hit (magic)
	StatMana                           // increases max mana
	StatManaRegen                      // increases mana gained per tick
	StatCritChance                     // increases the range considered critical and decreases the range considered noncritical
	StatAttackSpeed                    // decreases the number of ticks between auto-attacks

	// defensive
	StatMeleeArmor    Stat = 2<<8 + iota // increases max damage soaking (melee)
	StatMagicArmor                       // increases max damage soaking (magic)
	StatHealth                           // increases max health
	StatHealthRegen                      // increases health gained per tick
	StatResistance                       // increases the chance of evading an incoming attack
	StatMovementSpeed                    // increases the weight cap

	// other
	StatGathering       Stat = 3<<8 + iota // increases output from resource nodes
	StatStructureHealth                    // increases durability of structures
)

type StatLike interface {
	Stat(Stat) *big.Int
}
