package world

import (
	"math/big"
)

var (
	TuningMaxStat         = big.NewInt(1000)
	TuningMetaStatDivisor = big.NewInt(2)

	TuningMinCrit       = big.NewInt(-5)
	TuningMaxCrit       = big.NewInt(35)
	TuningCritDivisor   = big.NewInt(2)
	TuningResistDivisor = big.NewInt(4)

	TuningDamageMiss1    = big.NewInt(10)
	TuningDamageHit      = big.NewInt(85)
	TuningDamageMiss2    = big.NewInt(95)
	TuningDamageMax      = big.NewInt(100)
	TuningCritMultiplier = big.NewInt(3)

	TuningHealthMultiplier            = big.NewInt(5)
	TuningHealthRegenDivisorCombat    = big.NewInt(500)
	TuningHealthRegenDivisorNonCombat = big.NewInt(100)

	TuningManaMultiplier            = big.NewInt(4)
	TuningManaRegenDivisorCombat    = big.NewInt(500)
	TuningManaRegenDivisorNonCombat = big.NewInt(100)
)
