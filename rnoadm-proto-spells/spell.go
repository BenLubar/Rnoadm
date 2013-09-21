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

	// TODO(BenLubar): generate random spell

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
				b = append(append(append(append(b, i.String()...), ": "...), n.String()...), '\n')
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
