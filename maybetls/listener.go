package maybetls

import (
	"crypto/tls"
	"net"
)

func Listener(ln net.Listener, config *tls.Config) net.Listener {
	return &listener{ln, config}
}

type listener struct {
	ln     net.Listener
	config *tls.Config
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.ln.Accept()
	if err != nil {
		return nil, err
	}

	return Conn(c, l.config)
}

func (l *listener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *listener) Close() error {
	return l.ln.Close()
}
