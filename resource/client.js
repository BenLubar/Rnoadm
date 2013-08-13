(function() {
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
draw = function(x, y, info) {
},
paint = function() {
	inRepaint = false;
	canvas.clearRect(0, 0, canvas.canvas.width, canvas.canvas.height);
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
};

connect();
})()
