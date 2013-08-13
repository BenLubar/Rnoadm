package main

import (
	"github.com/BenLubar/Rnoadm/resource"
	"hash/crc64"
	"net/http"
)

func init() {
	hash := crc64.New(crc64.MakeTable(crc64.ISO))
	hash.Write(resource.Resource["client.js"])
	packetClientHash.ClientHash = hash.Sum64()
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/blank" {
		return
	}
	if r.URL.Path == "/" {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Rnoadm</title>
<link href="http://fonts.googleapis.com/css?family=Jolly+Lodger|Open+Sans+Condensed:300&amp;subset=latin,latin-ext,cyrillic,cyrillic-ext,greek-ext,greek,vietnamese" rel="stylesheet">
<style>
html {
	background: #000;
}
form {
	text-align: center;
	padding-top: 20px;
}
h1 {
	font-family: 'Jolly Lodger';
	color: #fff;
	font-weight: normal;
}
iframe {
	display: none;
}
input {
	background-color: #000;
	color: #aaa;
	font: 1em 'Open Sans Condensed';
	border: 0;
	border-bottom: 1px solid #aaa;
}
input:hover, input:focus {
	color: #fff;
	border-bottom-color: #fff;
	outline: 0;
}
input[type="submit"] {
	border-bottom: 0;
}
#password2 {
	position: absolute;
	left: -99999px;
}
</style>
</head>
<body>
<iframe name="dummy" src="/blank"></iframe>
<form id="signup" method="post" action="/blank" target="dummy">
<h1>Rnoadm</h1>
<input type="text" name="username" id="username" placeholder="me@example.com" tabindex="1"><br>
<input type="password" name="password" id="password" placeholder="password" tabindex="2"><br>
<input type="password" name="password2" id="password2" tabindex="-1">
<input type="submit" tabindex="3" value="Register or log in">
</form>
<script src="client.js"></script>
</body>
</html>`))
		return
	}
	if b, ok := resource.Resource[r.URL.Path[1:]]; ok {
		w.Write(b)
		return
	}
	http.NotFound(w, r)
}
