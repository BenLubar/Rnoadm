package critter

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"sync"
)

type Critter struct {
	world.CombatObject

	modules []world.Module

	kind   CritterType
	colors []string

	facing         uint   // not saved
	animation      string // not saved
	animationTicks uint   // not saved

	mtx sync.Mutex
}

func init() {
	world.Register("critter", world.Combat((*Critter)(nil)))

	world.RegisterSpawnFunc(func(s string) world.Visible {
		living, s := world.SpawnModules(s)

		for t := CritterType(0); t < critterTypeCount; t++ {
			if s == t.Name() {
				return &Critter{
					CombatObject: world.CombatObject{
						LivingObject: living,
					},
					kind:   t,
					colors: t.GenerateColors(),
				}
			}
		}
		return nil
	})
}

func (c *Critter) Save() (uint, interface{}, []world.ObjectLike) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	attached := []world.ObjectLike{&c.CombatObject}

	return 0, map[string]interface{}{
		"k": uint64(c.kind),
		"c": c.colors,
	}, attached
}

func (c *Critter) Load(version uint, data interface{}, attached []world.ObjectLike) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		c.CombatObject = *attached[0].(*world.CombatObject)
		c.kind = CritterType(dataMap["k"].(uint64))
		c.colors = dataMap["c"].([]string)

	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}

}

func (c *Critter) Name() string {
	return c.Type().Name()
}

func (c *Critter) Examine() (string, [][][2]string) {
	_, info := c.CombatObject.Examine()

	return c.Type().Examine(), info
}

func (c *Critter) Type() CritterType {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.kind
}

func (c *Critter) Sprite() string {
	return c.Type().Sprite()
}

func (c *Critter) Colors() []string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.colors
}

func (c *Critter) Think() {
	c.CombatObject.Think()

	c.mtx.Lock()
	if c.animationTicks > 0 {
		c.animationTicks--
		if c.animationTicks == 0 {
			c.animation = ""
			if t := c.Position(); t != nil {
				c.mtx.Unlock()
				t.Zone().Update(t, c.Outer())
				return
			}
		}
	}
	c.mtx.Unlock()
}

func (c *Critter) AnimationType() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.animation
}

func (c *Critter) SpritePos() (uint, uint) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	return c.facing, 0
}

func (c *Critter) MaxHealth() uint64 {
	return c.Type().MaxHealth()
}

func (c *Critter) NotifyPosition(old, new *world.Tile) {
	if old == nil || new == nil {
		c.mtx.Lock()
		c.animationTicks = 0
		c.animation = ""
		c.facing = 0
		c.mtx.Unlock()
		return
	}
	ox, oy := old.Position()
	nx, ny := new.Position()

	c.mtx.Lock()
	switch {
	case ox-1 == nx && oy == ny:
		c.animationTicks = 3
		c.animation = "wa" // walk (alternating)
		c.facing = 6
	case ox+1 == nx && oy == ny:
		c.animationTicks = 3
		c.animation = "wa" // walk (alternating)
		c.facing = 9
	case ox == nx && oy-1 == ny:
		c.animationTicks = 3
		c.animation = "wa" // walk (alternating)
		c.facing = 3
	case ox == nx && oy+1 == ny:
		c.animationTicks = 3
		c.animation = "wa" // walk (alternating)
		c.facing = 0
	}
	c.mtx.Unlock()

	new.Zone().Update(new, c.Outer())
}
