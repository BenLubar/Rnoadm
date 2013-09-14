package material

import (
	"github.com/BenLubar/Rnoadm/world"
	"sync"
)

type Door struct {
	world.VisibleObject

	open bool

	mtx sync.Mutex
}

func init() {
	world.Register("door", world.Visible((*Door)(nil)))
	world.RegisterSpawnFunc(func(s string) world.Visible {
		if s == "door" {
			return &Door{}
		}
		return nil
	})
}

func (d *Door) Name() string {
	return "door"
}

func (d *Door) Sprite() string {
	return "door"
}

func (d *Door) Colors() []string {
	return []string{"#888", "#888"}
}

func (d *Door) SpritePos() (uint, uint) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if d.open {
		return 1, 0
	}
	return 0, 0
}

func (d *Door) SpriteSize() (uint, uint) {
	return 48, 64
}

func (d *Door) Blocking() bool {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	return !d.open
}

func (d *Door) Actions(player world.PlayerLike) []string {
	actions := d.VisibleObject.Actions(player)

	d.mtx.Lock()
	defer d.mtx.Unlock()

	if d.Position() != nil {
		if d.open {
			actions = append([]string{"close"}, actions...)
		} else {
			actions = append([]string{"open"}, actions...)
		}
	}
	return actions
}

func (d *Door) Interact(player world.PlayerLike, action string) {
	switch action {
	case "open":
		pos := d.Position()
		if pos == nil {
			return
		}
		x, y := pos.Position()
		px, py := player.Position().Position()
		if (px == x && py != y-1 && py != y+1) || (py == y && px != x-1 && px != x+1) || (px != x && py != y) {
			player.SetSchedule(&world.ScheduleSchedule{
				Schedules: []world.Schedule{
					world.NewWalkSchedule(x, y, true, uint(player.Weight()/player.WeightMax())),
					&world.ActionSchedule{
						Action:  "open",
						Target_: d.Outer().(world.Visible),
					},
				},
			})
			return
		}
		d.mtx.Lock()
		d.open = true
		d.mtx.Unlock()
		pos.Zone().Update(pos, d)
	case "close":
		pos := d.Position()
		if pos == nil {
			return
		}
		x, y := pos.Position()
		px, py := player.Position().Position()
		if (px == x && (py < y-1 || py > y+1)) || (py == y && (px < x-1 || px > x+1)) || (px != x && py != y) {
			player.SetSchedule(&world.ScheduleSchedule{
				Schedules: []world.Schedule{
					world.NewWalkSchedule(x, y, true, uint(player.Weight()/player.WeightMax())),
					&world.ActionSchedule{
						Action:  "close",
						Target_: d.Outer().(world.Visible),
					},
				},
			})
			return
		}
		d.mtx.Lock()
		d.open = false
		d.mtx.Unlock()
		pos.Zone().Update(pos, d)
	default:
		d.VisibleObject.Interact(player, action)
	}
}
