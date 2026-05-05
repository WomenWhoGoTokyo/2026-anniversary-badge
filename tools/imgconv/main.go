// PNG → RGB565 変換
// usage:
//   go run ./tools/imgconv -in src/images/badge240.png -out src/images/badge.rgb565
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"os"
)

func main() {
	in := flag.String("in", "", "input PNG path")
	out := flag.String("out", "", "output raw RGB565 path")
	flag.Parse()

	if *in == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(*in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	buf := make([]byte, 0, w*h*2)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			r5 := (r >> 11) & 0x1F
			g6 := (g >> 10) & 0x3F
			b5 := (bl >> 11) & 0x1F
			v := uint16(r5<<11 | g6<<5 | b5)
			buf = append(buf, byte(v>>8), byte(v))
		}
	}

	if err := os.WriteFile(*out, buf, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s: %d bytes (%dx%d)\n", *out, len(buf), w, h)
}
