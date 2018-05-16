package main

import (
	"fmt"
	"sync"
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

type Card struct {
	Number int64
	mu     sync.Mutex
	count  int
}

var (
	pfd    *piface.PiFaceDigital
	reader Reader
)

const packetGap = time.Millisecond * 200

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
	event := make(chan struct{})
	card := &Card{}
	go ReadD0(card, event)
	go ReadD1(card, event)

	t := time.Now()
	for {
		select {
		case <-event:
			t = time.Now()
		default:
			if time.Now().Sub(t) > packetGap {
				card.mu.Lock()
				fmt.Printf("%b, %[1]x, (%d)\n", card.Number, card.count)
				card.Number = 0
				card.count = 0
				card.mu.Unlock()
			}
		}
	}
}

func ReadD0(card *Card, event chan<- struct{}) {
	var prev, cur byte
	for {
		if cur = reader.D0.Value(); prev < cur {
			card.mu.Lock()
			card.Number = card.Number << 1
			card.count++
			card.mu.Unlock()
			event <- struct{}{}
		}
	}
}

func ReadD1(card *Card, event chan<- struct{}) {
	var prev, cur byte
	for {
		if cur = reader.D1.Value(); prev < cur {
			card.mu.Lock()
			card.Number = card.Number<<1 | 1
			card.count++
			card.mu.Unlock()
			event <- struct{}{}
		}
	}
}

func reportCard(cardCh chan Card) {
	for {
		card := <-cardCh
		fmt.Printf("%b, %[1]x\n", card.Number)
	}
}

func blinkGreen() {
	green := pfd.OutputPins[3]
	green.Toggle()
	time.Sleep(time.Millisecond * 500)
	green.Toggle()
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
