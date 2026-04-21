package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.D7
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// led1 := machine.D7
	// led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		// led1.High()
		time.Sleep(time.Millisecond * 1000)

		led.High()
		time.Sleep(time.Millisecond * 1000)
		// led1.Low()
	}
}
