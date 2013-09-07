package material

import (
	"image/color"
	"math/big"
)

var defaultDurability = big.NewInt(1000)
var defaultStat = big.NewInt(250)
var zero = big.NewInt(0)

type MaterialData struct {
	name       string
	color0     color.Color
	color1     color.Color
	skin       uint8
	density    uint64   // centigrams per cubic centimeter
	durability *big.Int // resistance to item degradation

	// combination
	power        *big.Int // melee damage, melee armor
	magic        *big.Int // magic damage, magic armor
	agility      *big.Int // attack speed, move speed
	luck         *big.Int // crit chance, resist chance
	intelligence *big.Int // mana, mana regen
	stamina      *big.Int // health, health regen

	// offensive
	melee_damage *big.Int // increases max hit (melee)
	magic_damage *big.Int // increases max hit (magic)
	mana         *big.Int // increases max mana
	mana_regen   *big.Int // increases mana gained per tick
	crit_chance  *big.Int // increases the range considered critical and decreases the range considered noncritical
	attack_speed *big.Int // decreases the number of ticks between auto-attacks

	// defensive
	melee_armor  *big.Int // increases max damage soaking (melee)
	magic_armor  *big.Int // increases max damage soaking (magic)
	health       *big.Int // increases max health
	health_regen *big.Int // increases health gained per tick
	resistance   *big.Int // increases the chance of evading an incoming attack
	move_speed   *big.Int // increases the weight cap
}

func (m *MaterialData) Name() string {
	return m.name
}

func (m *MaterialData) Color() color.Color {
	return m.color0
}

func (m *MaterialData) ExtraColor() color.Color {
	return m.color1
}

func (m *MaterialData) Skin() uint8 {
	return m.skin
}

// Density in centigrams per cubic centimeter.
func (m *MaterialData) Density() uint64 {
	return m.density
}

// Durability increases the amount of wear an item can take before degrading.
func (m *MaterialData) Durability() *big.Int {
	if m.durability != nil {
		return m.durability
	}
	return defaultDurability
}

// Power is a combination stat affecting melee damage and melee armor.
func (m *MaterialData) Power() *big.Int {
	if m.power != nil {
		return m.power
	}
	return defaultStat
}

// Magic is a combination stat affecting magic damage and magic armor.
func (m *MaterialData) Magic() *big.Int {
	if m.magic != nil {
		return m.magic
	}
	return defaultStat
}

// Agility is a combination stat affecting attack speed and movement speed.
func (m *MaterialData) Agility() *big.Int {
	if m.agility != nil {
		return m.agility
	}
	return defaultStat
}

// Luck is a combination stat affecting chance of critical hits and resistance.
func (m *MaterialData) Luck() *big.Int {
	if m.luck != nil {
		return m.luck
	}
	return defaultStat
}

// Intelligence is a combination stat affecting mana and mana regeneration.
func (m *MaterialData) Intelligence() *big.Int {
	if m.intelligence != nil {
		return m.intelligence
	}
	return defaultStat
}

// Stamina is a combination stat affecting health and health regeneration.
func (m *MaterialData) Stamina() *big.Int {
	if m.stamina != nil {
		return m.stamina
	}
	return defaultStat
}

// MeleeDamage increases max hit (melee)
func (m *MaterialData) MeleeDamage() *big.Int {
	if m.melee_damage != nil {
		return m.melee_damage
	}
	return zero
}

// MagicDamage increases max hit (magic)
func (m *MaterialData) MagicDamage() *big.Int {
	if m.magic_damage != nil {
		return m.magic_damage
	}
	return zero
}

// Mana increases max mana
func (m *MaterialData) Mana() *big.Int {
	if m.mana != nil {
		return m.mana
	}
	return zero
}

// ManaRegen increases mana gained per tick
func (m *MaterialData) ManaRegen() *big.Int {
	if m.mana_regen != nil {
		return m.mana_regen
	}
	return zero
}

// CritChance increases the range considered critical and decreases the range
// considered noncritical
func (m *MaterialData) CritChance() *big.Int {
	if m.crit_chance != nil {
		return m.crit_chance
	}
	return zero
}

// AttackSpeed decreases the number of ticks between auto-attacks
func (m *MaterialData) AttackSpeed() *big.Int {
	if m.attack_speed != nil {
		return m.attack_speed
	}
	return zero
}

// MeleeArmor increases max damage soaking (melee)
func (m *MaterialData) MeleeArmor() *big.Int {
	if m.melee_armor != nil {
		return m.melee_armor
	}
	return zero
}

// MagicArmor increases max damage soaking (magic)
func (m *MaterialData) MagicArmor() *big.Int {
	if m.magic_armor != nil {
		return m.magic_armor
	}
	return zero
}

// Health increases max health
func (m *MaterialData) Health() *big.Int {
	if m.health != nil {
		return m.health
	}
	return zero
}

// Health increases health gained per tick
func (m *MaterialData) HealthRegen() *big.Int {
	if m.health_regen != nil {
		return m.health_regen
	}
	return zero
}

// Resistance increases the chance of evading an incoming attack
func (m *MaterialData) Resistance() *big.Int {
	if m.resistance != nil {
		return m.resistance
	}
	return zero
}

// MovementSpeed increases the weight cap
func (m *MaterialData) MovementSpeed() *big.Int {
	if m.move_speed != nil {
		return m.move_speed
	}
	return zero
}
