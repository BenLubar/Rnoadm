package world

import (
	"fmt"
)

type Schedule interface {
	ObjectLike

	// Act returns the number of ticks the object should wait before
	// calling Act again. 0 means to call Act on the next tick.
	// Returning false for the second value cancels the schedule.
	Act(Living) (uint, bool)

	// ShouldSave returns false if this schedule is impossible to save (for
	// example if it requires a pointer to another object).
	ShouldSave() bool
}

type cancelSchedule struct {
	Object
}

func init() {
	Register("nosched", Schedule((*cancelSchedule)(nil)))
}

func (s *cancelSchedule) Act(o Living) (uint, bool) {
	return 0, false
}

func (s *cancelSchedule) ShouldSave() bool {
	return true
}

type walkSchedule struct {
	Object
	ex, ey    uint8
	path      [][2]uint8
	stopEarly bool
	delay     uint
}

func init() {
	Register("walksched", Schedule((*walkSchedule)(nil)))
}

func NewWalkSchedule(ex, ey uint8, stopEarly bool, delay uint) Schedule {
	return InitObject(&walkSchedule{
		ex:        ex,
		ey:        ey,
		stopEarly: stopEarly,
		delay:     delay,
	}).(Schedule)
}

func (s *walkSchedule) Act(o Living) (uint, bool) {
	t := o.Position()
	if t == nil {
		return 0, false
	}
	if s.path == nil {
		s.path = t.Zone().Path(t, t.Zone().Tile(s.ex, s.ey), s.stopEarly)
	}
	if len(s.path) <= 1 {
		return 0, false
	}
	x, y := s.path[0][0], s.path[0][1]
	s.path = s.path[1:]
	if tx, ty := t.Position(); tx != x && ty != y {
		return 0, false
	}
	x, y = s.path[0][0], s.path[0][1]
	if t.Zone().Tile(x, y).Blocked() {
		return 0, false
	}
	t.Move(o, t.Zone().Tile(x, y))
	return 2 + s.delay, true
}

func (s *walkSchedule) ShouldSave() bool {
	return true
}

func (s *walkSchedule) Save() (uint, interface{}, []ObjectLike) {
	return 0, map[string]interface{}{
		"ex": s.ex,
		"ey": s.ey,
		"se": s.stopEarly,
		"d":  s.delay,
	}, nil
}

func (s *walkSchedule) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		s.ex = dataMap["ex"].(uint8)
		s.ey = dataMap["ey"].(uint8)
		s.stopEarly, _ = dataMap["se"].(bool)
		s.delay, _ = dataMap["d"].(uint)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

type ScheduleSchedule struct {
	Object

	Schedules []Schedule
}

func (s *ScheduleSchedule) Act(o Living) (uint, bool) {
	if len(s.Schedules) > 0 {
		delay, keep := s.Schedules[0].Act(o)
		if !keep {
			s.Schedules = s.Schedules[1:]
		}
		return delay, true
	}
	return 0, false
}

func (s *ScheduleSchedule) ShouldSave() bool {
	return false
}

type DelaySchedule struct {
	Object

	Delay uint
}

func (s *DelaySchedule) Act(o Living) (uint, bool) {
	delay := s.Delay
	if delay == 0 {
		return 0, false
	}
	s.Delay = 0
	return delay, true
}

func (s *DelaySchedule) ShouldSave() bool {
	return false
}

type ActionSchedule struct {
	Object

	Action string
	Target Visible
}

func (s *ActionSchedule) Act(o Living) (uint, bool) {
	s.Target.Interact(o.(PlayerLike), s.Action)
	return 0, false
}

func (s *ActionSchedule) ShouldSave() bool {
	return false
}
