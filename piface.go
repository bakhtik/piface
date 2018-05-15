package main

import (
	"fmt"
	"strconv"
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
	// green := pfd.OutputPins[3]
	// red := pfd.OutputPins[4]
	zero := pfd.InputPins[0]
	one := pfd.InputPins[1]

	for {
		switch {
		case zero.Value() == 1:
			readCard(pfd, 0)
		case one.Value() == 1:
			readCard(pfd, 1)
		}
	}

}

func readCard(pfd *piface.PiFaceDigital, firstDigit int) {
	cardNumber := strconv.Itoa(firstDigit)
	u, t := time.Now(), time.Now()
	for t.Sub(u) < time.Millisecond*50 {
		u = t
		time.Sleep(time.Microsecond * 10)
		switch {
		case pfd.InputPins[0].Value() == 1:
			cardNumber += "0"
			t = time.Now()
		case pfd.InputPins[1].Value() == 1:
			cardNumber += "1"
			t = time.Now()
		}
	}
	fmt.Println(cardNumber)
	fmt.Println()
	time.Sleep(time.Millisecond * 100)
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

func switchLeds(pfd *piface.PiFaceDigital) {
	for {
		if pfd.Switches[0].Value() != 0 {
			pfd.OutputPins[3].AllOff()
		} else {
			pfd.OutputPins[3].AllOn()
		}
		if pfd.Switches[1].Value() != 0 {
			pfd.OutputPins[4].AllOff()
		} else {
			pfd.OutputPins[4].AllOn()
		}
	}
}
