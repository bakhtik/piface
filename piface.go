package main

import (
	"fmt"
	"time"

	"github.com/luismesas/goPi/MCP23S17"
	"github.com/luismesas/goPi/piface"
	"github.com/luismesas/goPi/spi"
)

var pfd *piface.PiFaceDigital

func init() {
	// creates a new pifacedigital instance
	pfd = piface.NewPiFaceDigital(spi.DEFAULT_HARDWARE_ADDR, spi.DEFAULT_BUS, spi.DEFAULT_CHIP)
}

func main() {
	// initializes pifacedigital board
	err := pfd.InitBoard()
	if err != nil {
		fmt.Printf("Error on init board: %s", err)
		return
	}

	// buzz := pfd.OutputPins[2]
	// green := pfd.OutputPins[3]
	// red := pfd.OutputPins[4]
	// zero := pfd.InputPins[0]
	// one := pfd.InputPins[1]

	reader1, cardCh := make(chan int), make(chan int)
	go ReadD0(reader1)
	go ReadD1(reader1)
	// go Read(reader1)
	go reportCard(cardCh)
	card, count := 0, 0
	// t := time.Now()
	for {
		card = card<<1 | <-reader1
		count++
		if count == 35 {
			cardCh <- card
			count = 0
		}

		// digit := <-reader1
		// if count > 0 && time.Now().Sub(t) > time.Second {
		// 	cardCh <- card
		// 	card, count = 0, 0
		// }
		// t = time.Now()
		// count++
		// cad = card<<1 | digit
	}
}

// func Read(reader chan int) {
// 	D0 := pfd.InputPins[0]
// 	D1 := pfd.InputPins[1]
// 	for {
// 		if D0.Value() == 1 {
// 			reader <- 0
// 		}
// 		if D1.Value() == 1 {
// 			reader <- 1
// 		}
// 	}
// }

func ReadD0(reader chan int) {
	D0 := pfd.InputPins[0]
	for {
		if D0.Value() == 1 {
			reader <- 0
		}
	}
}

func ReadD1(reader chan int) {
	D1 := pfd.InputPins[1]
	for {
		if D1.Value() == 1 {
			reader <- 1
		}
	}
}

func reportCard(card chan int) {
	for {
		fmt.Printf("%b, %[1]x\n", <-card)
		blinkGreen()
	}
}

func blinkGreen() {
	green := pfd.OutputPins[3]
	green.Toggle()
	time.Sleep(time.Millisecond * 500)
	green.Toggle()
}

// func readCard(pfd *piface.PiFaceDigital, firstDigit int) {
// 	cardNumber := strconv.Itoa(firstDigit)
// 	u, t := time.Now(), time.Now()
// 	for t.Sub(u) < time.Millisecond*50 {
// 		u = t
// 		time.Sleep(time.Microsecond * 10)
// 		switch {
// 		case pfd.InputPins[0].Value() == 1:
// 			cardNumber += "0"
// 			t = time.Now()
// 		case pfd.InputPins[1].Value() == 1:
// 			cardNumber += "1"
// 			t = time.Now()
// 		}
// 	}
// 	fmt.Println(cardNumber)
// 	fmt.Println()
// 	time.Sleep(time.Millisecond * 100)
// }

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
