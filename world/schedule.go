package world

type Schedule interface {
	ObjectLike

	// Act returns the number of ticks the object should wait before
	// calling Act again. 0 means to call Act on the next tick.
	Act(Living) uint

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

func (s *cancelSchedule) Act(o Living) uint {
	o.ClearSchedule()
	return 0
}

func (s *cancelSchedule) ShouldSave() bool {
	return true
}
