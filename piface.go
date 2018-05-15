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

	buzz, green, red := pfd.OutputPins[2], pfd.OutputPins[3], pfd.OutputPins[4]

	// blink time!!
	fmt.Println("Bilnking HID reader")
	for {
		green.Toggle()
		time.Sleep(time.Second)
		green.Toggle()
		red.Toggle()
		buzz.Toggle()
		time.Sleep(time.Second)
		red.Toggle()
		buzz.Toggle()
	}
}
