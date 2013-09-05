// Package maybetls implements a net.Listener wrapper that handles both secure
// (TLS) and insecure (cleartext) connections.
//
// The heuristic is very simple: the first byte of the stream is read. If it is
// ASCII SYN (22 dec, 0x16 hex), the connection is wrapped by crypto/tls.Server.
// Otherwise, the connection is returned as-is. The connection is always wrapped
// to allow re-reading of the first byte.
//
// This package will not work if the first byte a client sends can be ASCII SYN
// in the cleartext protocol.
package maybetls
