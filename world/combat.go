package world

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
)

type Combat interface {
	Living

	Health() *big.Int
	SetHealth(*big.Int)
	MaxHealth() *big.Int
	HealthRegen() *big.Int

	MaxQuality() *big.Int
	MeleeDamage() *big.Int
	MeleeArmor() *big.Int
	CritChance() *big.Int
	Resistance() *big.Int

	Hurt(*big.Int, Combat)
	Die()
}

type CombatObject struct {
	LivingObject

	damaged big.Int

	combatTicks uint8 // not saved

	mtx sync.Mutex
}

func init() {
	Register("combatobj", Combat((*CombatObject)(nil)))
}

func (o *CombatObject) Save() (uint, interface{}, []ObjectLike) {
	return 1, map[string]interface{}{
		"d": &o.damaged,
	}, []ObjectLike{&o.LivingObject}
}

func (o *CombatObject) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		o.LivingObject = *attached[0].(*LivingObject)
	case 1:
		dataMap := data.(map[string]interface{})
		o.LivingObject = *attached[0].(*LivingObject)
		o.damaged = *dataMap["d"].(*big.Int)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (o *CombatObject) Think() {
	o.LivingObject.Think()

	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	if o.combatTicks > 0 {
		o.combatTicks--
	}

	if o.damaged.Sign() > 0 && o.damaged.Cmp(max) < 0 {
		var regen big.Int
		if o.combatTicks > 0 {
			regen.Div(o.Outer().(Combat).HealthRegen(), TuningHealthRegenDivisorCombat)
		} else {
			regen.Div(o.Outer().(Combat).HealthRegen(), TuningHealthRegenDivisorNonCombat)
		}
		o.damaged.Sub(&o.damaged, &regen)
		if o.damaged.Sign() < 0 {
			o.damaged.SetUint64(0)
		}
		if pos := o.Position(); pos != nil {
			o.mtx.Unlock()
			pos.Zone().Update(pos, o.Outer())
			return
		}
	}

	if o.damaged.Cmp(max) >= 0 {
		o.mtx.Unlock()
		o.Outer().(Combat).Die()
		return
	}
	o.mtx.Unlock()
}

func (o *CombatObject) Health() *big.Int {
	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	defer o.mtx.Unlock()

	var health big.Int

	if o.damaged.Cmp(max) >= 0 {
		return &health
	}

	return health.Sub(max, &o.damaged)
}

func (o *CombatObject) SetHealth(health *big.Int) {
	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	if max.Cmp(health) < 0 {
		// TODO: allow negative damage?
		o.damaged.SetUint64(0)
	} else {
		o.damaged.Sub(max, health)
	}
	o.mtx.Unlock()
}

func (o *CombatObject) MeleeArmor() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) MeleeDamage() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) CritChance() *big.Int {
	return big.NewInt(0)
}

func (o *CombatObject) Resistance() *big.Int {
	return big.NewInt(0)
}

func (o *CombatObject) MaxHealth() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) HealthRegen() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) MaxQuality() *big.Int {
	return big.NewInt(1)
}

var (
	DamageMissed   = big.NewInt(0)
	DamageBlocked  = big.NewInt(0)
	DamageResisted = big.NewInt(0)
)

func (o *CombatObject) Hurt(amount *big.Int, attacker Combat) {
	o.mtx.Lock()
	o.combatTicks = 50
	o.damaged.Add(&o.damaged, amount)
	o.mtx.Unlock()

	if pos := o.Position(); pos != nil {
		pos.Zone().Damage(attacker, o.Outer().(Combat), amount)
	}
}

func (o *CombatObject) Die() {
	o.Position().Remove(o.Outer())
}

func (o *CombatObject) Actions(player PlayerLike) []string {
	if player == o.Outer() {
		return o.LivingObject.Actions(player)
	}
	return append([]string{"fight"}, o.LivingObject.Actions(player)...)
}

func (o *CombatObject) Interact(player PlayerLike, action string) {
	switch action {
	default:
		o.LivingObject.Interact(player, action)
	case "fight":
		if player == o.Outer() {
			return
		}
		player.SetSchedule(&CombatSchedule{Target_: o.Outer().(Combat)})
	}
}

func (o *CombatObject) Examine() (string, [][][2]string) {
	message, info := o.LivingObject.Examine()

	info = append(info, [][2]string{
		{o.Outer().(Combat).Health().String(), "#4fc"},
		{"/", "#ccc"},
		{o.Outer().(Combat).MaxHealth().String(), "#4fc"},
		{" health", "#ccc"},
	})

	return message, info
}

type CombatSchedule struct {
	Object

	Target_ Combat
}

func (s *CombatSchedule) Act(o Living) (uint, bool) {
	c, ok := o.(Combat)
	if !ok {
		return 0, false
	}

	p1, p2 := c.Position(), s.Target_.Position()
	if p1 == nil || p2 == nil || p1.Zone() != p2.Zone() {
		return 0, false
	}
	x1, y1 := p1.Position()
	x2, y2 := p2.Position()
	if (x1 == x2 && y1 == y2) || (x1 != x2 && y1 != y2) || (x1 == x2 && y1 != y2-1 && y1 != y2+1) || (y1 == y2 && x1 != x2-1 && x1 != x2+1) {
		c.SetSchedule(&ScheduleSchedule{Schedules: []Schedule{NewWalkSchedule(x2, y2, true, 0), s}})
		return 0, true
	}

	r := rand.New(rand.NewSource(rand.Int63()))
	maxDamage := c.MeleeDamage()
	if maxDamage.Sign() <= 0 {
		// can't attack
		return 0, false
	}
	damage := (&big.Int{}).Rand(r, maxDamage)
	armor := (&big.Int{}).Set(s.Target_.MeleeArmor())
	if armor.Sign() <= 0 {
		armor.SetUint64(0)
	} else {
		armor.Rand(r, armor)
	}
	crit := (&big.Int{}).Div((&big.Int{}).Div(c.CritChance(), c.MaxQuality()), TuningCritDivisor)
	if crit.Cmp(TuningMinCrit) < 0 {
		crit.Set(TuningMinCrit)
	} else if crit.Cmp(TuningMaxCrit) > 0 {
		crit.Set(TuningMaxCrit)
	}
	crit_ := crit.Int64()
	damage_ := (&big.Int{}).Div((&big.Int{}).Mul(TuningDamageMax, damage), maxDamage).Int64()
	switch {
	case damage_ < TuningDamageMiss1:
		// miss
		damage.SetInt64(0)
	case damage_ < TuningDamageHit+crit_:
		// normal hit
	case damage_ < TuningDamageMiss2+crit_:
		// miss
		damage.SetInt64(0)
	default:
		damage.Mul(damage, TuningCritMultiplier)
	}

	if damage.Sign() <= 0 {
		// miss
		s.Target_.Hurt(DamageMissed, c)
	} else if damage.Cmp(armor) <= 0 {
		// block
		s.Target_.Hurt(DamageBlocked, c)
	} else {
		resistance := s.Target_.Resistance()
		if resistance.Sign() > 0 && (&big.Int{}).Rand(r, resistance).Cmp((&big.Int{}).Mul(s.Target_.MaxQuality(), TuningResistDivisor)) > 0 {
			s.Target_.Hurt(DamageResisted, c)
		} else {
			s.Target_.Hurt(damage.Sub(damage, armor), c)
		}
	}

	// TODO: attack speed
	return 2, true
}

func (s *CombatSchedule) ShouldSave() bool   { return false }
func (s *CombatSchedule) Target() ObjectLike { return s.Target_ }
