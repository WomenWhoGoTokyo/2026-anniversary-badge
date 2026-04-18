package main

import (
	"machine"

	"tinygo.org/x/drivers/encoders"
)

func main() {
	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO2,
		machine.GPIO3,
	)

	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	for oldValue := 0; ; {
		if newValue := enc.Position(); newValue != oldValue {
			println("value: ", newValue)
			oldValue = newValue
		}
	}
}
