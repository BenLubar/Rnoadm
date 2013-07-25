package main

import (
	"encoding/base32"
	"encoding/binary"
)

// cosine of 30 degrees (3^0.5 / 2)
const root3_2 = 0.8660254037844386467637231707529361834714026269051903

func zoneFilename(x, y int64) string {
	var buf [binary.MaxVarintLen64 * 2]byte
	i := binary.PutVarint(buf[:], x)
	i += binary.PutVarint(buf[i:], y)
	encoded := base32.StdEncoding.EncodeToString(buf[:i])

	l := len(encoded)
	for encoded[l-1] == '=' {
		l--
	}
	return "z" + encoded[:l] + ".gz"
}
