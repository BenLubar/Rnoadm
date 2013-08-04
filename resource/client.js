const tileSize = 32;
const w = 32, h = 16;
var gameState;
var images = {};
var images_recolor = {};
var clientHash;
var huds = {};
var canvas = document.querySelector('canvas');
canvas.width = w*tileSize;
canvas.height = h*tileSize;
canvas = canvas.getContext('2d');
canvas.font = '11px sans-serif';

var frame = 0;
setInterval(function() {
	frame++;
	repaint();
}, 500);

var inRepaint = false;
function repaint() {
	if (inRepaint)
		return;

	requestAnimationFrame(function() {
		inRepaint = false;

		canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);

		var draw = function(x, y, p) {
			if (p.Text) {
				var tx = x*tileSize + tileSize/8;
				var ty = y*tileSize + tileSize*7/8;
				canvas.fillStyle = '#000';
				canvas.fillText(p.Text, tx, ty + 1);
				canvas.fillStyle = p.Color;
				canvas.fillText(p.Text, tx, ty);
			}
			if (p.Sprite) {
				if (!images[p.Sprite]) {
					images[p.Sprite] = true;
					var img = new Image();
					img.onload = function() {
						images[p.Sprite] = img;
						images_recolor[p.Sprite] = {};
						repaint();
					};
					img.src = p.Sprite + '.png';
				}
				if (images[p.Sprite] === true)
					return;
				if (!images_recolor[p.Sprite][p.Color]) {
					var buffer = document.createElement('canvas');
					buffer.width = images[p.Sprite].width || 1;
					buffer.height = images[p.Sprite].height || 1;
					images_recolor[p.Sprite][p.Color] = buffer;
					buffer = buffer.getContext('2d');
					buffer.fillStyle = p.Color;
					buffer.fillRect(0, 0, 1, 1);
					var data = buffer.getImageData(0, 0, 1, 1);
					var r = data.data[0], g = data.data[1], b = data.data[2], a = data.data[3];
					buffer.clearRect(0, 0, 1, 1);
					buffer.drawImage(images[p.Sprite], 0, 0);
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
				canvas.drawImage(images_recolor[p.Sprite][p.Color],
					(p.X || 0) * tileSize,
					(p.Y || 0) * (p.SpriteHeight || tileSize),
					tileSize,
					p.Height || tileSize,
					x*tileSize,
					y*tileSize + tileSize - (p.Height || tileSize),
					tileSize,
					p.Height || tileSize);
			}
		};

		if (gameState.hud)
			gameState.hud(draw);
	});
}

var ws = new WebSocket('ws://' + location.host + '/ws');
var wsonopen = ws.onopen = function() {
	gameState = {
		hud: loginHud
	};
	repaint();
};

var wsonmessage = ws.onmessage = function(e) {
	var msg = JSON.parse(e.data);
	if (msg['ClientHash']) {
		if (clientHash) {
			if (clientHash != msg['ClientHash']) {
				location.reload(true);
			}
		} else {
			clientHash = msg['ClientHash'];
		}
	}
	if (msg['SetHUD']) {
		gameState.hud = huds[msg['SetHUD']['Name']](msg['SetHUD']['Data']);
		repaint();
	}
	if (msg['Kick']) {
		ws.onclose = wsonclose = function() {};
		gameState = {};
		repaint();
		alert('Kicked: ' + msg['Kick']);
	}
	console.log(msg);
};
var wsonclose = ws.onclose = function() {
	setTimeout(function() {
		ws = new WebSocket('ws://' + location.host + '/ws');
		ws.onopen = wsonopen;
		ws.onmessage = wsonmessage;
		ws.onclose = wsonclose;
	}, 1000);

	gameState.hud = lostConnectionHud;
	repaint();
};
function send(packet) {
	ws.send(JSON.stringify(packet));
}
document.onkeydown = function(e) {
	if (e.altKey || e.ctrlKey)
		return;
	if (e.keyCode < 20)
		e.preventDefault();
	if (gameState.hud === loginHud) {
		switch (e.keyCode) {
		case 8: // backspace
			switch (loginHudFocus) {
			case 0:
				if (loginHudUsername)
					loginHudUsername = loginHudUsername.substring(0, loginHudUsername.length - 1);
				break;
			case 1:
				if (loginHudPassword)
					loginHudPassword = loginHudPassword.substring(0, loginHudPassword.length - 1);
				break;
			}
			repaint();
			break;
		case 9: // tab
			loginHudFocus = (loginHudFocus + 1) % 2;
			repaint();
			break;
		case 13: // enter
			loginHudSubmit();
			break;
		}
	}
};
document.onkeypress = function(e) {
	if (gameState.hud === loginHud) {
		switch (loginHudFocus) {
		case 0:
			loginHudUsername += String.fromCharCode(e.charCode);
			localStorage['login'] = loginHudUsername;
			break;
		case 1:
			loginHudPassword += String.fromCharCode(e.charCode);
			break;
		}
		repaint();
	}
};
document.querySelector('canvas').onclick = document.querySelector('canvas').oncontextmenu = function(e) {
	var x = Math.floor(e.offsetX * 4 / tileSize)/4 - w/2;
	var y = Math.floor(e.offsetY * 4 / tileSize)/4 - h/2;

	if (gameState.hud && gameState.hud.click) {
		if (gameState.hud.click(x, y) === false)
			return false;
	}
	return false;
};
/*var mouseX = -1, mouseY = -1;
var mouseTimeout;
document.querySelector('canvas').onmouseout = function() {
	if (mouseX == -1 && mouseY == -1) {
		return;
	}
	mouseX = -1;
	mouseY = -1;
	if (mouseTimeout) {
		clearTimeout(mouseTimeout);
		mouseTimeout = null;
	} else {
		send({MouseMove:{X:-1, Y:-1}});
	}
};
document.querySelector('canvas').onmousemove = function(e) {
	var x = Math.floor(e.offsetX/tileSize);
	var y = Math.floor(e.offsetY/tileSize);
	if (mouseX == x && mouseY == y) {
		return;
	}
	if (mouseTimeout) {
		clearTimeout(mouseTimeout);
	} else {
		send({MouseMove:{X:-1, Y:-1}});
	}
	mouseX = x;
	mouseY = y;
	mouseTimeout = setTimeout(function() {
		send({MouseMove:{X:x, Y:y}});
		mouseTimeout = null;
	}, 500);
};*/

var loginHudUsername = localStorage['login'] || '';
var loginHudPassword = '';
var loginHudFocus = 0;
var loginHud = function(draw) {
	for (var x = w/2 - 4; x < w/2 + 4; x++) {
		for (var y = h/2 - 4; y < h/2 + 1; y++) {
			draw(x, y, {
				Sprite: 'ui_r1',
				Color:  y == h/2 - 4 ? '#444' : '#222',
				X:      y == h/2 - 4 ? x == w/2 - 4 ? 3 : x == w/2 + 3 ? 4 : 0 : 0
			});
		}
		draw(x + 0.5, h/2 + 0.5, {
			Sprite: 'ui_r1',
			Color:  '#444',
			Y:      1
		});
		draw(x, h/2 + 0.5, {
			Sprite: 'ui_r1',
			Color:  '#444',
			Y:      1
		});
	}
	for (var x = w/2 - 3.75; x < w/2 + 3.75; x += 0.5) {
		draw(x, h/2 - 2.75, {
			Sprite: 'ui_r1',
			Color:  '#111',
			Y:      1
		});
		draw(x, h/2 - 1.25, {
			Sprite: 'ui_r1',
			Color:  '#111',
			Y:      1
		});
	}
	draw(w/2 - 1.625, h/2 - 4, {
		Sprite: 'ui_title',
		Color:  '#888',
		X:      0
	});
	var permutations = [1, 2, 3, 4].sort(function(a, b) { return Math.random() * 2 - 1; });
	for (var i in permutations) {
		draw(w/2 - 1 + i/2, h/2 - 4, {
			Sprite: 'ui_title',
			Color:  '#888',
			X:      permutations[i]
		});
	}
	draw(w/2 + 1, h/2 - 4, {
		Sprite: 'ui_title',
		Color:  '#888',
		X:      5
	});
	draw(w/2 - 3.75, h/2 - 3.25, {
		Text:  'Login',
		Color: '#aaa'
	});
	draw(w/2 - 3.75, h/2 - 2.75, {
		Text:  loginHudUsername + (loginHudFocus === 0 ? '_' : ''),
		Color: '#fff'
	});
	draw(w/2 - 3.75, h/2 - 1.75, {
		Text:  'Password',
		Color: '#aaa'
	});
	draw(w/2 - 3.75, h/2 - 1.25, {
		Text:  loginHudPassword.replace(/./g, '*') + (loginHudFocus === 1 ? '_' : ''),
		Color: '#fff'
	});
	draw(w/2 - 3, h/2 - 0.25, {
		Text:  'Log in or register',
		Color: '#aaa'
	});
};

loginHud.click = function(x, y) {
	if (x >= -4 && x < 4) {
		if (y >= -2.75 && y <= -1.75) {
			loginHudFocus = 0; // login
			repaint();
			return false;
		}
		if (y >= -1.25 && y <= -0.25) {
			loginHudFocus = 1; // password
			repaint();
			return false;
		}
		if (y >= 0.00 && y <= 0.75) {
			loginHudSubmit();
			return false;
		}
	}
};

var loginHudSubmit = function() {
	if (!loginHudUsername) {
		loginHudFocus = 0;
		repaint();
		return;
	}
	if (!loginHudPassword) {
		loginHudFocus = 1;
		repaint();
		return;
	}
	send({'Auth': {'Login': loginHudUsername, 'Password': loginHudPassword}});
	loginHudPassword = '';
	loginHudFocus = 0;
}

huds['character_creation'] = function(data) {
	var f = function(draw) {
		for (var x = w/2 - 4; x < w/2 + 4; x++) {
			for (var y = h/2 - 4; y < h/2 + 2; y++) {
				draw(x, y, {
					Sprite: 'ui_fill',
					Color:  '#777'
				});
			}
		}
		draw(w/2 - 2, h/2 - 2, {
			Text:  'Race:',
			Color: '#fff'
		});
		draw(w/2, h/2 - 2, {
			Text:  data['race'],
			Color: '#fff'
		});
		draw(w/2 - 2, h/2 - 1, {
			Text:  'Gender:',
			Color: '#fff'
		});
		draw(w/2, h/2 - 1, {
			Text:  data['gender'],
			Color: '#fff'
		});
		draw(w/2 - 2, h/2, {
			Text:  'Skin:',
			Color: '#fff'
		});
		var rotate = [0, 6, 3, 9];
		draw(w/2, h/2 + 0.25, {
			Sprite: 'body_' + data['race'],
			Color:  data['skin'],
			X:      rotate[frame % 4]
		});
		draw(w/2, h/2 + 0.25, {
			Sprite: 'censor_' + data['race'],
			Color:  data['skin'],
			X:      rotate[frame % 4]
		});
		if (data['gender'] == 'female') {
			draw(w/2, h/2 + 0.25, {
				Sprite: 'censor_' + data['race'],
				Color:  data['skin'],
				X:      rotate[frame % 4],
				Y:      1
			});
		}
	};
	f.click = function(x, y) {
		if (x >= -4 && x < 4) {
			if (y >= -1.75 && y <= -1.25) {
				send({'CharacterCreation': {'Command': 'race'}});
				return false;
			} else if (y >= -0.75 && y <= -0.25) {
				send({'CharacterCreation': {'Command': 'gender'}});
				return false;
			} else if (y >= 0 && y <= 1) {
				send({'CharacterCreation': {'Command': 'skin'}});
				return false;
			} 
		}
	};
	return f;
};

var lostConnectionHud = function(draw) {
	for (var x = w/2 - 4; x < w/2 + 4; x++) {
		draw(x, h/2, {
			Sprite: 'ui_r1',
			Color:  '#fff'
		});
	}
	draw(w/2 - 2, h/2 - 0.25, {
		Text:  'Lost connection.',
		Color: '#666'
	});
};
