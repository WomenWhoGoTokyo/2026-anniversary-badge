package main

import (
	"machine"
	"time"
)

func main() {

	button1 := machine.D5
	button1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	for {
		if !button1.Get() {
			println("encorder sw is pressed!!")
		}

		time.Sleep(time.Millisecond * 100)
	}
}
