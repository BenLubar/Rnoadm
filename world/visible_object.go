package world

type Visible interface {
	ObjectLike

	// Name is the user-visible name of this object. Title Case is used for
	// important living objects. All other objects are fully lowercase.
	Name() string

	Examine() string

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
}

type VisibleObject struct {
	Object
}

func init() {
	Register("visobj", Visible((*VisibleObject)(nil)))
}

func (o *VisibleObject) Name() string {
	return "unknown"
}

func (o *VisibleObject) Examine() string {
	return "what could it be?"
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
