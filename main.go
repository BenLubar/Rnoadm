package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var Seed int64

func main() {
	flag.Int64Var(&Seed, "seed", time.Now().UnixNano(), "the world seed (default: the number of nanoseconds since midnight UTC on 1970-01-01)")

	flag.Parse()

	go func() {
		for _ = range time.Tick(time.Minute) {
			EachLoadedZone(func(z *Zone) {
				z.Save()
			})
		}
	}()

	go func() {
		for _ = range time.Tick(time.Second / 4) {
			EachLoadedZone(func(z *Zone) {
				z.Think()
			})
		}
	}()

	for {
		log.Print(http.ListenAndServe(":2064", nil))
		time.Sleep(time.Second)
	}
}

type InteractHUD struct {
	Player       *Player
	TileX, TileY uint8
	Objects      []Object
	Offset       int
}

func (h *InteractHUD) Paint(setcell func(int, int, rune, Color)) {
	if h.Player.TileX != h.TileX || h.Player.TileY != h.TileY || h.Objects == nil {
		h.TileX, h.TileY = h.Player.TileX, h.Player.TileY
		minX := h.TileX - 1
		if minX == 255 {
			minX = 0
		}
		maxX := h.TileX + 1
		if maxX == 0 {
			maxX = 255
		}
		minY := h.TileY - 1
		if minY == 255 {
			minY = 0
		}
		maxY := h.TileY + 1
		if maxY == 0 {
			maxY = 255
		}
		z := GrabZone(h.Player.ZoneX, h.Player.ZoneY)
		z.Lock()
		var objects []Object
		for x := minX; x >= minX && x <= maxX; x++ {
			for y := minY; y >= minY && y <= maxY; y++ {
				objects = append(objects, z.Tile(x, y).Objects...)
			}
		}
		z.Unlock()
		ReleaseZone(z)
		h.Objects = objects
		h.Offset = 0
	}
	const keys = "12345678"
	for i, o := range h.Objects[h.Offset:] {
		if i >= len(keys) {
			break
		}
		setcell(0, i, rune(keys[i]), "#fff")
		setcell(1, i, ' ', "#fff")
		j := 1
		for _, r := range o.Name() {
			j++
			setcell(j, i, r, "#fff")
		}
	}
	if h.Offset > 0 {
		setcell(0, 8, '9', "#fff")
		setcell(1, 8, ' ', "#fff")
		j := 1
		for _, r := range "previous" {
			j++
			setcell(j, 8, r, "#fff")
		}
	}
	if len(h.Objects) > h.Offset+len(keys) {
		setcell(0, 9, '0', "#fff")
		setcell(1, 9, ' ', "#fff")
		j := 1
		for _, r := range "next" {
			j++
			setcell(j, 9, r, "#fff")
		}
	}
}

func (h *InteractHUD) Key(code int) bool {
	switch code {
	case '1', '2', '3', '4', '5', '6', '7', '8':
		// TODO
		return true
	case '9':
		if h.Offset > 0 {
			h.Offset--
			h.Player.Repaint()
		}
		return true
	case '0':
		if h.Offset+8 < len(h.Objects) {
			h.Offset++
			h.Player.Repaint()
		}
		return true

	case 27: // esc
		h.Player.hud = nil
		h.Player.Repaint()
		return true
	}
	return false
}
