(function() {
const GRAPHICS_REVISION = 6;
const tileSize = 32;
var ws,
clientHash,
username,
password,
loggedIn = false,
connected = false,
inRepaint = true,
w = 32, h = 16,
canvas = document.createElement('canvas').getContext('2d'),
floor = function(n) {
	return Math.floor(n);
},
send = function(packet) {
	if (!connected)
		return;

	ws.send(JSON.stringify(packet));
},
wsopen = function() {
	connected = true;
	if (loggedIn) {
		send({'Auth': {'U': username, 'P': password}});
	}
},
wsclose = function() {
	connected = false;
	setTimeout(connect, 100);
},
wsmessage = function(e) {
	var msg = JSON.parse(e.data), p;
	if (p = msg['Kick']) {
		ws.onclose = wsopen = wsclose = wsmessage = function() {};
		ws.close();
		alert(p);
	}
	if (p = msg['ClientHash']) {
		if (clientHash) {
			if (p != clientHash) {
				location.reload(true);
			}
		} else {
			clientHash = p;
		}
	}
	console.log(msg);
},
connect = function() {
	ws = new WebSocket('ws://' + location.host + '/ws');
	ws.onopen = wsopen;
	ws.onclose = wsclose;
	ws.onmessage = wsmessage;
},
repaint = function() {
	if (inRepaint) {
		return;
	}

	inRepaint = true;
	requestAnimationFrame(paint);
},
images = {},
images_resize = {},
images_recolor = {},
drawSprite = function(x, y, sprite, color, extra, ctx) {
	extra = extra || {};
	ctx = ctx || canvas;
	var xoff = ctx == canvas ? w/2*tileSize : 0;
	var yoff = ctx == canvas ? h/2*tileSize : 0;
	var scale = extra['scale'] || 1;
	var height = extra['height'] || tileSize;
	var width = extra['width'] || tileSize;
	var col = extra['x'] || 0;
	var row = extra['y'] || 0;
	if (!images[sprite]) {
		images[sprite] = true;
		var img = new Image();
		img.onload = function() {
			images[sprite] = img;
			images_resize[sprite] = {};
			images_recolor[sprite] = {};
			repaint();
		};
		img.src = sprite + '.png?' + GRAPHICS_REVISION;
	}
	if (images[sprite] === true)
		return;
	if (!images_resize[sprite][scale]) {
		var buffer = document.createElement('canvas');
		buffer.width = images[sprite].width;
		buffer.height = images[sprite].height;
		images_resize[sprite][scale] = buffer;
		images_recolor[sprite][scale] = {};
		buffer = buffer.getContext('2d');
		buffer.drawImage(images[sprite], 0, 0);
		var base = buffer.getImageData(0, 0, images[sprite].width, images[sprite].height);
		buffer = buffer.canvas;
		buffer.width = images[sprite].width * scale;
		buffer.height = images[sprite].height * scale;
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
	if (!images_recolor[sprite][scale][color]) {
		var buffer = document.createElement('canvas');
		buffer.width = images_resize[sprite][scale].width || 1;
		buffer.height = images_resize[sprite][scale].height || 1;
		images_recolor[sprite][scale][color] = buffer;
		buffer = buffer.getContext('2d');

		buffer.fillStyle = color;
		buffer.fillRect(0, 0, 1, 1);
		var data = buffer.getImageData(0, 0, 1, 1);
		var r = data.data[0],
		    g = data.data[1],
		    b = data.data[2],
		    a = data.data[3];
		buffer.clearRect(0, 0, 1, 1);

		buffer.drawImage(images_resize[sprite][scale], 0, 0);
		data = buffer.getImageData(0, 0, buffer.canvas.width, buffer.canvas.height);
		var fade = function(x, y) {
			if (x >= 128)
				return 255 - (255-x)*(255-y)/127;
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
	ctx.drawImage(images_recolor[sprite][scale][color],
		floor(col * width * scale),
		floor(row * height * scale),
		floor(width * scale),
		floor(height * scale),
		floor(xoff + x*tileSize),
		floor(yoff + y*tileSize + (tileSize - height) * scale),
		floor(width * scale),
		floor(height * scale));
},
paint = function() {
	inRepaint = false;
	canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);

	drawSprite(0, 0, 'ui_r1', '#f0f');
},
loginForm = document.querySelector('form'),
loginField = loginForm.querySelector('#username'),
passField = loginForm.querySelector('#password'),
pass2Field = loginForm.querySelector('#password2');

onresize = function() {
	w = floor(innerWidth / tileSize);
	h = floor(innerHeight / tileSize);
	canvas.canvas.width  = tileSize*w;
	canvas.canvas.height = tileSize*h;
	repaint();
};
onresize();

this['admin'] = function(command) {
	send({'Admin': command});
};

passField.onchange = function() {
	pass2Field.value = passField.value;
};
loginForm.onsubmit = function() {
	if (loggedIn) {
		return;
	}
	username = loginField.value;
	password = passField.value;
	if (!username) {
		loginField.focus();
		return;
	}
	if (password.length <= 2) {
		passField.focus();
		return;
	}
	loggedIn = true;
	inRepaint = false;
	send({'Auth': {'U': username, 'P': password}});
	var parent = loginForm.parentNode;
	parent.removeChild(loginForm);
	parent.appendChild(canvas.canvas);
	repaint();
};

connect();
})()
