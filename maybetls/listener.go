package maybetls

import (
	"crypto/tls"
	"net"
)

// Listener returns a net.Listener that calls the Conn function in this package
// on accepted connections.
func Listener(ln net.Listener, config *tls.Config) net.Listener {
	return &listener{ln, config}
}

type listener struct {
	net.Listener
	config *tls.Config
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return Conn(c, l.config)
}
