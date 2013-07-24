package main

import (
	"encoding/base32"
	"encoding/binary"
)

// cosine of 30 degrees (3^0.5 / 2)
const root3_2 = 0.8660254037844386467637231707529361834714026269051903

func zoneFilename(x, y int64) string {
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(x))
	binary.LittleEndian.PutUint64(buf[8:], uint64(y))
	return base32.StdEncoding.EncodeToString(buf[:])
}
