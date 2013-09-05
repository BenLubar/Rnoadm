package maybetls

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
)

// Conn returns a net.Conn that is wrapped by the crypto/tls.Server function if
// a heursitic is satisfied. See the package documentation for more details.
func Conn(c net.Conn, config *tls.Config) (net.Conn, error) {
	var b [1]byte
	_, err := io.ReadFull(c, b[:])
	if err != nil {
		return nil, err
	}

	c = &conn{io.MultiReader(bytes.NewReader(b[:]), c), c}

	if b[0] == 22 { // TLS handshake starts with ASCII SYN
		return tls.Server(c, config), nil
	}
	return c, nil
}

type conn struct {
	io.Reader
	net.Conn
}

func (c *conn) Read(b []byte) (n int, err error) {
	return c.Reader.Read(b)
}
