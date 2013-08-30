package material

import (
	"github.com/BenLubar/Rnoadm/world"
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

func (f *Forge) Actions(player world.CombatInventoryMessageAdminHUD) []string {
	return append([]string{"smelt"}, f.VisibleObject.Actions(player)...)
}

func (f *Forge) Interact(player world.CombatInventoryMessageAdminHUD, action string) {
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
					world.NewWalkSchedule(x, y, false, 0),
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
			if _, ok := v.(*Ore); ok {
				ores = append(ores, v.NetworkID())
			}
		}
		for _, v := range player.Inventory() {
			if _, ok := v.(*Stone); ok {
				ores = append(ores, v.NetworkID())
			}
		}
		player.SetHUD("forge", map[string]interface{}{
			"O": ores,
		})
	}
}
