package main

import (
	"fmt"
	"time"

	"github.com/luismesas/goPi/MCP23S17"
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

	// buzz := pfd.OutputPins[2]
	green := pfd.OutputPins[3]
	red := pfd.OutputPins[4]
	// blink(green, red)
	for {
		if pfd.Switches[0].Value() != 0 {
			green.AllOn()
		} else {
			green.AllOff()
		}
		if pfd.Switches[1].Value() != 0 {
			red.AllOn()
		} else {
			red.AllOff()
		}
	}
}

func blink(green, red *MCP23S17.MCP23S17RegisterBit) {
	// blink time!!
	fmt.Println("Bilnking HID reader")
	for {
		green.Toggle()
		time.Sleep(time.Second)
		green.Toggle()
		red.Toggle()
		time.Sleep(time.Second)
		red.Toggle()
		time.Sleep(time.Second)
	}
}
