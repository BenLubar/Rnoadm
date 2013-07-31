package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/BenLubar/Rnoadm/resource"
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
		if b, ok := resource.Resource[r.URL.Path[1:]]; ok {
			w.Write(b)
			return
		}
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
	background: #000;
	text-align: center;
}
</style>
</head>
<body>
<canvas></canvas>
<script>
var tileSize = 32;
var authkey = localStorage['rnoadm-auth'] || (localStorage['rnoadm-auth'] = generateAuthKey());
var canvas;
var images = {};
var images_recolor = {};
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
			canvas = document.querySelector('canvas');
			canvas.width = 32*tileSize;
			canvas.height = 16*tileSize;
			canvas = canvas.getContext('2d');
			canvas.font = '18px sans-serif';
		}

		canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);

		for (var i = 0; i < 16; i++) {
			for (var j = 32-1; j >= 0; j--) {
				for (var k in (msg.Paint[j][i] || [])) {
					var p = msg.Paint[j][i][k];
					if (p.R) {
						canvas.fillStyle = '#000';
						canvas.fillText(p.R, j*tileSize+tileSize/4, i*tileSize+tileSize*3/4+1);
						canvas.fillStyle = p.C;
						canvas.fillText(p.R, j*tileSize+tileSize/4, i*tileSize+tileSize*3/4);
					}
					if (p.I) {
						if (!images[p.I]) {
							images[p.I] = true;
							(function(img, p) {
								img.onload = function() {
									images[p.I] = img;
									images_recolor[p.I] = {};
									send({ForceRepaint:true});
								};
								img.src = p.I + '.png';
							})(new Image(), p);
						}
						if (images[p.I] === true)
							continue;
						if (!images_recolor[p.I][p.C]) {
							var buffer = document.createElement('canvas');
							buffer.width = images[p.I].width || 1;
							buffer.height = images[p.I].height || 1;
							images_recolor[p.I][p.C] = buffer;
							buffer = buffer.getContext('2d');
							buffer.fillStyle = p.C;
							buffer.fillRect(0, 0, 1, 1);
							var data = buffer.getImageData(0, 0, 1, 1);
							var r = data.data[0], g = data.data[1], b = data.data[2], a = data.data[3];
							buffer.clearRect(0, 0, 1, 1);
							buffer.drawImage(images[p.I], 0, 0);
							data = buffer.getImageData(0, 0, buffer.canvas.width, buffer.canvas.height);
							var fade = function(x, y) {
								if (x >= 128)
									return 255 - fade(255-x, 255-y);
								return x*y/127;
							};
							for (var l = 0; l < data.data.length; l += 4) {
								data.data[l+0] = fade(data.data[l+0], r);
								data.data[l+1] = fade(data.data[l+1], g);
								data.data[l+2] = fade(data.data[l+2], b);
								data.data[l+3] = data.data[l+3]*a/255;
							}
							buffer.putImageData(data, 0, 0);

						}
						canvas.drawImage(images_recolor[p.I][p.C], j*tileSize, i*tileSize+tileSize-Math.max(tileSize, images[p.I].height), Math.max(tileSize, images[p.I].width), Math.max(tileSize, images[p.I].height));
					}
				}
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
		e.preventDefault();
	}
};
document.onkeypress = function(e) {
	send({Key:{Code:e.charCode, Special:false}});
};
document.querySelector('canvas').onclick = document.querySelector('canvas').oncontextmenu = function(e) {
	send({Click:{X:Math.floor(e.offsetX/tileSize), Y:Math.floor(e.offsetY/tileSize)}});
	return false;
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
	Click *struct {
		X, Y int
	}
	ForceRepaint bool
}

type packetPaintCell struct {
	R, I, C string `json:",omitempty"`
}
type packetPaint struct {
	Paint [40][16][]packetPaintCell
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

	onlinePlayersLock.Lock()
	if OnlinePlayers[playerID] != nil {
		player = OnlinePlayers[playerID]
	} else {
		OnlinePlayers[playerID] = player
	}
	onlinePlayersLock.Unlock()
	defer func() {
		onlinePlayersLock.Lock()
		delete(OnlinePlayers, playerID)
		onlinePlayersLock.Unlock()
	}()
	if player.Joined.IsZero() {
		player.Joined = time.Now().UTC()
	}
	player.LastLogin = time.Now().UTC()

	var messageBuffer [4]string
	messages := make(chan string)
	player.messages = messages

	player.Repaint()
	player.Lock()
	zone := GrabZone(player.ZoneX, player.ZoneY)
	player.zone = zone
	tx, ty := player.TileX, player.TileY
	if zone.Tile(tx, ty) == nil {
		tx, ty, player.TileX, player.TileY = 127, 127, 127, 127
	}
	player.Unlock()
	zone.Lock()
	zone.Tile(tx, ty).Add(player)
	zone.Unlock()
	zone.Repaint()
	defer func() {
		player.Lock()
		zone := player.zone
		tx, ty := player.TileX, player.TileY
		player.Unlock()
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

	const w, h = 20, 12
	for {
		select {
		case p, ok := <-packets:
			if !ok {
				return
			}

			if p.ForceRepaint {
				player.Repaint()
			}

			if p.Key != nil {
				player.Lock()
				player.schedule = nil
				player.Unlock()
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
			if p.Click != nil {
				player.Lock()
				player.schedule = nil
				player.Unlock()
				if player.hud != nil && player.hud.Click(p.Click.X, p.Click.Y) {
					break
				}
				if p.Click.X >= 0 && p.Click.X < w && p.Click.Y >= 0 && p.Click.Y < h {
					player.hud = &ClickHUD{
						X:      p.Click.X,
						Y:      p.Click.Y,
						W:      w,
						H:      h,
						Player: player,
					}
					player.Repaint()
				}
			}

		case <-player.repaint:
			var paint packetPaint
			setcell := func(x, y int, ch string, sprite string, color Color) {
				if x >= 0 && x < len(paint.Paint) {
					if y >= 0 && y < len(paint.Paint[x]) {
						paint.Paint[x][y] = append(paint.Paint[x][y], packetPaintCell{
							R: ch,
							I: sprite,
							C: string(color),
						})
					}
				}
			}

			camX := int(player.TileX)
			camY := int(player.TileY)

			player.Lock()
			z := player.zone
			player.Unlock()
			z.Lock()

			for x := 0; x < w; x++ {
				xCoord := x - w/2 + camX
				x8 := uint8(xCoord)
				for y := 0; y < h; y++ {
					yCoord := y - h/2 + camY
					if xCoord < 0 || xCoord > 255 || yCoord < 0 || yCoord > 255 {
						setcell(x, y, "?", "", "#111")
						continue
					}
					y8 := uint8(yCoord)
					if tile := z.Tile(x8, y8); tile != nil {
						tile.Paint(z, x, y, setcell)
					} else {
						setcell(x, y, "?", "", "#111")
					}
				}
			}

			if player.hud == nil {
				player.hud = ZoneEntryHUD(z.Name())
			}
			z.Unlock()

			for i, message := range messageBuffer {
				setcell(0, h+i, message, "", "#fff")
			}

			player.Lock()
			setcell(w+1, 0, "WEARING", "", "#aaa")
			if player.Head != nil {
				player.Head.Paint(w+1, 1, setcell)
			}
			if player.Top != nil {
				player.Top.Paint(w+2, 1, setcell)
			}
			if player.Legs != nil {
				player.Legs.Paint(w+3, 1, setcell)
			}
			if player.Feet != nil {
				player.Feet.Paint(w+4, 1, setcell)
			}

			setcell(w+6, 0, "TOOLBELT", "", "#aaa")
			if player.Toolbelt.Pickaxe != nil {
				player.Toolbelt.Pickaxe.Paint(w+6, 1, setcell)
			}
			if player.Toolbelt.Hatchet != nil {
				player.Toolbelt.Hatchet.Paint(w+7, 1, setcell)
			}

			setcell(w+1, 2, "INVENTORY", "", "#aaa")
			for i, o := range player.Backpack {
				o.Paint(i%10+w+1, i/10+3, setcell)
			}
			player.Unlock()

			player.hud.Paint(setcell)
			websocket.JSON.Send(conn, &paint)

		case message := <-messages:
			copy(messageBuffer[:], messageBuffer[1:])
			messageBuffer[len(messageBuffer)-1] = message
		}
	}
}
