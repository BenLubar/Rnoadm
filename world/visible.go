package world

import (
	"github.com/dustin/go-humanize"
	"sync/atomic"
)

type Visible interface {
	ObjectLike

	NetworkID() uint64

	// Name is the user-visible name of this object. Title Case is used for
	// important living objects. All other objects are fully lowercase.
	Name() string

	Examine() (string, [][][2]string)

	// Sprite is sent to the client as the name of the image (minus ".png"
	// extension) to use as a sprite sheet.
	Sprite() string

	// SpriteSize is the height and width of each sprite on this object's
	// sprite sheet.
	SpriteSize() (uint, uint)

	// AnimationType returns a string defined in client code. Empty string
	// means no animation.
	AnimationType() string

	// SpritePos is the position of the base sprite in the sprite sheet.
	// Each color increases the y value on the client, and each animation
	// frame increases the x value.
	SpritePos() (uint, uint)

	Scale() uint

	Colors() []string

	Attached() []Visible

	Blocking() bool

	ExtraBlock() [][2]int8

	Actions(PlayerLike) []string

	Interact(PlayerLike, string)
}

type VisibleObject struct {
	Object

	networkID uint64 // not saved
}

func init() {
	Register("visobj", Visible((*VisibleObject)(nil)))
}

var nextNetworkID uint64

func (o *VisibleObject) NetworkID() uint64 {
	if id := atomic.LoadUint64(&o.networkID); id != 0 {
		// simple case: we already have a network ID; return it
		return id
	}
	// set our network ID to the next available ID, but do nothing if
	// we already have an ID set.
	atomic.CompareAndSwapUint64(&o.networkID, 0, atomic.AddUint64(&nextNetworkID, 1))
	// we definitely have a network ID at this point; return it.
	return atomic.LoadUint64(&o.networkID)
}

func (o *VisibleObject) Name() string {
	return "unknown"
}

func (o *VisibleObject) Examine() (string, [][][2]string) {
	var info [][][2]string

	if item, ok := o.Outer().(Item); ok {
		info = append(info, [][2]string{
			{humanize.Comma(int64(item.Volume())), "#4fc"},
			{" cubic centimeters", "#ccc"},
		}, [][2]string{
			{humanize.Comma(int64(item.Weight())), "#4fc"},
			{" grams", "#ccc"},
		})
	}

	return "what could it be?", info
}

func (o *VisibleObject) Sprite() string {
	return "ui_r1"
}

func (o *VisibleObject) SpriteSize() (uint, uint) {
	return 32, 32
}

func (o *VisibleObject) AnimationType() string {
	return ""
}

func (o *VisibleObject) SpritePos() (uint, uint) {
	return 0, 0
}

func (o *VisibleObject) Scale() uint {
	return 1
}

func (o *VisibleObject) Colors() []string {
	return []string{"#f0f"}
}

func (o *VisibleObject) Attached() []Visible {
	return nil
}

func (o *VisibleObject) Blocking() bool {
	return false
}

func (o *VisibleObject) ExtraBlock() [][2]int8 {
	return nil
}

func (o *VisibleObject) Actions(player PlayerLike) []string {
	var actions []string
	if _, ok := o.Outer().(Item); ok {
		if o.Position() == nil {
			actions = append(actions, "drop")
		} else {
			actions = append(actions, "take")
		}
	}
	actions = append(actions, "examine")
	if player.IsAdmin() {
		actions = append(actions, "impersonate")
	}
	return actions
}

func (o *VisibleObject) Interact(player PlayerLike, action string) {
	switch action {
	case "drop":
		if _, ok := o.Outer().(Item); ok {
			if player.RemoveItem(o.Outer().(Visible)) {
				player.Position().Add(o.Outer())
			}
		}
	case "take":
		pos := o.Position()
		if pos == nil {
			player.SendMessage("somebody else took that!")
			return
		}
		if player.Position() != pos {
			ex, ey := pos.Position()
			player.SetSchedule(&ScheduleSchedule{
				Schedules: []Schedule{
					NewWalkSchedule(ex, ey, false, uint(player.Weight()/player.WeightMax())),
					&ActionSchedule{
						Action:  "take",
						Target_: o.Outer().(Visible),
					},
				},
			})
			return
		}
		if _, ok := o.Outer().(Item); ok {
			if pos.Remove(o.Outer()) {
				if !player.GiveItem(o.Outer().(Visible)) {
					pos.Add(o.Outer())
				}
			}
		}
	case "examine":
		e, i := o.Outer().(Visible).Examine()
		player.SetHUD("examine", map[string]interface{}{
			"N": o.Outer().(Visible).Name(),
			"E": e,
			"I": i,
		})
	case "impersonate":
		if player.IsAdmin() {
			player.Impersonate(o.Outer().(Visible))
		}
	}
}
