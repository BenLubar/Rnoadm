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
	sx, sy uint8
	ex, ey uint8
}

func init() {
	Register("walksched", Schedule((*walkSchedule)(nil)))
}

func NewWalkSchedule(sx, sy, ex, ey uint8) Schedule {
	return InitObject(&walkSchedule{
		sx: sx,
		sy: sy,
		ex: ex,
		ey: ey,
	}).(Schedule)
}

func (s *walkSchedule) Act(o Living) (uint, bool) {
	t := o.Position()
	if t == nil {
		return 0, false
	}
	x, y := t.Position()
	if x != s.sx || y != s.sy {
		return 0, false
	}
	dx, dy := int(s.ex)-int(x), int(s.ey)-int(y)
	if dx > 0 {
		if dy > dx {
			s.sy++
		} else if dy < 0 && -dy > dx {
			s.sy--
		} else {
			s.sx++
		}
	} else if dx < 0 {
		if dy > -dx {
			s.sy++
		} else if dy < 0 && dy < dx {
			s.sy--
		} else {
			s.sx--
		}
	} else {
		if dy > 0 {
			s.sy++
		} else if dy < 0 {
			s.sy--
		} else {
			return 0, false
		}
	}
	if t.Zone().Tile(s.sx, s.sy).Blocked() {
		return 0, false
	}
	t.Move(o, t.Zone().Tile(s.sx, s.sy))
	return 2, true
}

func (s *walkSchedule) ShouldSave() bool {
	return true
}

func (s *walkSchedule) Save() (uint, interface{}, []ObjectLike) {
	return 0, map[string]interface{}{
		"sx": s.sx,
		"sy": s.sy,
		"ex": s.ex,
		"ey": s.ey,
	}, nil
}

func (s *walkSchedule) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		sx, sy := dataMap["sx"].(uint8), dataMap["sy"].(uint8)
		ex, ey := dataMap["ex"].(uint8), dataMap["ey"].(uint8)
		*s = *NewWalkSchedule(sx, sy, ex, ey).(*walkSchedule)
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

type ActionSchedule struct {
	Object

	Action string
	Target Visible
}

func (s *ActionSchedule) Act(o Living) (uint, bool) {
	s.Target.Interact(o.(CombatInventoryMessageAdmin), s.Action)
	return 0, false
}

func (s *ActionSchedule) ShouldSave() bool {
	return false
}
