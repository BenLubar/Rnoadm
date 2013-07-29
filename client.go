package main

import (
	"code.google.com/p/go.net/websocket"
	"hash/crc64"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var OnlinePlayers = make(map[uint64]*Player)
var onlinePlayersLock sync.Mutex

func init() {
	http.HandleFunc("/", httpHandler)
	http.Handle("/ws", websocket.Handler(websocketHandler))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Rnoadm</title>
<script>document.title = 'R' + 'ando'.split('').sort(function(a, b) {return Math.random()-.5}).join('') + 'm'</script>
<style>
html {
	background:  #000;
	text-align:  center;
}
table {
	background:  #000;
	font-family: monospace;
	margin:      0 auto;
}
td {
	width:       1em;
	height:      1em;
}
</style>
</head>
<body>
<script>
var authkey = localStorage['rnoadm-auth'] || (localStorage['rnoadm-auth'] = generateAuthKey());
var canvas;
function generateAuthKey() {
	return '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'.split('').sort(function(a, b) {return Math.random()-.5}).join('');
}
var ws = new WebSocket('ws://' + location.host + '/ws');
var wsonopen = ws.onopen = function() {
	send({Auth:{Key:authkey}});
};
var wsonmessage = ws.onmessage = function(e) {
	var msg = JSON.parse(e.data);
	if (msg.Paint) {
		if (!canvas) {
			canvas = document.createElement('table');
			for (var i = 0; i < 24; i++) {
				var row = document.createElement('tr');
				for (var j = 0; j < 80; j++) {
					row.appendChild(document.createElement('td'));
				}
				canvas.appendChild(row);
			}
			document.body.appendChild(canvas);
		}

		for (var i = 0; i < 24; i++) {
			var row = canvas.children[i];
			for (var j = 0; j < 80; j++) {
				var cell = row.children[j];
				cell.textContent = msg.Paint[j][i].R;
				cell.style.color = msg.Paint[j][i].C;
			}
		}
	}
};
var wsonclose = ws.onclose = function() {
	setTimeout(function() {
		ws = new WebSocket('ws://' + location.host + '/ws');
		ws.onopen = wsonopen;
		ws.onmessage = wsonmessage;
		ws.onclose = wsonclose;
	}, 1000);
};
function send(packet) {
	ws.send(JSON.stringify(packet));
}
document.onkeydown = function(e) {
	send({Key:{Code:e.keyCode, Special:true}});
	if (e.keyCode == 8) {
		e.preventDefault()
	}
};
document.onkeypress = function(e) {
	send({Key:{Code:e.charCode, Special:false}});
};
</script>
</body>
</html>`))
}

type packetIn struct {
	Auth *struct {
		Key string
	}
	Key *struct {
		Code    int
		Special bool
	}
}
type packetPaint struct {
	Paint [80][24]struct {
		R, C string
	}
}

func websocketHandler(conn *websocket.Conn) {
	defer conn.Close()

	addr := conn.Request().RemoteAddr
	for i := range addr {
		if addr[i] == ':' {
			addr = addr[:i]
			break
		}
	}

	var playerID uint64
	{
		var p packetIn
		err := websocket.JSON.Receive(conn, &p)
		if err != nil {
			log.Printf("[%s] %v", addr, err)
			return
		}
		if p.Auth == nil {
			log.Printf("[%s] noauth", addr)
			return
		}
		playerID = crc64.Checksum([]byte(p.Auth.Key), crc64.MakeTable(crc64.ISO))
	}
	player, err := LoadPlayer(playerID)
	if err != nil {
		player = &Player{
			ID:    playerID,
			TileX: 127,
			TileY: 127,
			Hero: Hero{
				Name_: GenerateName(rand.New(rand.NewSource(int64(playerID))), NameHero),
			},
			repaint: make(chan struct{}, 1),
			Joined:  time.Now().UTC(),
		}
	}
	if player.Joined.IsZero() {
		player.Joined = time.Now().UTC()
	}
	player.LastLogin = time.Now().UTC()

	onlinePlayersLock.Lock()
	if OnlinePlayers[playerID] != nil {
		onlinePlayersLock.Unlock()
		log.Printf("[%s:%d] multilog", addr, playerID)
		return
	}
	OnlinePlayers[playerID] = player
	onlinePlayersLock.Unlock()
	defer func() {
		onlinePlayersLock.Lock()
		delete(OnlinePlayers, playerID)
		onlinePlayersLock.Unlock()
	}()

	player.Repaint()
	player.lock.Lock()
	zone := GrabZone(player.ZoneX, player.ZoneY)
	player.zone = zone
	tx, ty := player.TileX, player.TileY
	player.lock.Unlock()
	zone.Lock()
	zone.Tile(tx, ty).Add(player)
	zone.Unlock()
	zone.Repaint()
	defer func() {
		player.lock.Lock()
		zone := player.zone
		tx, ty := player.TileX, player.TileY
		player.lock.Unlock()
		zone.Lock()
		zone.Tile(tx, ty).Remove(player)
		zone.Unlock()
		zone.Repaint()
		ReleaseZone(zone)
		player.Save()
	}()

	packets := make(chan packetIn)
	go func() {
		for {
			var p packetIn
			err := websocket.JSON.Receive(conn, &p)
			if err != nil {
				log.Printf("[%s:%d] %v", addr, playerID, err)
				close(packets)
				return
			}
			select {
			case packets <- p:
			case <-time.After(time.Second):
				log.Printf("[%s:%d] dropped a packet (server)", addr, playerID)
			}
		}
	}()

	for {
		select {
		case p, ok := <-packets:
			if !ok {
				return
			}

			if p.Key != nil {
				if player.hud != nil && player.hud.Key(p.Key.Code, p.Key.Special) {
					break
				}
				if !p.Key.Special {
					break
				}
				switch p.Key.Code {
				case 38, 'W':
					player.Move(0, -1)
				case 37, 'A':
					player.Move(-1, 0)
				case 40, 'S':
					player.Move(0, 1)
				case 39, 'D':
					player.Move(1, 0)
				case 'E':
					player.hud = &InteractHUD{Player: player}
					player.Repaint()
				case 'I':
					player.hud = &InventoryHUD{Player: player}
					player.Repaint()
				case 192: // `
					if player.Admin {
						player.hud = &AdminHUD{Player: player}
						player.Repaint()
					}
				default:
					//log.Printf("[%s:%d] %d", addr, playerID, p.Key.Code)
				}
			}
		case <-player.repaint:
			var paint packetPaint
			setcell := func(x, y int, ch rune, color Color) {
				if x >= 0 && x < len(paint.Paint) {
					if y >= 0 && y < len(paint.Paint[x]) {
						paint.Paint[x][y].R = string(ch)
						paint.Paint[x][y].C = string(color)
					}
				}
			}

			w := len(paint.Paint)
			h := len(paint.Paint[0])
			camX := int(player.TileX)
			camY := int(player.TileY)

			player.lock.Lock()
			z := player.zone
			player.lock.Unlock()
			z.Lock()

			for x := 0; x < w; x++ {
				xCoord := x - w/2 + camX
				x8 := uint8(xCoord)
				for y := 0; y < h; y++ {
					yCoord := y - h/2 + camY
					if xCoord < 0 || xCoord > 255 || yCoord < 0 || yCoord > 255 {
						setcell(x, y, '?', "#111")
						continue
					}
					y8 := uint8(yCoord)
					if tile := z.Tile(x8, y8); tile != nil {
						r, fg := tile.Paint(z)
						setcell(x, y, r, fg)
					} else {
						setcell(x, y, '?', "#111")
					}
				}
			}

			if player.hud == nil {
				player.hud = ZoneEntryHUD(z.Name())
			}
			z.Unlock()

			player.hud.Paint(setcell)
			websocket.JSON.Send(conn, &paint)
		}
	}
}
