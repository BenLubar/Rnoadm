package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

func main() {
	flagColor := flag.String("color", "FFFFFF", "base color (6 digit hex)")
	flagOutput := flag.String("o", "out.png", "output file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	c, err := strconv.ParseUint(*flagColor, 16, 24)
	if err != nil {
		flag.Usage()
		os.Exit(2)
	}
	ar := uint32(c >> 16 & 255)
	ag := uint32(c >> 8 & 255)
	ab := uint32(c & 255)

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(3)
	}
	base, err := png.Decode(f)
	f.Close()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(4)
	}
	bounds := base.Bounds()
	result := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, a := base.At(x, y).RGBA()
			result.Set(x, y, color.RGBA{
				R: fade(r, ar),
				G: fade(g, ag),
				B: fade(b, ab),
				A: uint8(a),
			})
		}
	}
	f, err = os.Create(*flagOutput)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(5)
	}
	defer f.Close()
	err = png.Encode(f, result)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(6)
	}
}

func fade(base, accent uint32) uint8 {
	if base > 127 {
		return 255 - fade(255-base, 255-accent)
	}
	return uint8(base * accent / 127)
}
