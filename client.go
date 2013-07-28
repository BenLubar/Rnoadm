package main

import (
	"code.google.com/p/go.net/websocket"
	"hash/crc64"
	"log"
	"math/rand"
	"net/http"
	"time"
)

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
	send({Key:{Code:e.keyCode}});
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
		Code int
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
	player.Repaint()
	zone := GrabZone(player.ZoneX, player.ZoneY)
	zone.Lock()
	zone.Tile(player.TileX, player.TileY).Add(player)
	player.hud = ZoneEntryHUD(zone.Name())
	zone.Unlock()
	defer func() {
		zone := GrabZone(player.ZoneX, player.ZoneY)
		zone.Lock()
		zone.Tile(player.TileX, player.TileY).Remove(player)
		zone.Unlock()
		ReleaseZone(zone)
		ReleaseZone(zone)
		player.Save()
	}()

	packets := make(chan packetIn)
	go func() {
		for {
			var p packetIn
			err := websocket.JSON.Receive(conn, &p)
			if err != nil {
				log.Printf("[%s] %v", addr, err)
				close(packets)
				return
			}
			packets <- p
		}
	}()

	for {
		select {
		case p, ok := <-packets:
			if !ok {
				return
			}

			if p.Key != nil {
				if player.hud != nil && player.hud.Key(p.Key.Code) {
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
				default:
					//log.Printf("[%s] %d", addr, p.Key.Code)
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

			z := GrabZone(player.ZoneX, player.ZoneY)
			z.Lock()

			for x := 0; x < w; x++ {
				xCoord := x - w/2 + camX
				if xCoord < 0 || xCoord > 255 {
					continue
				}
				x8 := uint8(xCoord)
				for y := 0; y < h; y++ {
					yCoord := y - h/2 + camY
					if yCoord < 0 || yCoord > 255 {
						continue
					}
					y8 := uint8(yCoord)
					if tile := z.Tile(x8, y8); tile != nil {
						r, fg := tile.Paint(z)
						setcell(x, y, r, fg)
					}
				}
			}

			if player.hud == nil {
				player.hud = ZoneEntryHUD(z.Name())
			}
			z.Unlock()
			ReleaseZone(z)

			player.hud.Paint(setcell)
			websocket.JSON.Send(conn, &paint)
		}
	}
}
