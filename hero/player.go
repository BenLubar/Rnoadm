package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Player struct {
	Hero
	characterCreation bool
	ancestry          PlayerAncestry
	zoneX, zoneY      int64
	tileX, tileY      uint8
	instances         PlayerInstances

	login      string
	password   []byte
	admin      bool
	firstAddr  string
	lastAddr   string
	registered time.Time
	lastLogin  time.Time

	kick      chan string
	hud       chan HUD
	inventory chan []world.Visible
	messages  chan []Message

	impersonating world.Visible
}

var _ world.NoSaveObject = (*Player)(nil)
var _ world.AdminLike = (*Player)(nil)
var _ world.SendMessageLike = (*Player)(nil)
var _ world.InventoryLike = (*Player)(nil)

func init() {
	world.Register("player", HeroLike((*Player)(nil)))
}

func (p *Player) SaveSelf() {
	SavePlayer(p)
}

func (p *Player) UpdatePosition() {
	t := p.Position()

	if t != nil {
		p.tileX, p.tileY = t.Position()
		z := t.Zone()
		p.zoneX, p.zoneY = z.X, z.Y
	}
}

func (p *Player) Save() (uint, interface{}, []world.ObjectLike) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	t := p.Position()
	zx, zy := p.zoneX, p.zoneY
	tx, ty := p.tileX, p.tileY

	if t != nil {
		tx, ty = t.Position()
		z := t.Zone()
		zx, zy = z.X, z.Y
	}

	return 1, map[string]interface{}{
		"creation": p.characterCreation,
		"tx":       tx,
		"ty":       ty,
		"zx":       zx,
		"zy":       zy,
		"u":        p.login,
		"p":        p.password,
		"admin":    p.admin,
		"rega":     p.firstAddr,
		"lasta":    p.lastAddr,
		"regt":     p.registered,
		"lastt":    p.lastLogin,
	}, []world.ObjectLike{&p.Hero, &p.ancestry, &p.instances}
}

func (p *Player) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		attached = append(attached, world.InitObject(&PlayerInstances{}))
		var err error
		dataMap["regt"], err = time.Parse(time.RFC3339Nano, dataMap["regt"].(string))
		if err != nil {
			panic(err)
		}
		dataMap["lastt"], err = time.Parse(time.RFC3339Nano, dataMap["lastt"].(string))
		if err != nil {
			panic(err)
		}
		fallthrough
	case 1:
		dataMap := data.(map[string]interface{})
		p.Hero = *attached[0].(*Hero)
		p.mtx.Lock()
		defer p.mtx.Unlock()
		p.ancestry = *attached[1].(*PlayerAncestry)
		p.instances = *attached[2].(*PlayerInstances)
		p.characterCreation = dataMap["creation"].(bool)
		p.zoneX, p.zoneY = dataMap["zx"].(int64), dataMap["zy"].(int64)
		p.tileX, p.tileY = dataMap["tx"].(uint8), dataMap["ty"].(uint8)
		p.login = dataMap["u"].(string)
		p.password = dataMap["p"].([]byte)
		p.admin, _ = dataMap["admin"].(bool)
		p.registered = dataMap["regt"].(time.Time)
		p.firstAddr = dataMap["rega"].(string)
		p.lastLogin = dataMap["lastt"].(time.Time)
		p.lastAddr = dataMap["lasta"].(string)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
	for _, e := range p.Hero.equipped {
		if e.AdminOnly() && !p.admin {
			delete(p.Hero.equipped, e.slot)
		} else {
			e.wearer = p
		}
	}
	if !p.admin {
		for i := 0; i < len(p.Hero.items); i++ {
			if item, ok := p.Hero.items[i].(world.Item); !ok || item.AdminOnly() {
				p.Hero.items = append(p.Hero.items[:i], p.Hero.items[i+1:]...)
				i--
			}
		}
	}
}

type PlayerAncestry struct {
	world.Object
	ancestors []*Hero
	mtx       sync.Mutex
}

func init() {
	world.Register("ancestry", world.ObjectLike((*PlayerAncestry)(nil)))
}

func (a *PlayerAncestry) Save() (uint, interface{}, []world.ObjectLike) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	objects := []world.ObjectLike{}
	for _, h := range a.ancestors {
		objects = append(objects, h)
	}

	return 1, uint(0), objects
}

func (a *PlayerAncestry) Load(version uint, data interface{}, attached []world.ObjectLike) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	switch version {
	case 0:
		attached = attached[1:]
		fallthrough
	case 1:
		a.ancestors = make([]*Hero, len(attached))
		for i, h := range attached {
			a.ancestors[i] = h.(*Hero)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

type PlayerInstances struct {
	world.Object
	instances map[instanceLocation]*PlayerInstance
	mtx       sync.Mutex
}

func init() {
	world.Register("playerinsts", world.ObjectLike((*PlayerInstances)(nil)))
}

func (pi *PlayerInstances) Save() (uint, interface{}, []world.ObjectLike) {
	pi.mtx.Lock()
	defer pi.mtx.Unlock()

	locations := []interface{}{}
	objects := []world.ObjectLike{}
	for l, o := range pi.instances {
		locations = append(locations, l.Save())
		objects = append(objects, o)
	}

	return 0, locations, objects
}

func (pi *PlayerInstances) Load(version uint, data interface{}, attached []world.ObjectLike) {
	pi.mtx.Lock()
	defer pi.mtx.Unlock()

	switch version {
	case 0:
		pi.instances = make(map[instanceLocation]*PlayerInstance, len(attached))
		var loc instanceLocation
		for i, l := range data.([]interface{}) {
			loc.Load(l.(map[string]interface{}))
			pi.instances[loc] = attached[i].(*PlayerInstance)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

type instanceLocation struct {
	ZoneX, ZoneY int64
	TileX, TileY uint8
}

func (i *instanceLocation) Save() map[string]interface{} {
	return map[string]interface{}{
		"zx": i.ZoneX,
		"zy": i.ZoneY,
		"tx": i.TileX,
		"ty": i.TileY,
	}
}

func (i *instanceLocation) Load(data map[string]interface{}) {
	i.ZoneX = data["zx"].(int64)
	i.ZoneY = data["zy"].(int64)
	i.TileX = data["tx"].(uint8)
	i.TileY = data["ty"].(uint8)
}

type PlayerInstance struct {
	world.Object

	items []world.Visible
	last  time.Time

	mtx sync.Mutex
}

func init() {
	world.Register("playerinst", world.ObjectLike((*PlayerInstance)(nil)))
}

func (pi *PlayerInstance) Save() (uint, interface{}, []world.ObjectLike) {
	pi.mtx.Lock()
	defer pi.mtx.Unlock()

	attached := make([]world.ObjectLike, len(pi.items))
	for i, o := range pi.items {
		attached[i] = o
	}

	return 0, map[string]interface{}{
		"last":  pi.last,
		"items": uint(len(pi.items)),
	}, attached
}

func (pi *PlayerInstance) Load(version uint, data interface{}, attached []world.ObjectLike) {
	pi.mtx.Lock()
	defer pi.mtx.Unlock()

	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		pi.last = dataMap["last"].(time.Time)
		pi.items = make([]world.Visible, dataMap["items"].(uint))
		for i := range pi.items {
			pi.items[i] = attached[i].(world.Visible)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (p *Player) Name() string {
	if i := p.impersonate(); i != nil {
		return i.Name()
	}
	return p.Hero.Name()
}

func (p *Player) Examine() (string, [][][2]string) {
	if i := p.impersonate(); i != nil {
		return i.Examine()
	}

	examine, info := p.Hero.Examine()

	if p.IsAdmin() {
		examine = "could this be a Founder?"
	}

	return examine, info
}

func (p *Player) Sprite() string {
	if i := p.impersonate(); i != nil {
		return i.Sprite()
	}
	return p.Hero.Sprite()
}

func (p *Player) SpriteSize() (uint, uint) {
	if i := p.impersonate(); i != nil {
		return i.SpriteSize()
	}
	return p.Hero.SpriteSize()
}

func (p *Player) AnimationType() string {
	if i := p.impersonate(); i != nil {
		return i.AnimationType()
	}
	return p.Hero.AnimationType()
}

func (p *Player) SpritePos() (uint, uint) {
	if i := p.impersonate(); i != nil {
		return i.SpritePos()
	}
	return p.Hero.SpritePos()
}

func (p *Player) Scale() uint {
	if i := p.impersonate(); i != nil {
		return i.Scale()
	}
	return p.Hero.Scale()
}

func (p *Player) Colors() []string {
	if i := p.impersonate(); i != nil {
		return i.Colors()
	}
	return p.Hero.Colors()
}

func (p *Player) Attached() []world.Visible {
	if i := p.impersonate(); i != nil {
		return i.Attached()
	}
	return p.Hero.Attached()
}

func (p *Player) Actions(player world.CombatInventoryMessageAdminHUD) []string {
	if i := p.impersonate(); i != nil {
		return i.Actions(player)
	}
	return p.Hero.Actions(player)
}

func (p *Player) SetHUD(name string, data map[string]interface{}) {
	for {
		select {
		case p.hud <- HUD{
			Name: name,
			Data: data,
		}:
			return
		case <-p.hud:
		}
	}
}

func (p *Player) ClearHUD() {
	p.SetHUD("", nil)
}

func (p *Player) InitPlayer() (kick <-chan string, hud <-chan HUD, inventory <-chan []world.Visible, messages <-chan []Message) {
	p.kick = make(chan string, 1)
	kick = p.kick
	p.hud = make(chan HUD, 1)
	hud = p.hud
	p.inventory = make(chan []world.Visible, 1)
	inventory = p.inventory
	p.messages = make(chan []Message, 1)
	messages = p.messages
	return
}

func (p *Player) Kick(message string) {
	select {
	case p.kick <- message:
	default: // already kicked
	}
}

func (p *Player) CanSpawn() bool {
	p.mtx.Lock()
	characterCreation := p.characterCreation
	admin := p.admin
	p.mtx.Unlock()

	return !characterCreation && (admin || p.Health() > 0)
}

func (p *Player) LoginPosition() (int64, uint8, int64, uint8) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.zoneX, p.tileX, p.zoneY, p.tileY
}

func (p *Player) CharacterCreation(command string) {
	switch command {
	default:
		return
	case "":
	case "gender":
		p.mtx.Lock()
		var index int
		for i, g := range p.race.Genders() {
			if g == p.gender {
				index = (i + 1) % len(p.race.Genders())
				break
			}
		}
		p.gender = p.race.Genders()[index]
		p.mtx.Unlock()
		fallthrough
	case "name":
		p.mtx.Lock()
		switch p.race {
		case RaceHuman:
			p.name = GenerateHumanName(rand.New(rand.NewSource(rand.Int63())), p.gender)
		}
		p.mtx.Unlock()
	case "skin":
		p.mtx.Lock()
		p.skinTone = (p.skinTone + 1) % uint(len(p.race.SkinTones()))
		p.mtx.Unlock()
	case "shirt":
		p.mtx.Lock()
		p.equipped[SlotShirt].customColors[0] = randomColor(rand.New(rand.NewSource(rand.Int63())))
		p.mtx.Unlock()
	case "pants":
		p.mtx.Lock()
		p.equipped[SlotPants].customColors[0] = randomColor(rand.New(rand.NewSource(rand.Int63())))
		p.mtx.Unlock()
	case "accept":
		p.mtx.Lock()
		z := world.GetZone(p.zoneX, p.zoneY)
		p.characterCreation = false
		tile := z.Tile(p.tileX, p.tileY)
		p.mtx.Unlock()
		tile.Add(p)
		world.ReleaseZone(z)
		p.ClearHUD()
		return
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	var sprites []map[string]interface{}

	w, h := p.race.SpriteSize()
	sprites = append(sprites, map[string]interface{}{
		"C": p.race.SkinTones()[p.skinTone],
		"S": p.race.Sprite(),
		"E": map[string]interface{}{
			"a": "ccr",
			"s": 4,
			"w": w,
			"h": h,
		},
	})

	for _, e := range p.equipped {
		w, h := e.SpriteSize()
		for _, c := range e.Colors() {
			sprites = append(sprites, map[string]interface{}{
				"C": c,
				"S": e.Sprite(),
				"E": map[string]interface{}{
					"a": "ccr",
					"s": 4,
					"w": w,
					"h": h,
				},
			})
		}
	}

	p.SetHUD("cc", map[string]interface{}{
		"S": sprites,
		"N": p.name.Name(),
		"G": p.gender.Name(),
	})
}

func (p *Player) AdminCommand(addr string, command ...string) {
	p.mtx.Lock()
	if !p.admin {
		p.Kick("I'm sorry, Dave. I'm afraid you can't do that.")
		p.mtx.Unlock()
		return
	}

	log.Printf("[admin_cmd] %q:%q %#v", addr, p.login, command)
	p.mtx.Unlock()

	if len(command) == 0 {
		return
	}
	switch command[0] {
	case "grant admin":
		if len(command) != 2 {
			return
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		player, err := getPlayerKick(command[1], "You have been given admin. Please refresh the page and log back in.")
		if err != nil {
			return
		}
		player.admin = true
		savePlayer(player)
	case "revoke admin":
		if len(command) != 2 {
			return
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		player, err := getPlayerKick(command[1], "demoted by admin")
		if err != nil {
			return
		}
		if player.login == "BenLubar" {
			p.Kick("lolnope")
			return
		}
		player.admin = false
		savePlayer(player)
	case "kick":
		if len(command) < 2 || len(command) > 3 {
			return
		}
		message := "kicked by admin"
		if len(command) > 2 {
			message = command[2]
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		getPlayerKick(command[1], message)
	case "nuke tile":
		if len(command) != 1 {
			return
		}
		pos := p.Position()
		if pos == nil {
			return
		}
		for _, o := range pos.Objects() {
			if _, ok := o.(*Player); !ok {
				pos.Remove(o)
			}
		}
	case "unequip all":
		if len(command) != 1 {
			return
		}
		p.mtx.Lock()
		for _, e := range p.equipped {
			p.mtx.Unlock()
			e.Interact(p, "unequip")
			p.mtx.Lock()
		}
		p.mtx.Unlock()
	case "spawn item":
		if len(command) != 2 {
			return
		}
		item := world.Spawn(command[1])
		if item != nil {
			p.mtx.Lock()
			p.giveItem(item)
			p.mtx.Unlock()
		}
	case "spawn drop":
		if len(command) != 2 {
			return
		}
		obj := world.Spawn(command[1])
		if obj != nil {
			p.Position().Add(obj)
		}
	case "spawn equip":
		if len(command) != 2 {
			return
		}
		obj := world.Spawn(command[1])
		if e, ok := obj.(*Equip); ok {
			p.mtx.Lock()
			if old, ok := p.equipped[e.slot]; ok {
				old.wearer = nil
				p.giveItem(old)
			}
			e.wearer = p
			p.equipped[e.slot] = e
			p.mtx.Unlock()
			p.Position().Zone().Update(p.Position(), p)
		}
	case "clear inventory":
		if len(command) != 1 {
			return
		}
		for _, item := range p.Inventory() {
			p.RemoveItem(item)
		}
	case "tp":
		if len(command) != 3 {
			return
		}
		x, err := strconv.ParseUint(command[1], 0, 8)
		if err != nil {
			return
		}
		y, err := strconv.ParseUint(command[2], 0, 8)
		if err != nil {
			return
		}
		pos := p.Position()
		if pos == nil {
			return
		}
		if pos.Remove(p) {
			pos.Zone().Tile(uint8(x), uint8(y)).Add(p)
		}
	case "tpt", "tpp":
		if len(command) != 2 {
			return
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		if o := onlinePlayers[command[1]]; o != nil {
			if pos := o.Position(); pos != nil {
				if old := p.Position(); old != nil {
					if old.Remove(p) {
						pos.Add(p)
					}
				}
			}
		}
	case "summon":
		if len(command) != 2 {
			return
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		if o := onlinePlayers[command[1]]; o != nil {
			if pos := p.Position(); pos != nil {
				if old := o.Position(); old != nil {
					if old.Remove(o) {
						pos.Add(o)
						o.SendMessage("you have been summoned!")
					}
				}
			}
		}
	case "online":
		if len(command) != 1 {
			return
		}
		loginLock.Lock()
		defer loginLock.Unlock()
		for _, o := range onlinePlayers {
			if pos := o.Position(); pos != nil {
				x, y := pos.Position()
				z := pos.Zone()
				p.SendMessage(fmt.Sprintf("L:%q I:%q N:%q X:%d:%d Y:%d:%d", o.login, o.lastAddr, o.Name(), z.X, x, z.Y, y))
			} else {
				p.SendMessage(fmt.Sprintf("L:%q I:%q SPAWNING", o.login, o.lastAddr))
			}
		}
	case "butcher":
		if len(command) != 2 {
			return
		}
		r, err := strconv.Atoi(command[1])
		if err != nil {
			return
		}
		pos := p.Position()
		if pos == nil {
			return
		}
		x, y := pos.Position()
		sx, sy := int(x)-r, int(y)-r
		ex, ey := int(x)+r, int(y)+r
		for x := sx; x <= ex; x++ {
			if x < 0 || x > 255 {
				continue
			}
			for y := sy; y <= ey; y++ {
				if y < 0 || y > 255 {
					continue
				}
				t := pos.Zone().Tile(uint8(x), uint8(y))
				for _, o := range t.Objects() {
					if c, ok := o.(world.Combat); ok {
						if _, ok = c.(*Player); !ok {
							c.Hurt(c.MaxHealth(), p)
						}
					}
				}
			}
		}
	case "first name":
		if len(command) != 2 {
			return
		}
		p.mtx.Lock()
		origT := p.name.FirstT
		orig := p.name.First
		for p.name.FirstT = 0; p.name.FirstT < nameSubtypeCount; p.name.FirstT++ {
			for p.name.First = 0; p.Hero.name.First < uint64(len(names[p.name.FirstT])); p.name.First++ {
				if names[p.name.FirstT][p.name.First] == command[1] {
					p.mtx.Unlock()
					if pos := p.Position(); pos != nil {
						pos.Zone().Update(pos, p)
					}
					return
				}
			}
		}
		p.name.FirstT = origT
		p.name.First = orig
		p.mtx.Unlock()
	case "nickname":
		if len(command) != 2 {
			return
		}
		p.mtx.Lock()
		p.name.Nickname = command[1]
		p.mtx.Unlock()
		if pos := p.Position(); pos != nil {
			pos.Zone().Update(pos, p)
		}
	case "last name 1":
		if len(command) != 2 {
			return
		}
		p.mtx.Lock()
		origT := p.name.Last1T
		orig := p.name.Last1
		for p.name.Last1T = 0; p.name.Last1T < nameSubtypeCount; p.name.Last1T++ {
			for p.name.Last1 = 0; p.Hero.name.Last1 < uint64(len(names[p.name.Last1T])); p.name.Last1++ {
				if names[p.name.Last1T][p.name.Last1] == command[1] {
					p.mtx.Unlock()
					if pos := p.Position(); pos != nil {
						pos.Zone().Update(pos, p)
					}
					return
				}
			}
		}
		p.name.Last1T = origT
		p.name.Last1 = orig
		p.mtx.Unlock()
	case "last name 2":
		if len(command) != 2 {
			return
		}
		p.mtx.Lock()
		origT := p.name.Last2T
		orig := p.name.Last2
		for p.name.Last2T = 0; p.name.Last2T < nameSubtypeCount; p.name.Last2T++ {
			for p.name.Last2 = 0; p.Hero.name.Last2 < uint64(len(names[p.name.Last2T])); p.name.Last2++ {
				if names[p.name.Last2T][p.name.Last2] == command[1] {
					p.mtx.Unlock()
					if pos := p.Position(); pos != nil {
						pos.Zone().Update(pos, p)
					}
					return
				}
			}
		}
		p.name.Last2T = origT
		p.name.Last2 = orig
		p.mtx.Unlock()
	case "last name 3":
		if len(command) != 2 {
			return
		}
		p.mtx.Lock()
		origT := p.name.Last3T
		orig := p.name.Last3
		for p.name.Last3T = 0; p.name.Last3T < nameSubtypeCount; p.name.Last3T++ {
			for p.name.Last3 = 0; p.Hero.name.Last3 < uint64(len(names[p.name.Last3T])); p.name.Last3++ {
				if names[p.name.Last3T][p.name.Last3] == command[1] {
					p.mtx.Unlock()
					if pos := p.Position(); pos != nil {
						pos.Zone().Update(pos, p)
					}
					return
				}
			}
		}
		p.name.Last3T = origT
		p.name.Last3 = orig
		p.mtx.Unlock()
	}
}

func (p *Player) notifyInventoryChanged() {
	inventory := make([]world.Visible, len(p.items))
	copy(inventory, p.items)
	for {
		select {
		case p.inventory <- inventory:
			return
		case <-p.inventory:
		}
	}
}

func (p *Player) SendMessage(message string) {
	p.SendMessageColor(message, "#ddd")
}

func (p *Player) SendMessageColor(message, color string) {
	messages := []Message{{
		Text:  message,
		Color: color,
	}}
	for {
		select {
		case p.messages <- messages:
			return
		default:
			select {
			case p.messages <- messages:
				return
			case other := <-p.messages:
				messages = append(other, messages...)
			}
		}
	}
}

func (p *Player) Chat(addr, message string) {
	if pos := p.Position(); pos != nil {
		log.Printf("[info_chat] %s:%q %q %q", addr, p.login, p.Name(), message)
		pos.Zone().Chat(p, message)
	}
}

func (p *Player) IsAdmin() bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.admin
}

func (p *Player) Impersonate(o world.Visible) {
	if o == nil || o == p {
		if pos := p.Position(); pos != nil {
			pos.Zone().Impersonate(p, nil)
		}

		p.mtx.Lock()
		if p.impersonating != nil {
			p.impersonating = nil
			p.mtx.Unlock()
			if pos := p.Position(); pos != nil {
				pos.Zone().Update(pos, p)
			}
		} else {
			p.mtx.Unlock()
		}

		return
	}

	if other, ok := o.(*Player); ok {
		h := other.Hero
		o = &h
	}

	// copy the object
	o = world.LoadConvert(world.SaveConvert(o)).(world.Visible)

	if pos := p.Position(); pos != nil {
		pos.Zone().Impersonate(p, o)
	}

	p.mtx.Lock()
	p.impersonating = o
	p.mtx.Unlock()
	if pos := p.Position(); pos != nil {
		pos.Zone().Update(pos, p)
	}
}

func (p *Player) impersonate() world.Visible {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.impersonating
}

func (p *Player) NotifyPosition(old, new *world.Tile) {
	if i := p.impersonate(); i != nil {
		i.NotifyPosition(old, new)
	}

	p.Hero.NotifyPosition(old, new)
}

func (p *Player) Think() {
	p.Hero.Think()

	if i := p.impersonate(); i != nil {
		i.Think()
	}
}

type HUD struct {
	Name string                 `json:"N"`
	Data map[string]interface{} `json:"D,omitempty"`
}

type Message struct {
	Text  string `json:"T"`
	Color string `json:"C"`
}
