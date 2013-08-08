var undefined;
const tileSize = 32;
var w = 32, h = 16;
var gameState = {};
var images = {};
var images_resize = {};
var images_recolor = {};
var clientHash;
var huds = {};
var canvas = document.querySelector('canvas');
canvas.width = w*tileSize;
canvas.height = h*tileSize;
canvas = canvas.getContext('2d');

var zoneBufferStatic = document.createElement('canvas');
var zoneBufferDynamic = document.createElement('canvas');
zoneBufferStatic.width = 256*tileSize;
zoneBufferDynamic.width = 256*tileSize;
zoneBufferStatic.height = 256*tileSize;
zoneBufferDynamic.height = 256*tileSize;
var zoneCtxStatic = zoneBufferStatic.getContext('2d');
var zoneCtxDynamic = zoneBufferDynamic.getContext('2d');
var zoneBufferStaticDirty = false;
var zoneBufferDynamicDirty = false;

const color_222 = '#222';
const color_444 = '#444';
const color_888 = '#888';
const color_aaa = '#aaa';
const color_fff = '#fff';

var frame = 0;
setInterval(function() {
	frame = +new Date() / 50;
	repaint();
}, 50);

setInterval(function() {
	inRepaint = false;
}, 60000);

function drawObject(draw, x, y, ctx, o) {
	o.colors.forEach(function(color, j) {
		if (color) {
			draw(x, y, {
				Sprite: o.sprite,
				Color:  color,
				Scale:  o.scale,
				Height: o.height,
				Y:      j
			}, ctx);
		}
	});
	o.attach.forEach(function(a) {
		drawObject(draw, x, y, ctx, a);
	});
}

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
				if (p.Title) {
					ctx.font = Math.floor(scale * 32) + 'px "Jolly Lodger"';
				} else {
					ctx.font = Math.floor(scale * 14) + 'px "Open Sans Condensed"';
				}
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
							zoneBufferStaticDirty = true;
							zoneBufferDynamicDirty = true;
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

		var playerX = getPlayerX();
		var playerY = getPlayerY();

		if (gameState.objects) {
			if (zoneBufferStaticDirty) {
				zoneBufferStaticDirty = false;
				zoneCtxStatic.clearRect(0, 0, 256*tileSize, 256*tileSize);
				for (var x = 0; x < 256; x++) {
					for (var y = 0; y < 256; y++) {
						draw(x, y, {
							Sprite: 'grass_r1',
							Color:  '#268f1e',
							X:      (x + y + x*y) % 4
						}, zoneCtxStatic);
					}
				}
				for (var i in gameState.objects) {
					var obj = gameState.objects[i];

					if (!obj.object.moves) {
						drawObject(draw, animateCoord(obj.x, obj.xnext, obj.frame), animateCoord(obj.y, obj.ynext, obj.frame), zoneCtxStatic, obj.object);
					}
				}
			}
			if (zoneBufferDynamicDirty) {
				zoneBufferDynamicDirty = false;
				zoneCtxDynamic.clearRect(0, 0, 256*tileSize, 256*tileSize);
				for (var i in gameState.objects) {
					var obj = gameState.objects[i];

					if (obj.object.moves) {
						drawObject(draw, animateCoord(obj.x, obj.xnext, obj.frame), animateCoord(obj.y, obj.ynext, obj.frame), zoneCtxDynamic, obj.object);
					}
				}
			}
			canvas.drawImage(zoneBufferStatic,
				Math.floor((w/2 - playerX) * tileSize),
				Math.floor((h/2 - playerY) * tileSize));
			canvas.drawImage(zoneBufferDynamic,
				Math.floor((w/2 - playerX) * tileSize),
				Math.floor((h/2 - playerY) * tileSize));
		}

		(gameState.inventory || []).forEach(function(item, i) {
			const cols = 8;
			drawObject(draw, w / 2 - cols + i % cols, h / 2 - Math.ceil(gameState.inventory.length / cols) + Math.floor(i / cols), undefined, item.object);
		});

		var y = h / 2 - 1;
		(gameState.messages || []).forEach(function(message) {
			var lines = message.split(/\n/g);
			y -= lines.length / 2;
			lines.forEach(function(line, i) {
				draw(-w / 2, y + i / 2, {
					Text:  line,
					Color: color_fff
				});
			});
		});

		if (gameState.hud) {
			gameState.hud(draw);
		}
	});
}

function animateCoord(start, end, anim) {
	start = start || 0;
	end   = isNaN(end) ? start : end;
	anim  = frame - (anim || 0);
	if (anim > 10) {
		return end;
	}
	zoneBufferDynamicDirty = true;
	return (start * (10 - anim) + end * anim) / 10;
}

function getPlayerX() {
	return animateCoord(gameState.playerX, gameState.playerXNext, gameState.playerXFrame);
}

function getPlayerY() {
	return animateCoord(gameState.playerY, gameState.playerYNext, gameState.playerYFrame);
}

var ws = new WebSocket('ws://' + location.host + '/ws');
var wsonopen = ws.onopen = function() {
	gameState = {
		hud: loginHud
	};
	loginHudFocus = loginHudUsername == '' ? 0 : 1;
	repaint();
};

function toObject(o) {
	return {
		name:    o['N'],
		options: o['O'] || [],
		sprite:  o['I'],
		colors:  o['C'],
		scale:   o['S'],
		height:  o['H'],
		moves:   !!o['M'],
		attach:  (o['A'] || []).map(toObject)
	};
}
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
	if (msg['Message']) {
		if (!gameState.messages) {
			gameState.messages = [];
		}
		if (gameState.messages.length >= 8) {
			gameState.messages.pop();
		}
		gameState.messages.unshift(msg['Message']);
		repaint();
	}
	if (msg['ResetZone']) {
		zoneBufferStaticDirty = true;
		zoneBufferDynamicDirty = true;
		gameState.objects = {};
		repaint();
	}
	if ('PlayerX' in msg && gameState.playerXNext != msg['PlayerX']) {
		gameState.playerX = gameState.playerXNext || msg['PlayerX'];
		gameState.playerXNext = msg['PlayerX'];
		gameState.playerXFrame = frame;
		repaint();
	}
	if ('PlayerY' in msg && gameState.playerYNext != msg['PlayerY']) {
		gameState.playerY = gameState.playerYNext || msg['PlayerY'];
		gameState.playerYNext = msg['PlayerY'];
		gameState.playerYFrame = frame;
		repaint();
	}
	if (msg['TileChange']) {
		msg['TileChange'].forEach(function(tile) {
			if (tile['R']) {
				if (gameState.objects[tile['ID']].object.moves) {
					zoneBufferDynamicDirty = true;
				} else {
					zoneBufferStaticDirty = true;
				}
				delete gameState.objects[tile['ID']];
			} else {
				if (tile['O']) {
					gameState.objects[tile['ID']] = {
						x:      tile['X'],
						xnext:  tile['X'],
						y:      tile['Y'],
						ynext:  tile['Y'],
						frame:  0,
						object: toObject(tile['O'])
					};
					if (gameState.objects[tile['ID']].object.moves) {
						zoneBufferDynamicDirty = true;
					} else {
						zoneBufferStaticDirty = true;
					}
				} else {
					gameState.objects[tile['ID']].x = gameState.objects[tile['ID']].xnext;
					gameState.objects[tile['ID']].y = gameState.objects[tile['ID']].ynext;
					gameState.objects[tile['ID']].xnext = tile['X'];
					gameState.objects[tile['ID']].ynext = tile['Y'];
					gameState.objects[tile['ID']].frame = frame;
					zoneBufferDynamicDirty = true;
				}
			}
		});
		repaint();
	}
	if (msg['Inventory']) {
		gameState.inventory = [];
		msg['Inventory'].forEach(function(item) {
			gameState.inventory.push({
				id:     item['I'],
				object: toObject(item['O'])
			});
		});
		repaint();
	}
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
	switch (e.keyCode) {
	case 38: // up
		send({'Walk': {'X': gameState.playerXNext, 'Y': gameState.playerYNext - 1}});
		break;
	case 40: // down
		send({'Walk': {'X': gameState.playerXNext, 'Y': gameState.playerYNext + 1}});
		break;
	case 37: // left
		send({'Walk': {'X': gameState.playerXNext - 1, 'Y': gameState.playerYNext}});
		break;
	case 39: // right
		send({'Walk': {'X': gameState.playerXNext + 1, 'Y': gameState.playerYNext}});
		break;
	}
};
document.onkeypress = function(e) {
	if (gameState.hud && gameState.hud.keyPress) {
		if (gameState.hud.keyPress(e.charCode) === false)
			return;
	}
};
canvas.canvas.onclick = function(e) {
	var x = Math.floor(e.offsetX * 4 / tileSize)/4 - w/2;
	var y = Math.floor(e.offsetY * 4 / tileSize)/4 - h/2;

	if (gameState.hud && gameState.hud.click) {
		if (gameState.hud.click(x, y) === false)
			return false;
	}

	var wx = Math.floor(x + getPlayerX());
	var wy = Math.floor(y + getPlayerY());
	if (wx >= 0 && wx < 256 && wy >= 0 && wy < 256) {
		send({'Walk': {'X': wx, 'Y': wy}});
	}
	return false;
}
canvas.canvas.oncontextmenu = function(e) {
	var x = Math.floor(e.offsetX * 4 / tileSize)/4 - w/2;
	var y = Math.floor(e.offsetY * 4 / tileSize)/4 - h/2;

	if (gameState.hud && gameState.hud.click) {
		if (gameState.hud.click(x, y) === false)
			return false;
	}

	var wx = Math.floor(x + getPlayerX());
	var wy = Math.floor(y + getPlayerY());
	if (wx >= 0 && wx < 256 && wy >= 0 && wy < 256) {
		gameState.hud = rightClickHud(wx, wy, x, y);
	}
	return false;
};
var mouseX = -w, mouseY = -h;
canvas.canvas.onmouseout = function() {
	mouseX = -w;
	mouseY = -h;
	repaint();
};
canvas.canvas.onmousemove = function(e) {
	mouseX = Math.floor(e.offsetX * 4 / tileSize)/4 - w/2;
	mouseY = Math.floor(e.offsetY * 4 / tileSize)/4 - h/2;
	repaint();
};

var loginHudUsername = localStorage['login'] || '';
var loginHudPassword = '';
var loginHudFocus = 0;
var loginHudPermutations = 'ando'.split('');
var loginHudPermutationsFrame = 0;
var loginHud = function(draw) {
	for (var x = -4; x < 4; x++) {
		for (var y = -4; y < 1; y++) {
			draw(x, y, {
				Sprite: 'ui_r1',
				Color:  y == -4 ? color_444 : color_222,
				X:      y == -4 ? x == -4 ? 3 : x == 3 ? 4 : 0 : 0
			});
		}
		draw(x, 0.5, {
			Sprite: 'ui_r1',
			Color:  color_444,
			Y:      1
		});
		draw(x + 0.5, 0.5, {
			Sprite: 'ui_r1',
			Color:  color_444,
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
		Text:  'R',
		Color: color_888,
		Title: true
	});
	if (loginHudPermutationsFrame < frame - 25) {
		for (var i in loginHudPermutations) {
			var tmp = loginHudPermutations[i];
			var j = Math.floor(Math.random() * loginHudPermutations.length);
			loginHudPermutations[i] = loginHudPermutations[j];
			loginHudPermutations[j] = tmp;
		}
		loginHudPermutationsFrame = frame;
	}
	for (var i in loginHudPermutations) {
		draw(-1 + i/2, -4, {
			Text:  loginHudPermutations[i],
			Color: color_888,
			Title: true
		});
	}
	draw(1, -4, {
		Text:  'm',
		Color: color_888,
		Title: true
	});
	draw(-3.75, -3.25, {
		Text:  'Login',
		Color: mouseX >= -4 && mouseX < 4 && mouseY >= -2.75 && mouseY <= -1.75 ? color_fff : color_aaa
	});
	draw(-3.75, -2.75, {
		Text:  loginHudUsername + (loginHudFocus === 0 ? '_' : ''),
		Color: color_fff
	});
	draw(-3.75, -1.75, {
		Text:  'Password',
		Color: mouseX >= -4 && mouseX < 4 && mouseY >= -1.25 && mouseY <= -0.25 ? color_fff : color_aaa
	});
	draw(-3.75, -1.25, {
		Text:  loginHudPassword.replace(/./g, '*') + (loginHudFocus === 1 ? '_' : ''),
		Color: color_fff
	});
	draw(-3, -0.25, {
		Text:  'Log in or register',
		Color: mouseX >= -4 && mouseX < 4 && mouseY >= 0.00 && mouseY <= 0.75 ? color_fff : color_aaa
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
	}
	return false;
};

loginHud.keyDown = function(code) {
	switch (code) {
	case 8: // backspace
		switch (loginHudFocus) {
		case 0:
			if (loginHudUsername) {
				loginHudUsername = loginHudUsername.substring(0, loginHudUsername.length - 1);
				localStorage['login'] = loginHudUsername;
			}
			break;
		case 1:
			if (loginHudPassword) {
				loginHudPassword = loginHudPassword.substring(0, loginHudPassword.length - 1);
			}
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
	if (loginHudPassword.length < 3) {
		loginHudFocus = 1;
		repaint();
		return;
	}
	send({'Auth': {'Login': loginHudUsername, 'Password': loginHudPassword}});
	loginHudPassword = '';
	loginHudFocus = 0;
}

var rightClickHud = function(wx, wy, sx, sy) {
	var options = [];
	for (var i in gameState.objects) {
		var o = gameState.objects[i];
		if (o.xnext == wx && o.ynext == wy) {
			o.object.options.forEach(function(option, j) {
				options.push({
					id:   i,
					obj:  o,
					name: o.object.name,
					cmd:  option,
					oid:  j
				});
			});
			options.push({
				id:   i,
				obj:  o,
				name: o.object.name,
				cmd:  'examine',
				oid:  -1
			});
		}
	}
	var f = function(draw) {
		if (mouseX < sx - 1 || mouseX > sx + 6 || mouseY < sy - 1 || mouseY > sy + options.length / 2 + 0.5) {
			delete gameState.hud;
			return;
		}
		options.forEach(function(option, y) {
			for (var x = 0; x < 10; x++) {
				draw(sx + x / 2, sy + (y - 1) / 2, {
					Sprite: 'ui_r1',
					Color:  mouseX >= sx && mouseX < sx + 5 && mouseY >= sy + y / 2 && mouseY < sy + (y + 1) / 2 ? color_444 : color_222,
					Y:      1
				});
			}
			draw(sx, sy + (y - 1) / 2, {
				Text:  option.cmd + ' ' + option.name,
				Color: mouseX >= sx && mouseX < sx + 5 && mouseY >= sy + y / 2 && mouseY < sy + (y + 1) / 2 ? color_fff : color_aaa
			});
		});
	};
	f.click = function(x, y) {
		options.forEach(function(option, i) {
			if (x >= sx && x < sx + 5 && y >= sy + i / 2 && y < sy + (i + 1) / 2) {
				send({'Interact': {
					'I': option.id,
					'O': option.oid,
					'X': wx,
					'Y': wy
				}});
			}
		});
		delete gameState.hud;
		repaint();
		return false;
	};
	return f;
};

huds['character_creation'] = function(data) {
	var f = function(draw) {
		gameState.playerXNext = gameState.playerX = 127 + Math.cos(frame / 10000 * 7) * 64;
		gameState.playerYNext = gameState.playerY = 127 + Math.sin(frame / 10000 * 6) * 64;
		for (var x = -6; x < 6; x++) {
			draw(x, -5, {
				Sprite: 'ui_r1',
				Color:  color_444,
				X:      x == -6 ? 3 : x == 5 ? 4 : 0
			});
			for (var y = -4; y < 2; y++) {
				draw(x, y, {
					Sprite: 'ui_r1',
					Color:  x >= -5 && x < -1 && y < 0 ? '#ccc' : color_222
				});
			}
		}
		for (var x = 1; x < 5; x += 0.5) {
			draw(x, 1.5, {
				Sprite: 'ui_r1',
				Color:  color_444,
				X:      x == 1 ? 1 : x == 4.5 ? 2 : 0,
				Y:      1
			});
		}
		draw(2.25, 1.5, {
			Text:  'Accept',
			Color: mouseX >= -1 && mouseX < 5 && mouseY >= 2 && mouseY < 2.5 ? color_fff : color_aaa
		});
		draw(-2.75, -5, {
			Text:  'Character Creation',
			Color: color_888,
			Title: true
		});
		draw(0, -4, {
			Text:  'Race:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -3.75 && mouseY <= -3.25 ? color_fff : color_aaa
		});
		draw(2, -4, {
			Text:  data['race'],
			Color: color_fff
		});
		draw(0, -3, {
			Text:  'Gender:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -2.75 && mouseY <= -2.25 ? color_fff : color_aaa
		});
		draw(2, -3, {
			Text:  data['gender'],
			Color: color_fff
		});
		draw(0, -2, {
			Text:  'Skin:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -2 && mouseY <= -1 ? color_fff : color_aaa
		});
		draw(2.125, -2.125, {
			Sprite: 'ui_r1',
			Color:  data['skin'],
			Y:      1
		});
		draw(0, -1, {
			Text:  'Shirt:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= -1 && mouseY <= 0 ? color_fff : color_aaa
		});
		draw(2.125, -1.125, {
			Sprite: 'ui_r1',
			Color:  data['shirt'],
			Y:      1
		});
		draw(0, 0, {
			Text:  'Pants:',
			Color: mouseX >= 0 && mouseX < 6 && mouseY >= 0 && mouseY <= 1 ? color_fff : color_aaa
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
			Color: mouseX >= -5 && mouseX < 0 && mouseY >= 0.25 && mouseY <= 0.75 ? color_fff : color_aaa
		});
		draw(-3.5, 0, {
			Text:  data['name'],
			Color: color_fff
		});
	};
	f.click = function(x, y) {
		if (x >= -1 && x < 5 && y >= 2 && y < 2.5) {
			send({'CharacterCreation': {'Command': 'accept'}});
		} else if (x >= 0 && x < 6) {
			if (y >= -3.75 && y <= -3.25) {
				send({'CharacterCreation': {'Command': 'race'}});
			} else if (y >= -2.75 && y <= -2.25) {
				send({'CharacterCreation': {'Command': 'gender'}});
			} else if (y >= -2 && y <= -1) {
				send({'CharacterCreation': {'Command': 'skin'}});
			} else if (y >= -1 && y <= 0) {
				send({'CharacterCreation': {'Command': 'shirt'}});
			} else if (y >= 0 && y <= 1) {
				send({'CharacterCreation': {'Command': 'pants'}});
			} 
		} else if (x >= -5 && x < 0 && y >= 0.25 && y <= 0.75) {
			send({'CharacterCreation': {'Command': 'name'}});
		}
		return false;
	};
	return f;
};

var lostConnectionHud = function(draw) {
	for (var x = -4; x < 4; x++) {
		draw(x, 0, {
			Sprite: 'ui_r1',
			Color:  color_fff
		});
	}
	draw(-2, -0.25, {
		Text:  'Lost connection.',
		Color: '#666'
	});
};

function admin(cmd) {
	send({'Admin': cmd});
}
window['admin'] = admin;
