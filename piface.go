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

type Card struct {
	Number int64
	Count  int
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
		D0:    pfd.InputPins[0],
		D1:    pfd.InputPins[1],
	}
}

func main() {
	// card := Card{}
	// readerCh := make(chan int, 35)
	// event := make(chan struct{})
	// // cardCh := make(chan Card)
	// // go reportCard(cardCh)
	// go ReadD0(readerCh, event)
	// go ReadD1(readerCh, event)
	// t := time.Now()
	// count := 0
	// for {
	// 	select {
	// 	case <-event:
	// 		count++
	// 		t = time.Now()
	// 	default:
	// 		if count > 0 && time.Now().Sub(t) > packetGap {
	// 			for digit := range readerCh {
	// 				fmt.Print(digit)
	// 				count--
	// 				if count == 0 {
	// 					break
	// 				}
	// 			}
	// 			fmt.Println()
	// 		}
	// 	}

	// }
	// zch, och := make(chan int), make(chan int)
	// go func() {
	// 	var t time.Time
	// 	zeroes := 0
	// 	for {
	// 		if reader.D0.Value() == 1 {
	// 			t = time.Now()
	// 			zeroes++
	// 		}
	// 		if zeroes > 0 && time.Now().Sub(t) > packetGap {
	// 			zch <- zeroes
	// 			zeroes = 0
	// 		}
	// 	}
	// }()
	// go func() {
	// 	var t time.Time
	// 	ones := 0
	// 	for {
	// 		if reader.D1.Value() == 1 {
	// 			t = time.Now()
	// 			ones++
	// 		}
	// 		if ones > 0 && time.Now().Sub(t) > packetGap {
	// 			och <- ones
	// 			ones = 0
	// 		}
	// 	}
	// }()

	count := 0
	t := time.Now()
	var prevD0, curD0, prevD1, curD1 byte
	for {
		curD0, curD1 = reader.D0.Value(), reader.D1.Value()
		if prevD0 == 0 && curD0 == 1 {
			count++
			t = time.Now()
		}
		if prevD1 == 0 && curD1 == 1 {
			count++
			t = time.Now()
		}
		prevD0, prevD1 = curD0, curD1
		if count > 0 && time.Now().Sub(t) > time.Millisecond*100 {
			fmt.Println(count)
			count = 0
		}
	}

}

func ReadD0(readerCh chan int, event chan struct{}) {
	for {
		if reader.D0.Value() == 1 {
			readerCh <- 0
			event <- struct{}{}

		}
	}
}

func ReadD1(readerCh chan int, event chan struct{}) {
	for {
		if reader.D1.Value() == 1 {
			readerCh <- 1
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
