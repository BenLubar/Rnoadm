package world

import (
	"math/big"
)

const TuningEquipSlotCount = 14

var (
	TuningMaxStat         = big.NewInt(1000)
	TuningMetaStatDivisor = big.NewInt(2)

	TuningMinCrit       = big.NewInt(-5)
	TuningMaxCrit       = big.NewInt(35)
	TuningCritDivisor   = big.NewInt(20 * TuningEquipSlotCount)
	TuningResistDivisor = big.NewInt(4000 * TuningEquipSlotCount)

	TuningDamageMiss1    = int64(10)
	TuningDamageHit      = int64(85)
	TuningDamageMiss2    = int64(95)
	TuningDamageMax      = big.NewInt(100)
	TuningCritMultiplier = big.NewInt(3)

	TuningHealthMultiplier            = big.NewInt(5)
	TuningHealthRegenDivisorCombat    = big.NewInt(500)
	TuningHealthRegenDivisorNonCombat = big.NewInt(100)

	TuningManaMultiplier            = big.NewInt(4)
	TuningManaRegenDivisorCombat    = big.NewInt(500)
	TuningManaRegenDivisorNonCombat = big.NewInt(100)

	TuningDefaultStatQuality = big.NewInt(1)

	TuningDefaultStat = map[Stat]*big.Int{
		StatMana:        big.NewInt(500),
		StatManaRegen:   big.NewInt(500),
		StatHealth:      big.NewInt(1000),
		StatHealthRegen: big.NewInt(500),
		StatMeleeDamage: big.NewInt(500),
		StatMeleeArmor:  big.NewInt(500),
		StatCritChance:  big.NewInt(100),
		StatResistance:  big.NewInt(100),
	}

	TuningNodeScoreMultiplier = big.NewInt(500)
	TuningMaxGatherVolume     = big.NewInt(10000)
)
