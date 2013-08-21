package hero

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	Hero
	characterCreation bool
	ancestry          PlayerAncestry
	zoneX, zoneY      int64
	tileX, tileY      uint8

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
}

var _ world.NoSaveObject = (*Player)(nil)
var _ world.AdminLike = (*Player)(nil)
var _ world.SendMessageLike = (*Player)(nil)

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

	return 0, map[string]interface{}{
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
		"regt":     p.registered.Format(time.RFC3339Nano),
		"lastt":    p.lastLogin.Format(time.RFC3339Nano),
	}, []world.ObjectLike{&p.Hero, &p.ancestry}
}

func (p *Player) Load(version uint, data interface{}, attached []world.ObjectLike) {
	switch version {
	case 0:
		dataMap := data.(map[string]interface{})
		p.Hero = *attached[0].(*Hero)
		p.mtx.Lock()
		defer p.mtx.Unlock()
		p.ancestry = *attached[1].(*PlayerAncestry)
		p.characterCreation = dataMap["creation"].(bool)
		p.zoneX, p.zoneY = dataMap["zx"].(int64), dataMap["zy"].(int64)
		p.tileX, p.tileY = dataMap["tx"].(uint8), dataMap["ty"].(uint8)
		p.login = dataMap["u"].(string)
		p.password = dataMap["p"].([]byte)
		p.admin, _ = dataMap["admin"].(bool)
		p.firstAddr = dataMap["rega"].(string)
		p.lastAddr = dataMap["lasta"].(string)
		var err error
		p.registered, err = time.Parse(time.RFC3339Nano, dataMap["regt"].(string))
		if err != nil {
			panic(err)
		}
		p.lastLogin, err = time.Parse(time.RFC3339Nano, dataMap["lastt"].(string))
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
	for _, e := range p.Hero.equipped {
		if e.AdminOnly() && !p.admin {
			delete(p.Hero.equipped, e.slot)
		} else {
			e.wearer = &p.Hero
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

	objects := []world.ObjectLike{&a.Object}
	for _, h := range a.ancestors {
		objects = append(objects, h)
	}

	return 0, uint(0), objects
}

func (a *PlayerAncestry) Load(version uint, data interface{}, attached []world.ObjectLike) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	switch version {
	case 0:
		a.Object = *attached[0].(*world.Object)
		attached = attached[1:]
		a.ancestors = make([]*Hero, len(attached))
		for i, h := range attached {
			a.ancestors[i] = h.(*Hero)
		}
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
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
			e.wearer = &p.Hero
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

func (p *Player) Chat(message string) {
	if pos := p.Position(); pos != nil {
		pos.Zone().Chat(p, message)
	}
}

func (p *Player) IsAdmin() bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	return p.admin
}

type HUD struct {
	Name string                 `json:"N"`
	Data map[string]interface{} `json:"D,omitempty"`
}

type Message struct {
	Text  string `json:"T"`
	Color string `json:"C"`
}
