var ws,
clientHash,
username,
password,
loggedIn = false,
connected = false,
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
		wsopen = wsclose = wsmessage = function() {};
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
};

window['admin'] = function(command) {
	send({'Admin': command});
};
connect();
(function(form) {
	var loginField = form.querySelector('#username');
	var passField = form.querySelector('#password');
	var pass2Field = form.querySelector('#password2');
	passField.onchange = function() {
		pass2Field.value = passField.value;
	};
	form.onsubmit = function() {
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
		send({'Auth': {'U': username, 'P': password}});
		form.parentNode.removeChild(form);
	};
})(document.querySelector('form'));
