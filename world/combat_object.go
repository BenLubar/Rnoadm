package world

import (
	"fmt"
	"sync"
)

type Combat interface {
	Living

	Health() uint64
	SetHealth(uint64)
	MaxHealth() uint64

	MaxDamage() uint64

	Accuracy() uint64
	Armor() uint64

	Hurt(uint64, Combat)
	Die()
}

type CombatObject struct {
	LivingObject

	damaged uint64

	mtx sync.Mutex
}

func init() {
	Register("combatobj", Combat((*CombatObject)(nil)))
}

func (o *CombatObject) Save() (uint, interface{}, []ObjectLike) {
	return 0, map[string]interface{}{
		"d": o.damaged,
	}, []ObjectLike{&o.LivingObject}
}

func (o *CombatObject) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		o.LivingObject = *attached[0].(*LivingObject)
		o.damaged = dataMap["d"].(uint64)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (o *CombatObject) Think() {
	o.LivingObject.Think()

	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	damaged := o.damaged
	if o.damaged > 0 && o.damaged < max {
		o.damaged--
	}
	o.mtx.Unlock()

	if damaged >= max {
		o.Outer().(Combat).Die()
	}
}

func (o *CombatObject) Health() uint64 {
	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	damaged := o.damaged
	o.mtx.Unlock()

	if damaged >= max {
		return 0
	}

	return max - damaged
}

func (o *CombatObject) SetHealth(health uint64) {
	max := o.Outer().(Combat).MaxHealth()

	o.mtx.Lock()
	if max < health {
		o.damaged = 0
	} else {
		o.damaged = max - health
	}
	o.mtx.Unlock()
}

func (o *CombatObject) MaxHealth() uint64 {
	return 100
}

func (o *CombatObject) MaxDamage() uint64 {
	return 10
}

func (o *CombatObject) Accuracy() uint64 {
	return 100
}
func (o *CombatObject) Armor() uint64 {
	return 100
}

func (o *CombatObject) Hurt(amount uint64, attacker Combat) {
	o.mtx.Lock()
	o.damaged += amount
	o.mtx.Unlock()
}

func (o *CombatObject) Die() {
	o.Position().Remove(o.Outer())
}

func (o *CombatObject) Actions() []string {
	return append(o.LivingObject.Actions(), "assault")
}
