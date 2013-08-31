package world

import (
	"fmt"
	"math/rand"
)

// Stationary - Do not move.
type StationaryModule struct {
	BaseModule
}

func init() {
	Register("modstay", Module((*StationaryModule)(nil)))
}

func (m *StationaryModule) ChooseSchedule() Schedule {
	return &DelaySchedule{Delay: 5}
}

// Wanderer - Target a random nearby tile and move to it, then wait a few ticks.
type WandererModule struct {
	BaseModule
}

func init() {
	Register("modwander", Module((*WandererModule)(nil)))
}

func (m *WandererModule) ChooseSchedule() Schedule {
	const distance = 6

	x8, y8 := m.Owner().Position().Position()
	x, y := int(x8), int(y8)
	x += rand.Intn(distance*2) - distance
	y += rand.Intn(distance*2) - distance
	if x < 0 || x > 255 || y < 0 || y > 255 {
		return nil
	}
	return &ScheduleSchedule{
		Schedules: []Schedule{
			&DelaySchedule{Delay: 5},
			NewWalkSchedule(uint8(x), uint8(y), false, 6),
		},
	}
}

// Guard dog - Move to a specific tile, which does not change location after it
// is set.
type GuardDogModule struct {
	BaseModule

	x, y *uint8
}

func init() {
	Register("modguarddog", Module((*GuardDogModule)(nil)))
}

func (m *GuardDogModule) ChooseSchedule() Schedule {
	x, y := m.Owner().Position().Position()
	if m.x == nil {
		m.x = &x
		m.y = &y
	}
	if x == *m.x && y == *m.y {
		return &DelaySchedule{Delay: 5}
	}
	return NewWalkSchedule(*m.x, *m.y, false, 0)
}

func (m *GuardDogModule) Save() (uint, interface{}, []ObjectLike) {
	if m.x == nil {
		return 0, map[string]interface{}{}, nil
	}
	return 0, map[string]interface{}{
		"x": *m.x,
		"y": *m.y,
	}, nil
}

func (m *GuardDogModule) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		if _, ok := dataMap["x"]; ok {
			x := dataMap["x"].(uint8)
			m.x = &x
			y := dataMap["y"].(uint8)
			m.y = &y
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}
