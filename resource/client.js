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

var zoneBuffer = document.createElement('canvas');
zoneBuffer.width = 256*tileSize;
zoneBuffer.height = 256*tileSize;
var zoneCtx = zoneBuffer.getContext('2d');
var zoneBufferDirty = false;

var frame = 0;
setInterval(function() {
	frame = +new Date() / 50;
	repaint();
}, 50);

var inRepaint = false;
function repaint() {
	if (inRepaint)
		return;

	requestAnimationFrame(function() {
		inRepaint = false;

		canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);

		var draw = function(x, y, p, ctx) {
			ctx = ctx || canvas;
			var woff = w*tileSize/2;
			var hoff = h*tileSize/2;
			if (ctx != canvas) {
				woff = 0;
				hoff = 0;
			}
			var scale = p.Scale || 1;
			var height = p.Height || tileSize;
			if (p.Text) {
				var tx = Math.floor(woff + x*tileSize + tileSize/8);
				var ty = Math.floor(hoff + y*tileSize + tileSize*7/8);
				ctx.font = Math.floor(scale * 16) + 'px rnoadm-expection';
				ctx.fillStyle = '#000';
				ctx.fillText(p.Text, tx, ty + 1);
				ctx.fillStyle = p.Color;
				ctx.fillText(p.Text, tx, ty);
			}
			if (p.Sprite) {
				if (!images[p.Sprite]) {
					images[p.Sprite] = true;
					var img = new Image();
					img.onload = function() {
						images[p.Sprite] = img;
						images_resize[p.Sprite] = {};
						images_recolor[p.Sprite] = {};
						if (ctx != canvas) {
							zoneBufferDirty = true;
						}
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
				ctx.drawImage(images_recolor[p.Sprite][scale][p.Color],
					Math.floor((p.X || 0) * tileSize * scale),
					Math.floor((p.Y || 0) * height * scale),
					Math.floor(tileSize * scale),
					Math.floor(height * scale),
					Math.floor(woff + x*tileSize),
					Math.floor(hoff + y*tileSize + (tileSize - height) * scale),
					Math.floor(tileSize * scale),
					Math.floor(height * scale));
			}
		};

		var playerX = gameState.playerX || 0;
		var playerY = gameState.playerY || 0;

		if (gameState.objects) {
			if (zoneBufferDirty) {
				zoneBufferDirty = false;
				zoneCtx.clearRect(0, 0, zoneBuffer.width, zoneBuffer.height);
				for (var x = 0; x < 256; x += 4) {
					for (var y = 0; y < 256; y += 4) {
						draw(x, y, {
							Sprite: 'grass_r1',
							Color:  '#268f1e',
							Scale:  4,
							X:      ((x + y + x*y/4) / 4) % 4
						}, zoneCtx);
					}
				}
				for (var i in gameState.objects) {
					var obj = gameState.objects[i];
					var drawObject = function(o) {
						o.colors.forEach(function(color, j) {
							if (color) {
								draw(obj.x, obj.y, {
									Sprite: o.sprite,
									Color:  color,
									Scale:  o.scale,
									Height: o.height,
									Y:      j
								}, zoneCtx);
							}
						});
						o.attach.forEach(drawObject);
					};

					drawObject(obj.object);
				}
			}
			canvas.drawImage(zoneBuffer,
				Math.floor((w/2 - playerX) * tileSize),
				Math.floor((h/2 - playerY) * tileSize));
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
		if (msg['SetHUD']['Name']) {
			gameState.hud = huds[msg['SetHUD']['Name']](msg['SetHUD']['Data']);
		} else {
			delete gameState.hud;
		}
		repaint();
	}
	if (msg['Kick']) {
		ws.onclose = wsonclose = function() {};
		gameState = {};
		repaint();
		alert('Kicked: ' + msg['Kick']);
	}
	if (msg['ResetZone']) {
		zoneBufferDirty = true;
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
		zoneBufferDirty = true;
		var toObject = function(o) {
			return {
				sprite: o['I'],
				colors: o['C'],
				scale:  o['S'],
				height: o['H'],
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
var mouseX = -w, mouseY = -h;
document.querySelector('canvas').onmouseout = function() {
	mouseX = -w;
	mouseY = -h;
	repaint();
};
document.querySelector('canvas').onmousemove = function(e) {
	mouseX = Math.floor(e.offsetX * 4 / tileSize)/4 - w/2;
	mouseY = Math.floor(e.offsetY * 4 / tileSize)/4 - h/2;
	repaint();
};

var loginHudUsername = localStorage['login'] || '';
var loginHudPassword = '';
var loginHudFocus = 0;
var loginHud = function(draw) {
	for (var x = -4; x < 4; x++) {
		for (var y = -4; y < 1; y++) {
			draw(x, y, {
				Sprite: 'ui_r1',
				Color:  y == -4 ? '#444' : '#222',
				X:      y == -4 ? x == -4 ? 3 : x == 3 ? 4 : 0 : 0
			});
		}
		draw(x, 0.5, {
			Sprite: 'ui_r1',
			Color:  '#444',
			Y:      1
		});
		draw(x + 0.5, 0.5, {
			Sprite: 'ui_r1',
			Color:  '#444',
			Y:      1
		});
	}
	for (var x = -3.75; x < 3.75; x += 0.5) {
		draw(x, -2.75, {
			Sprite: 'ui_r1',
			Color:  '#111',
			Y:      1
		});
		draw(x, -1.25, {
			Sprite: 'ui_r1',
			Color:  '#111',
			Y:      1
		});
	}
	draw(-1.625, -4, {
		Text:   'R',
		Color:  '#888',
		Scale:  2
	});
	var permutations = 'ando'.split('').sort(function(a, b) { return Math.random() * 2 - 1; });
	for (var i in permutations) {
		draw(-1 + i/2, -4, {
			Text:   permutations[i],
			Color:  '#888',
			Scale:  2
		});
	}
	draw(1, -4, {
		Text:  'm',
		Color: '#888',
		Scale: 2
	});
	draw(-3.75, -3.25, {
		Text:  'Login',
		Color: '#aaa'
	});
	draw(-3.75, -2.75, {
		Text:  loginHudUsername + (loginHudFocus === 0 ? '_' : ''),
		Color: '#fff'
	});
	draw(-3.75, -1.75, {
		Text:  'Password',
		Color: '#aaa'
	});
	draw(-3.75, -1.25, {
		Text:  loginHudPassword.replace(/./g, '*') + (loginHudFocus === 1 ? '_' : ''),
		Color: '#fff'
	});
	draw(-3, -0.25, {
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
		for (var x = -6; x < 6; x++) {
			draw(x, -5, {
				Sprite: 'ui_r1',
				Color:  '#444',
				X:      x == -6 ? 3 : x == 5 ? 4 : 0
			});
			for (var y = -4; y < 2; y++) {
				draw(x, y, {
					Sprite: 'ui_r1',
					Color:  x >= -5 && x < -1 && y < 0 ? '#ccc' : '#222'
				});
			}
		}
		for (var x = 1; x < 5; x += 0.5) {
			draw(x, 1.5, {
				Sprite: 'ui_r1',
				Color:  '#444',
				X:      x == 1 ? 1 : x == 4.5 ? 2 : 0,
				Y:      1
			});
		}
		draw(2.25, 1.5, {
			Text:  'Accept',
			Color: mouseX >= -1 && mouseX < 5 && mouseY >= 2 && mouseY < 2.5 ? '#fff' : '#aaa'
		});
		draw(-2.75, -5, {
			Text:  'Character Creation',
			Color: '#888',
			Scale: 2
		});
		draw(0, -4, {
			Text:  'Race:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -3.75 && mouseY <= -3.25 ? '#fff' : '#aaa'
		});
		draw(2, -4, {
			Text:  data['race'],
			Color: '#fff'
		});
		draw(0, -3, {
			Text:  'Gender:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -2.75 && mouseY <= -2.25 ? '#fff' : '#aaa'
		});
		draw(2, -3, {
			Text:  data['gender'],
			Color: '#fff'
		});
		draw(0, -2, {
			Text:  'Skin:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -2 && mouseY <= -1 ? '#fff' : '#aaa'
		});
		draw(2.125, -2.125, {
			Sprite: 'ui_r1',
			Color:  data['skin'],
			Y:      1
		});
		draw(0, -1, {
			Text:  'Shirt:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -1 && mouseY <= 0 ? '#fff' : '#aaa'
		});
		draw(2.125, -1.125, {
			Sprite: 'ui_r1',
			Color:  data['shirt'],
			Y:      1
		});
		draw(0, 0, {
			Text:  'Pants:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= 0 && mouseY <= 1 ? '#fff' : '#aaa'
		});
		draw(2.125, -0.125, {
			Sprite: 'ui_r1',
			Color:  data['pants'],
			Y:      1
		});
		var rotate = [0, 6, 3, 9];
		draw(-5, -4, {
			Sprite: 'body_' + data['race'],
			Color:  data['skin'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(-5, -4, {
			Sprite: 'shoes_basic',
			Color:  '#eef8f0',
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(-5, -4, {
			Sprite: 'pants_basic',
			Color:  data['pants'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(-5, -4, {
			Sprite: 'shirt_basic',
			Color:  data['shirt'],
			X:      rotate[Math.floor(frame/10) % 4],
			Scale:  4
		});
		draw(-5, 0, {
			Text:  'Name:',
			Color: mouseX >= -5 && mouseX < 0 && mouseY >= 0.25 && mouseY <= 0.75 ? '#fff' : '#aaa'
		});
		draw(-3.5, 0, {
			Text:  data['name'],
			Color: '#fff'
		});
	};
	f.click = function(x, y) {
		if (x >= -1 && x < 5 && y >= 2 && y < 2.5) {
			send({'CharacterCreation': {'Command': 'accept'}});
			return false;
		} else if (x >= 0 && x < 6) {
			if (y >= -3.75 && y <= -3.25) {
				send({'CharacterCreation': {'Command': 'race'}});
				return false;
			} else if (y >= -2.75 && y <= -2.25) {
				send({'CharacterCreation': {'Command': 'gender'}});
				return false;
			} else if (y >= -2 && y <= -1) {
				send({'CharacterCreation': {'Command': 'skin'}});
				return false;
			} else if (y >= -1 && y <= 0) {
				send({'CharacterCreation': {'Command': 'shirt'}});
				return false;
			} else if (y >= 0 && y <= 1) {
				send({'CharacterCreation': {'Command': 'pants'}});
				return false;
			} 
		} else if (x >= -5 && x < 0 && y >= 0.25 && y <= 0.75) {
			send({'CharacterCreation': {'Command': 'name'}});
			return false;
		}
	};
	return f;
};

var lostConnectionHud = function(draw) {
	for (var x = -4; x < 4; x++) {
		draw(x, 0, {
			Sprite: 'ui_r1',
			Color:  '#fff'
		});
	}
	draw(-2, -0.25, {
		Text:  'Lost connection.',
		Color: '#666'
	});
};
