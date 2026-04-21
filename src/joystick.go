package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	machine.InitADC()

	ax := machine.ADC{Pin: machine.GPIO29}
	ax.Configure(machine.ADCConfig{})
	ay := machine.ADC{Pin: machine.GPIO28}
	ay.Configure(machine.ADCConfig{})

	for {
		x := ax.Get()
		y := ay.Get()
		fmt.Printf("%04X %04X\n", x, y)
		time.Sleep(200 * time.Millisecond)
	}
}
