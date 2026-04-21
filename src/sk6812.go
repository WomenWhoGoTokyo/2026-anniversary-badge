package main

import (
	"machine"
	"time"

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

type WS2812B struct {
	Pin machine.Pin
	ws  *piolib.WS2812B
}

func NewWS2812B(pin machine.Pin) *WS2812B {
	s, _ := pio.PIO0.ClaimStateMachine()
	ws, _ := piolib.NewWS2812B(s, pin)
	ws.EnableDMA(true)
	return &WS2812B{
		ws: ws,
	}
}

func (ws *WS2812B) WriteRaw(rawGRB []uint32) error {
	return ws.ws.WriteRaw(rawGRB)
}

func main() {
	ws := NewWS2812B(machine.GPIO4)

	ws.WriteRaw([]uint32{0xFFFFFFFF, 0xFFFFFFFF})
	time.Sleep(500 * time.Millisecond)

	for {
		// green
		ws.WriteRaw([]uint32{0xFF0000FF, 0xFFFFFFFF})
		time.Sleep(500 * time.Millisecond)

		ws.WriteRaw([]uint32{0xFF0000FF, 0xFF0000FF})
		time.Sleep(500 * time.Millisecond)

		// red
		ws.WriteRaw([]uint32{0x00FF00FF, 0xFF0000FF})
		time.Sleep(500 * time.Millisecond)

		ws.WriteRaw([]uint32{0x00FF00FF, 0x00FF00FF})
		time.Sleep(500 * time.Millisecond)

		// blue
		ws.WriteRaw([]uint32{0x0000FFFF, 0x00FF00FF})
		time.Sleep(500 * time.Millisecond)

		ws.WriteRaw([]uint32{0x0000FFFF, 0x0000FFFF})
		time.Sleep(500 * time.Millisecond)

		// white
		ws.WriteRaw([]uint32{0xFFFFFFFF, 0x0000FFFF})
		time.Sleep(500 * time.Millisecond)

		ws.WriteRaw([]uint32{0xFFFFFFFF, 0xFFFFFFFF})
		time.Sleep(500 * time.Millisecond)
	}
}
