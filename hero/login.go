package hero

import (
	"code.google.com/p/go.crypto/bcrypt"
	"compress/gzip"
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const CryptoCost = bcrypt.DefaultCost

var onlinePlayers = make(map[string]*Player)
var loginAttempts = make(map[string]int)
var loginLock sync.Mutex

func init() {
	go func() {
		for {
			loginLock.Lock()
			for a, n := range loginAttempts {
				if n <= 1 || n == 6 {
					delete(loginAttempts, a)
				} else {
					loginAttempts[a]--
				}
			}
			loginLock.Unlock()
			time.Sleep(time.Minute)
		}
	}()
}

type LoginPacket struct {
	Login string `json:"U"`
	Pass  string `json:"P"`
}

// Login returns a non-nil player and an empty error string OR a nil player
// and an error to show to the user.
func Login(addr string, packet *LoginPacket) (*Player, string) {
	login := strings.TrimSpace(packet.Login)
	if login == "" {
		return nil, "A username is required."
	}
	pass := []byte(packet.Pass)
	if len(pass) <= 2 {
		return nil, "A password is required."
	}
	filename := loginToFilename(login)

	loginLock.Lock()
	defer loginLock.Unlock()

	loginAttempts[addr]++
	if loginAttempts[addr] == 5 {
		loginAttempts[addr] += 60
	}
	if loginAttempts[addr] > 5 {
		return nil, fmt.Sprintf("Too many login attempts. Come back in %d minutes.", loginAttempts[addr]-5)
	}

	f, err := os.Open(filename)
	if err != nil {
		hashedPass, err := bcrypt.GenerateFromPassword(pass, CryptoCost)
		if err != nil {
			panic(err)
		}
		p := &Player{
			Hero:              *GenerateHero(rand.New(rand.NewSource(rand.Int63()))),
			characterCreation: true,
			tileX:             127,
			tileY:             127,
			login:             login,
			password:          hashedPass,
			firstAddr:         addr,
			registered:        time.Now().UTC(),
			lastAddr:          addr,
			lastLogin:         time.Now().UTC(),
		}
		world.InitObject(p)
		for _, e := range p.Hero.equipped {
			e.wearer = p
		}
		savePlayer(p)
		onlinePlayers[login] = p
		return p, ""
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	var data interface{}
	err = gob.NewDecoder(g).Decode(&data)
	if err != nil {
		panic(err)
	}

	p := world.LoadConvert(data).(*Player)

	if bcrypt.CompareHashAndPassword(p.password, pass) != nil {
		return nil, "Username or password incorrect."
	}

	loginAttempts[addr]--

	if other, ok := onlinePlayers[p.login]; ok {
		savePlayer(other)
		other.Kick("Logged in from a different location.")
		_, err = f.Seek(0, 0)
		if err != nil {
			panic(err)
		}

		g, err = gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer g.Close()

		data = nil
		err = gob.NewDecoder(g).Decode(&data)
		if err != nil {
			panic(err)
		}
		p = world.LoadConvert(data).(*Player)
	}

	cost, err := bcrypt.Cost(p.password)
	if err != nil {
		panic(err)
	}
	if cost != CryptoCost {
		p.password, err = bcrypt.GenerateFromPassword(pass, CryptoCost)
		if err != nil {
			panic(err)
		}
		savePlayer(p)
	}

	onlinePlayers[p.login] = p

	return p, ""
}

func loginToFilename(login string) string {
	canonicalLogin := base32.StdEncoding.EncodeToString([]byte(strings.ToLower(login)))
	for i := range canonicalLogin {
		if canonicalLogin[i] == '=' {
			canonicalLogin = canonicalLogin[:i]
			break
		}
	}
	return filepath.Join("rnoadm-AA", "player"+canonicalLogin+".gz")
}

func SavePlayer(p *Player) {
	loginLock.Lock()
	defer loginLock.Unlock()

	savePlayer(p)
}

func savePlayer(p *Player) {
	data := world.SaveConvert(p)

	f, err := os.Create(loginToFilename(p.login))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	err = gob.NewEncoder(g).Encode(&data)
	if err != nil {
		panic(err)
	}
}

func PlayerDisconnected(p *Player) {
	if pos := p.Position(); pos != nil {
		pos.Zone().Impersonate(p, nil)
	}
	loginLock.Lock()
	defer loginLock.Unlock()

	if onlinePlayers[p.login] == p {
		savePlayer(p)
		delete(onlinePlayers, p.login)
	}
}

func SaveAllPlayers() {
	loginLock.Lock()
	defer loginLock.Unlock()

	for _, p := range onlinePlayers {
		savePlayer(p)
	}
}

func getPlayerKick(login, message string) (*Player, error) {
	filename := loginToFilename(login)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	g, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	var data interface{}
	err = gob.NewDecoder(g).Decode(&data)
	if err != nil {
		panic(err)
	}

	p := world.LoadConvert(data).(*Player)

	if o, ok := onlinePlayers[p.login]; ok {
		p = o
		savePlayer(p)
		p.Kick(message)
		delete(onlinePlayers, p.login)
	}

	return p, nil
}
