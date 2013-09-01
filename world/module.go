package world

import (
	"sync"
)

func SpawnModules(s string) (LivingObject, string) {
	var o LivingObject

spawnModules:
	for {
		for name, f := range moduleByName {
			if len(s) > len(name) && s[:len(name)] == name && s[len(name)] == ' ' {
				o.modules = append(o.modules, f())
				s = s[len(name) + 1:]
				continue spawnModules
			}
		}
		return o, s
	}
}

var moduleByName = make(map[string]func() Module)

type Module interface {
	ObjectLike

	ChooseSchedule() Schedule

	notifyOwner(Living)
	Owner() Living
	NotifyOwner(Living)
}

// BaseModule implements the Module interface.
type BaseModule struct {
	Object

	owner Living // not saved

	mtx sync.Mutex
}

func init() {
	Register("module", Module((*BaseModule)(nil)))
}

func (m *BaseModule) ChooseSchedule() Schedule {
	return nil
}

func (m *BaseModule) notifyOwner(owner Living) {
	m.mtx.Lock()
	if m.owner == owner {
		m.mtx.Unlock()
		return
	}
	m.owner = owner
	m.mtx.Unlock()

	m.NotifyOwner(owner)
}

func (m *BaseModule) Owner() Living {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.owner
}

func (m *BaseModule) NotifyOwner(owner Living) {
	// do nothing
}
