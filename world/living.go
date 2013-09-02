package world

import (
	"fmt"
	"sync"
)

type Living interface {
	Visible

	HasSchedule() bool
	SetSchedule(Schedule)
	ClearSchedule()
}

type LivingObject struct {
	VisibleObject
	modules []Module

	delay    uint
	schedule Schedule

	mtx sync.Mutex
}

func init() {
	Register("livingobj", Living((*LivingObject)(nil)))
}

func (o *LivingObject) Save() (uint, interface{}, []ObjectLike) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	attached := []ObjectLike{o.schedule}
	if o.schedule == nil || !o.schedule.ShouldSave() {
		attached[0] = &cancelSchedule{}
	}

	for _, m := range o.modules {
		attached = append(attached, m)
	}

	return 1, map[string]interface{}{
		"d": o.delay,
		"m": uint(len(o.modules)),
	}, attached
}

func (o *LivingObject) Load(version uint, data interface{}, attached []ObjectLike) {
	switch version {
	case 0:
		data.(map[string]interface{})["m"] = uint(0)
		attached = attached[1:]
	case 1:
		dataMap := data.(map[string]interface{})
		o.delay = dataMap["d"].(uint)
		o.schedule = attached[0].(Schedule)
		o.modules = make([]Module, dataMap["m"].(uint))
		for i := range o.modules {
			o.modules[i] = attached[i+1].(Module)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (o *LivingObject) HasSchedule() bool {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	return o.schedule != nil
}

func (o *LivingObject) SetSchedule(s Schedule) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.schedule = s
}

func (o *LivingObject) ClearSchedule() {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.schedule = nil
}

func (o *LivingObject) Think() {
	o.VisibleObject.Think()

	if o.delay > 0 {
		o.delay--
		return
	}

	o.mtx.Lock()
	sched := o.schedule
	o.mtx.Unlock()

	if sched == nil {
		o.mtx.Lock()
		defer o.mtx.Unlock()

		for _, m := range o.modules {
			m.notifyOwner(o.Outer().(Living))
			if s := m.ChooseSchedule(); s != nil {
				o.schedule = s
				return
			}
		}
		o.schedule = nil
		return
	}

	var ok bool
	o.delay, ok = sched.Act(o.Outer().(Living))
	if !ok {
		o.ClearSchedule()
	}
}

func (o *LivingObject) Actions(player PlayerLike) []string {
	if player == o.Outer() {
		return o.VisibleObject.Actions(player)
	}
	return append([]string{"follow"}, o.VisibleObject.Actions(player)...)
}

func (o *LivingObject) Interact(player PlayerLike, action string) {
	switch action {
	default:
		o.VisibleObject.Interact(player, action)

	case "follow":
		pos := o.Position()
		if pos == nil {
			return
		}
		x, y := pos.Position()
		player.SetSchedule(&ScheduleSchedule{
			Schedules: []Schedule{
				NewWalkSchedule(x, y, true, uint(player.Weight()/player.WeightMax())),
				&DelaySchedule{Delay: 1},
				&ActionSchedule{Target: o.Outer().(Living), Action: action},
			},
		})
	}
}
