package material

import (
	"github.com/BenLubar/Rnoadm/world"
	"strconv"
	"time"
)

type Forge struct {
	world.VisibleObject
}

func init() {
	world.Register("forge", world.Visible((*Forge)(nil)))
	world.RegisterSpawnFunc(func(s string) world.Visible {
		if s == "forge" {
			return &Forge{}
		}
		return nil
	})
}

func (f *Forge) Name() string {
	return "forge"
}

func (f *Forge) Examine() (string, [][][2]string) {
	_, info := f.VisibleObject.Examine()

	return "smelts ore into metal ingots. not very good for making pizzas.", info
}

func (f *Forge) Sprite() string {
	return "forge"
}

func (f *Forge) SpriteSize() (uint, uint) {
	return 64, 32
}

func (f *Forge) Scale() uint {
	return 2
}

func (f *Forge) Colors() []string {
	return []string{"no"}
}

func (f *Forge) Blocking() bool {
	return true
}

func (f *Forge) ExtraBlock() [][2]int8 {
	return [][2]int8{{-1, 0}, {1, 0}}
}

func (f *Forge) Actions(player world.PlayerLike) []string {
	return append([]string{"smelt"}, f.VisibleObject.Actions(player)...)
}

func (f *Forge) Interact(player world.PlayerLike, action string) {
	switch action {
	default:
		f.VisibleObject.Interact(player, action)
	case "smelt":
		pos := f.Position()
		ppos := player.Position()
		if pos == nil || ppos == nil {
			return
		}
		x, y := pos.Position()
		if px, py := ppos.Position(); px != x || py != y+1 {
			player.SetSchedule(&world.ScheduleSchedule{
				Schedules: []world.Schedule{
					world.NewWalkSchedule(x, y, false, uint(player.Weight()/player.WeightMax())),
					&world.ActionSchedule{
						Target: f,
						Action: action,
					},
				},
			})
			return
		}
		ores := []interface{}{}
		for _, v := range player.Inventory() {
			if o, ok := v.(*Ore); ok {
				ores = append(ores, map[string]interface{}{
					"i": o.NetworkID(),
					"q": o.Material().Quality(),
				})
			}
		}
		for _, v := range player.Inventory() {
			if s, ok := v.(*Stone); ok {
				ores = append(ores, map[string]interface{}{
					"i": s.NetworkID(),
					"q": s.Material().Quality(),
				})
			}
		}
		instance := player.Instance(f.Position())
		var done time.Time
		instance.Last(func(last time.Time) time.Time {
			done = last
			return last
		})
		contents := []interface{}{}
		instance.Items(func(items []world.Visible) []world.Visible {
			for _, i := range items {
				var sprites []map[string]interface{}
				for j, c := range i.Colors() {
					sprites = append(sprites, map[string]interface{}{
						"S": i.Sprite(),
						"C": c,
						"E": map[string]interface{}{
							"y": j,
						},
					})
				}
				contents = append(contents, sprites)
				// 4 minutes per 5kg of ore
				done = done.Add(time.Minute * 4 * time.Duration(i.(world.Item).Weight()) / 5000)
			}
			return items
		})
		player.SetHUD("forge", map[string]interface{}{
			"O": ores,
			"C": contents,
			"T": done,
		})
	}
}

func (f *Forge) Command(player world.PlayerLike, data map[string]interface{}) {
	defer func() {
		recover() // ignore errors caused by malformed packets.
	}()
	instance := player.Instance(f.Position())
	switch data["A"].(string) {
	case "a":
		id, err := strconv.ParseUint(data["I"].(string), 10, 64)
		if err != nil {
			return
		}
		for _, item := range player.Inventory() {
			if item.NetworkID() == id && player.RemoveItem(item) {
				instance.Items(func(items []world.Visible) []world.Visible {
					return append(items, item)
				})
				instance.Last(func(t time.Time) time.Time {
					return time.Now().UTC()
				})
				f.Interact(player, "smelt")
				break
			}
		}
	}
}
