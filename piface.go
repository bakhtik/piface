package main

import (
	"fmt"
	"time"

	"github.com/luismesas/goPi/piface"
	"github.com/luismesas/goPi/spi"
)

func main() {

	// creates a new pifacedigital instance
	pfd := piface.NewPiFaceDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)

	// initializes pifacedigital board
	err := pfd.InitBoard()
	if err != nil {
		fmt.Printf("Error on init board: %s", err)
		return
	}

	// blink time!!
	fmt.Println("Bilnking HID reader")
	for {
		// pfd.Leds[5].Toggle()
		pfd.OutputPins[0].Toggle()
		pfd.OutputPins[2].Toggle()
		time.Sleep(time.Second)
	}
}
