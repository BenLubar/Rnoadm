package material

import (
	"github.com/BenLubar/Rnoadm/world"
)

type Wall struct {
	world.VisibleObject
}

func init() {
	world.Register("wall", world.Visible((*Wall)(nil)))
	world.RegisterSpawnFunc(func(s string) world.Visible {
		if s == "wall" {
			return world.InitObject(&Wall{}).(world.Visible)
		}
		return nil
	})
}

func (w *Wall) Name() string {
	return "wall"
}

func (w *Wall) Sprite() string {
	return "wall"
}

func (w *Wall) Colors() []string {
	return []string{"#888"}
}

func containsWall(objects []world.ObjectLike) bool {
	for _, o := range objects {
		if _, ok := o.(IsWall); ok {
			return true
		}
	}
	return false
}

type IsWall interface {
	isWall() IsWall
}

func (w *Wall) isWall() IsWall {
	return w
}

func (w *Wall) SpritePos() (uint, uint) {
	var sides uint
	if pos := w.Position(); pos != nil {
		x, y := pos.Position()
		z := pos.Zone()

		if x > 0 && containsWall(z.Tile(x-1, y).Objects()) {
			sides |= 1
		}
		if x < 255 && containsWall(z.Tile(x+1, y).Objects()) {
			sides |= 2
		}
		if y > 0 && containsWall(z.Tile(x, y-1).Objects()) {
			sides |= 4
		}
		if y < 255 && containsWall(z.Tile(x, y+1).Objects()) {
			sides |= 8
		}
	}
	return sides, 0
}

func (w *Wall) SpriteSize() (uint, uint) {
	return 32, 48
}

func (w *Wall) Blocking() bool {
	return true
}
