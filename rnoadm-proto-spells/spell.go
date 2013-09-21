package main

import (
	"math/big"
	"math/rand"
	"time"
)

type BaseSpell uint64

const (
	SpellChangeHealth BaseSpell = iota

	BaseSpellCount
)

var baseSpellNames [BaseSpellCount]string = [...]string{
	"change health",
}

func (b BaseSpell) String() string {
	return baseSpellNames[b]
}

type Spell struct {
	SelfDirect   map[BaseSpell]*big.Int
	SelfEffect   map[BaseSpell]*big.Int
	TargetDirect map[BaseSpell]*big.Int
	TargetEffect map[BaseSpell]*big.Int

	Warmup   time.Duration
	Channel  time.Duration
	Effect   time.Duration
	Cooldown time.Duration
}

func NewSpell(r *rand.Rand) *Spell {
	s := &Spell{}

	// randomly executes one or more of the argument functions.
	bits := func(funcs ...func()) {
		// bitfield. lowest bit corresponds to first argument.
		run := r.Intn(1<<uint(len(funcs))-1) + 1

		for i, f := range funcs {
			if (run>>uint(i))&1 == 1 {
				f()
			}
		}
	}

	// returns a random time.Duration within [min, max).
	randTime := func(min, max time.Duration) time.Duration {
		return time.Duration(r.Int63n(int64(max-min+1))) + min
	}

	bits(func() {
		// warmup
		s.Warmup = randTime(time.Second/5, 5*time.Second)
	}, func() {
		// cooldown
		s.Cooldown = randTime(time.Second/5, 5*time.Second)
	})

	initStats := func(m map[BaseSpell]*big.Int, negative bool) {
		max := big.NewInt(100000)

		var statFuncs [BaseSpellCount]func()
		for i := range statFuncs {
			statFuncs[i] = func(s BaseSpell) func() {
				return func() {
					m[s] = (&big.Int{}).Rand(r, max)
					if negative {
						m[s].Neg(m[s])
					}
				}
			}(BaseSpell(i))
		}

		bits(statFuncs[:]...)
	}

	bits(func() {
		// effect
		s.Effect = randTime(time.Second/5, 5*time.Second)

		bits(func() {
			// self
			s.SelfEffect = make(map[BaseSpell]*big.Int)
			initStats(s.SelfEffect, false)
		}, func() {
			// target
			s.TargetEffect = make(map[BaseSpell]*big.Int)
			initStats(s.TargetEffect, true)
		})
	}, func() {
		if r.Intn(2) < 1 {
			// channel
			s.Channel = randTime(time.Second/5, 5*time.Second)
		} // else instant

		bits(func() {
			// self
			s.SelfDirect = make(map[BaseSpell]*big.Int)
			initStats(s.SelfDirect, false)
		}, func() {
			// target
			s.TargetDirect = make(map[BaseSpell]*big.Int)
			initStats(s.TargetDirect, true)
		})
	})

	return s
}

func (s *Spell) String() string {
	var b []byte
	if s.Warmup != 0 {
		b = append(append(append(b, "Cast time: "...), s.Warmup.String()...), '\n')
	}
	if s.Channel != 0 {
		b = append(append(append(b, "Channel time: "...), s.Channel.String()...), '\n')
	}
	if s.Warmup == 0 && s.Channel == 0 {
		b = append(b, "Instant cast\n"...)
	}
	if s.Cooldown != 0 {
		b = append(append(append(b, "Cooldown: "...), s.Cooldown.String()...), '\n')
	}

	in_order := func(m map[BaseSpell]*big.Int) {
		for i := BaseSpell(0); i < BaseSpellCount; i++ {
			if n, ok := m[i]; ok {
				b = append(append(append(append(b, i.String()...), ": "...), CommaPlus(n)...), '\n')
			}
		}
	}

	if len(s.SelfDirect) != 0 {
		b = append(b, "\nOn self:\n"...)
		in_order(s.SelfDirect)
	}
	if len(s.SelfEffect) != 0 {
		b = append(append(append(b, "\nApply effect on self ("...), s.Effect.String()...), "):\n"...)
		in_order(s.SelfEffect)
	}
	if len(s.TargetDirect) != 0 {
		b = append(b, "\nOn target:\n"...)
		in_order(s.TargetDirect)
	}
	if len(s.TargetEffect) != 0 {
		b = append(append(append(b, "\nApply effect on target ("...), s.Effect.String()...), "):\n"...)
		in_order(s.TargetEffect)
	}

	return string(b)
}
