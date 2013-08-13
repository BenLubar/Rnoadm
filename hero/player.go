package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"sync"
)

type Player struct {
	Hero
	characterCreation bool
	ancestry          PlayerAncestry
	zoneX, zoneY      int64
	tileX, tileY      uint8
}

var _ world.NoSaveObject = (*Player)(nil)

func init() {
	world.Register("player", HeroLike((*Player)(nil)))
}

func (p *Player) NoSaveObject() {}

func (p *Player) Save() (uint, interface{}, []world.ObjectLike) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	t := p.Position()
	zx, zy := p.zoneX, p.zoneY
	tx, ty := p.tileX, p.tileY

	if t != nil {
		tx, ty = t.Position()
		z := t.Zone()
		zx, zy = z.X, z.Y
	}

	return 0, map[string]interface{}{
		"creation": p.characterCreation,
		"tx":       tx,
		"ty":       ty,
		"zx":       zx,
		"zy":       zy,
	}, []world.ObjectLike{&p.Hero, &p.ancestry}
}

func (p *Player) Load(version uint, data interface{}, attached []world.ObjectLike) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		p.Hero = *attached[0].(*Hero)
		p.ancestry = *attached[1].(*PlayerAncestry)
		p.characterCreation = dataMap["creation"].(bool)
		p.zoneX, p.zoneY = dataMap["zx"].(int64), dataMap["zy"].(int64)
		p.tileX, p.tileY = dataMap["tx"].(uint8), dataMap["ty"].(uint8)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

type PlayerAncestry struct {
	world.Object
	ancestors []*Hero
	mtx       sync.Mutex
}

func init() {
	world.Register("ancestry", world.ObjectLike((*PlayerAncestry)(nil)))
}

func (a *PlayerAncestry) Save() (uint, interface{}, []world.ObjectLike) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	objects := []world.ObjectLike{&a.Object}
	for _, h := range a.ancestors {
		objects = append(objects, h)
	}

	return 0, uint(0), objects
}

func (a *PlayerAncestry) Load(version uint, data interface{}, attached []world.ObjectLike) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	switch version {
	case 0:
		a.Object = *attached[0].(*world.Object)
		attached = attached[1:]
		a.ancestors = make([]*Hero, len(attached))
		for i, h := range attached {
			a.ancestors[i] = h.(*Hero)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
