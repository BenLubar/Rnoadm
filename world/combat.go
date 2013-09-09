package world

import (
	"fmt"
	"math/big"
	"sync"
)

type Combat interface {
	Living

	Health() *big.Int
	SetHealth(*big.Int)
	MaxHealth() *big.Int
	HealthRegen() *big.Int

	MaxQuality() *big.Int
	Damage() *big.Int
	Armor() *big.Int
	Crit() *big.Int

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
	}

	if o.damaged.Cmp(max) >= 0 {
		o.mtx.Unlock()
		o.Outer().(Combat).Die()
	} else {
		o.mtx.Unlock()
	}
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

func (o *CombatObject) Armor() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) Damage() *big.Int {
	return big.NewInt(500)
}

func (o *CombatObject) Crit() *big.Int {
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

func (o *CombatObject) Hurt(amount *big.Int, attacker Combat) {
	o.mtx.Lock()
	o.combatTicks = 50
	o.damaged.Add(&o.damaged, amount)
	o.mtx.Unlock()
}

func (o *CombatObject) Die() {
	o.Position().Remove(o.Outer())
}

/*func (o *CombatObject) Actions(player PlayerLike) []string {
	if player == o.Outer() {
		return o.LivingObject.Actions(player)
	}
	return append([]string{"assault"}, o.LivingObject.Actions(player)...)
}*/
