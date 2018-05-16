package main

import (
	"fmt"
	"time"

	"github.com/luismesas/goPi/MCP23S17"
	"github.com/luismesas/goPi/piface"
	"github.com/luismesas/goPi/spi"
)

type Reader struct {
	Buzz  *MCP23S17.MCP23S17RegisterBit
	Green *MCP23S17.MCP23S17RegisterBit
	Red   *MCP23S17.MCP23S17RegisterBit
	D0    *MCP23S17.MCP23S17RegisterBitNeg
	D1    *MCP23S17.MCP23S17RegisterBitNeg
}

type Toggler interface {
	Toggle()
}

var (
	pfd    *piface.PiFaceDigital
	reader Reader
)

const packetGap = time.Millisecond * 500

func init() {
	// creates a new pifacedigital instance
	pfd = piface.NewPiFaceDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
	// initializes pifacedigital board
	err := pfd.InitBoard()
	if err != nil {
		fmt.Printf("Error on init board: %s", err)
		return
	}
	reader = Reader{
		Buzz:  pfd.OutputPins[2],
		Green: pfd.OutputPins[3],
		Red:   pfd.OutputPins[4],
		D0:    pfd.InputPins[4],
		D1:    pfd.InputPins[5],
	}
}

func main() {
	go SwitchFunc(0, reader.Buzz)()
	go SwitchFunc(1, reader.Green)()
	go SwitchFunc(2, reader.Red)()

	select {}
}

func switch1() {
	var prev, cur byte
	prev = 1
	for {
		cur = pfd.Switches[0].Value()
		if prev != cur {
			reader.Buzz.Toggle()
		}
		prev = cur
	}
}

func switch2() {
	var prev, cur byte
	prev = 1
	for {
		cur = pfd.Switches[0].Value()
		if prev != cur {
			reader.Green.Toggle()
		}
		prev = cur
	}
}

func SwitchFunc(swithIndex int, device Toggler) func() {
	return func() {
		var prev, cur byte
		prev = 1
		for {
			cur = pfd.Switches[swithIndex].Value()
			if prev != cur {
				device.Toggle()
			}
			prev = cur
		}
	}
}
