package main

import (
	"code.google.com/p/go.net/websocket"
	"hash/crc64"
	"log"
	"net/http"
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
</style>
</head>
<body>
<script>
var front, back;
var ws = new WebSocket('ws://' + location.host + '/ws');
var wsonopen = ws.onopen = function() {
	console.log("open");
};
var wsonmessage = ws.onmessage = function(e) {
	var msg = JSON.parse(e.data);
	if (msg.Paint) {
		back = document.createElement('table');
		back.style.display = 'none';
		document.body.appendChild(back);

		for (var i = 0; i < 24; i++) {
			var row = document.createElement('tr');
			for (var j = 0; j < 80; j++) {
				var cell = document.createElement('td');
				cell.textContent = msg.Paint[j][i].Char;
				cell.style.color = msg.Paint[j][i].Color;
				row.appendChild(cell);
			}
			back.appendChild(row);
		}

		if (front) {
			front.parentNode.removeChild(front);
		}
		back.style.display = 'inline-block';
		front = back;
		back = null;
	}
};
var wsonclose = ws.onclose = function() {
	console.log("close");
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
	Key *struct {
		Code int
	}
}
type packetPaint struct {
	Paint [80][24]struct {
		Char, Color string
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

	playerID := crc64.Checksum([]byte(addr), crc64.MakeTable(crc64.ISO))
	player, err := LoadPlayer(playerID)
	if err != nil {
		player = &Player{
			ID:      playerID,
			TileX:   127,
			TileY:   127,
			repaint: make(chan struct{}, 1),
		}
	}
	player.Name_ = &Name{Name: addr}
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
					log.Printf("[%s] %d", addr, p.Key.Code)
				}
			}
		case <-player.repaint:
			var paint packetPaint
			for i := range paint.Paint {
				for j := range paint.Paint[i] {
					paint.Paint[i][j].Char = " "
					paint.Paint[i][j].Color = "#fff"
				}
			}
			setcell := func(x, y int, ch rune, color Color) {
				if x >= 0 && x < len(paint.Paint) {
					if y >= 0 && y < len(paint.Paint[x]) {
						paint.Paint[x][y].Char = string(ch)
						paint.Paint[x][y].Color = string(color)
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
						r, fg := tile.Paint()
						setcell(x, y, r, fg)
					}
				}
			}

			z.Unlock()
			ReleaseZone(z)

			if player.hud != nil {
				player.hud.Paint(setcell)
			}
			websocket.JSON.Send(conn, &paint)
		}
	}
}
