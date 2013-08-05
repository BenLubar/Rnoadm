const tileSize = 32;
const w = 32, h = 16;
var gameState;
var images = {};
var images_resize = {};
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
	frame = Math.floor(+new Date() / 50);
	repaint();
}, 50);

var inRepaint = false;
function repaint() {
	if (inRepaint)
		return;

	requestAnimationFrame(function() {
		inRepaint = false;

		canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);

		var draw = function(x, y, p) {
			var scale = p.Scale || 1;
			var height = p.Height || tileSize;
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
						images_resize[p.Sprite] = {};
						images_recolor[p.Sprite] = {};
						repaint();
					};
					img.src = p.Sprite + '.png';
				}
				if (images[p.Sprite] === true)
					return;
				if (!images_resize[p.Sprite][scale]) {
					var buffer = document.createElement('canvas');
					buffer.width = images[p.Sprite].width;
					buffer.height = images[p.Sprite].height;
					images_resize[p.Sprite][scale] = buffer;
					images_recolor[p.Sprite][scale] = {};
					buffer = buffer.getContext('2d');
					buffer.drawImage(images[p.Sprite], 0, 0);
					var base = buffer.getImageData(0, 0, images[p.Sprite].width, images[p.Sprite].height);
					buffer = buffer.canvas;
					buffer.width = images[p.Sprite].width * scale;
					buffer.height = images[p.Sprite].height * scale;
					buffer = buffer.getContext('2d');
					var scaled = buffer.getImageData(0, 0, buffer.canvas.width, buffer.canvas.height);
					var rowIndex = 0;
					var baseIndex = 0;
					var scaledIndex = 0;
					for (var sy = 0; sy < buffer.canvas.height; sy++) {
						for (var sx = 0; sx < buffer.canvas.width; sx++) {
							scaled.data[scaledIndex+0] = base.data[baseIndex+0];
							scaled.data[scaledIndex+1] = base.data[baseIndex+1];
							scaled.data[scaledIndex+2] = base.data[baseIndex+2];
							scaled.data[scaledIndex+3] = base.data[baseIndex+3];
							if (sx % scale == scale - 1) {
								baseIndex += 4;
							}
							scaledIndex += 4;
						}
						if (sy % scale == scale - 1) {
							rowIndex = baseIndex;
						} else {
							baseIndex = rowIndex;
						}
					}
					buffer.putImageData(scaled, 0, 0);
				}
				if (!images_recolor[p.Sprite][scale][p.Color]) {
					var buffer = document.createElement('canvas');
					buffer.width = images_resize[p.Sprite][scale].width || 1;
					buffer.height = images_resize[p.Sprite][scale].height || 1;
					images_recolor[p.Sprite][scale][p.Color] = buffer;
					buffer = buffer.getContext('2d');
					buffer.fillStyle = p.Color;
					buffer.fillRect(0, 0, 1, 1);
					var data = buffer.getImageData(0, 0, 1, 1);
					var r = data.data[0], g = data.data[1], b = data.data[2], a = data.data[3];
					buffer.clearRect(0, 0, 1, 1);
					buffer.drawImage(images_resize[p.Sprite][scale], 0, 0);
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
				canvas.drawImage(images_recolor[p.Sprite][scale][p.Color],
					(p.X || 0) * tileSize * scale,
					(p.Y || 0) * height * scale,
					tileSize * scale,
					height * scale,
					x*tileSize,
					y*tileSize + (tileSize - height) * scale,
					tileSize * scale,
					height * scale);
			}
		};

		var playerX = gameState.playerX || 0;
		var playerY = gameState.playerY || 0;

		if (gameState.objects) {
			for (var x = 0; x < 256; x += 4) {
				for (var y = 0; y < 256; y += 4) {
					draw(x - playerX + w/2, y - playerY + h/2, {
						Sprite: 'grass_r1',
						Color:  '#268f1e',
						Scale:  4,
						X:      ((x + y + x*y/4) / 4) % 4
					});
				}
			}
			for (var i in gameState.objects) {
				var obj = gameState.objects[i];
				var drawObject = function(o) {
					o.colors.forEach(function(color, j) {
						draw(obj.x - playerX + w/2, obj.y - playerY + h/2, {
							Sprite: o.sprite,
							Color:  color,
							Scale:  o.scale,
							Y:      j
						});
					});
					o.attach.forEach(drawObject);
				};

				drawObject(obj.object);
			}
		}

		if (gameState.hud) {
			gameState.hud(draw);
		}
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
	if (msg['ResetZone']) {
		gameState.objects = {};
		repaint();
	}
	if ('PlayerX' in msg) {
		gameState.playerX = msg['PlayerX'];
		repaint();
	}
	if ('PlayerY' in msg) {
		gameState.playerY = msg['PlayerY'];
		repaint();
	}
	if (msg['TileChange']) {
		var toObject = function(o) {
			return {
				sprite: o['I'],
				colors: o['C'],
				scale:  o['S'] || 1,
				attach: (o['A'] || []).map(toObject)
			};
		};
		msg['TileChange'].forEach(function(tile) {
			if (tile['R']) {
				delete gameState.objects[tile['ID']];
			} else {
				if (tile['O']) {
					gameState.objects[tile['ID']] = {
						x:      tile['X'],
						y:      tile['Y'],
						object: toObject(tile['O'])
					};
				} else {
					gameState.objects[tile['ID']].x = tile['X'];
					gameState.objects[tile['ID']].y = tile['Y'];
				}
			}
		});
		repaint();
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
	if (gameState.hud && gameState.hud.keyDown) {
		if (gameState.hud.keyDown(e.keyCode) === false)
			return;
	}
};
document.onkeypress = function(e) {
	if (gameState.hud && gameState.hud.keyPress) {
		if (gameState.hud.keyPress(e.charCode) === false)
			return;
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
		} else if (y >= -1.25 && y <= -0.25) {
			loginHudFocus = 1; // password
			repaint();
		} else if (y >= 0.00 && y <= 0.75) {
			loginHudSubmit();
		}
		if (y >= -4 && y <= 1) {
			return false;
		}
	}
};

loginHud.keyDown = function(code) {
	switch (code) {
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
		return false;
	case 9: // tab
		loginHudFocus = (loginHudFocus + 1) % 2;
		repaint();
		return false;
	case 13: // enter
		loginHudSubmit();
		return false;
	}
};

loginHud.keyPress = function(code) {
	code = String.fromCharCode(code);
	switch (loginHudFocus) {
	case 0:
		loginHudUsername += code;
		localStorage['login'] = loginHudUsername;
		break;
	case 1:
		loginHudPassword += code;
		break;
	}
	repaint();
	return false;
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
		gameState.playerX = 127 + Math.cos(frame / 10000 * 7) * 64;
		gameState.playerY = 127 + Math.sin(frame / 10000 * 6) * 64;
		for (var x = w/2 - 6; x < w/2 + 6; x++) {
			draw(x, h/2 - 5, {
				Sprite: 'ui_r1',
				Color:  '#444',
				X:      x == w/2 - 6 ? 3 : x == w/2 + 5 ? 4 : 0
			});
			for (var y = h/2 - 4; y < h/2 + 2; y++) {
				draw(x, y, {
					Sprite: 'ui_r1',
					Color:  x >= w/2 - 5 && x < w/2 - 1 && y < h/2 ? '#ccc' : '#222'
				});
			}
		}
		draw(w/2, h/2 - 4, {
			Text:  'Race:',
			Color: '#aaa'
		});
		draw(w/2 + 2, h/2 - 4, {
			Text:  data['race'],
			Color: '#fff'
		});
		draw(w/2, h/2 - 3, {
			Text:  'Gender:',
			Color: '#aaa'
		});
		draw(w/2 + 2, h/2 - 3, {
			Text:  data['gender'],
			Color: '#fff'
		});
		draw(w/2, h/2 - 2, {
			Text:  'Skin:',
			Color: '#aaa'
		});
		draw(w/2 + 2.125, h/2 - 2.125, {
			Sprite: 'ui_r1',
			Color:  data['skin'],
			Y:      1
		});
		draw(w/2, h/2 - 1, {
			Text:  'Shirt:',
			Color: '#aaa'
		});
		draw(w/2 + 2.125, h/2 - 1.125, {
			Sprite: 'ui_r1',
			Color:  data['shirt'],
			Y:      1
		});
		draw(w/2, h/2, {
			Text:  'Pants:',
			Color: '#aaa'
		});
		draw(w/2 + 2.125, h/2 - 0.125, {
			Sprite: 'ui_r1',
			Color:  data['pants'],
			Y:      1
		});
		var rotate = [0, 6, 3, 9];
		draw(w/2 - 5, h/2 - 4, {
			Sprite: 'body_' + data['race'],
			Color:  data['skin'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(w/2 - 5, h/2 - 4, {
			Sprite: 'shoes_basic',
			Color:  '#eef8f0',
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(w/2 - 5, h/2 - 4, {
			Sprite: 'pants_basic',
			Color:  data['pants'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(w/2 - 5, h/2 - 4, {
			Sprite: 'shirt_basic',
			Color:  data['shirt'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(w/2 - 4.5, h/2, {
			Text:  data['name'],
			Color: '#fff'
		});
	};
	f.click = function(x, y) {
		if (x >= 0 && x < 6) {
			if (y >= -3.75 && y <= -3.25) {
				send({'CharacterCreation': {'Command': 'race'}});
				return false;
			} else if (y >= -2.75 && y <= -2.25) {
				send({'CharacterCreation': {'Command': 'gender'}});
				return false;
			} else if (y >= -2.125 && y <= -1) {
				send({'CharacterCreation': {'Command': 'skin'}});
				return false;
			} else if (y >= -1.125 && y <= 0) {
				send({'CharacterCreation': {'Command': 'shirt'}});
				return false;
			} else if (y >= 0.125 && y <= 1) {
				send({'CharacterCreation': {'Command': 'pants'}});
				return false;
			} 
		} else if (x >= -4 && x < 0 && y >= 0.25 && y <= 0.75) {
			send({'CharacterCreation': {'Command': 'name'}});
			return false;
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
