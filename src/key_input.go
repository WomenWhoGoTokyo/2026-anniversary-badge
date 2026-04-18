package main

import (
	"machine"
	"time"
)

func main() {
	machine.GPIO0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	for {
		if !machine.GPIO0.Get() {
			println("上側のボタンが押されました")
		} else {
		}

		if !machine.GPIO1.Get() {
			println("下側のボタンが押されました")
		} else {
		}

		time.Sleep(100 * time.Millisecond)
	}
}
