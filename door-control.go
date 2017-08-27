/*



*/

package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var (
	// Use mcu pin 17, corresponds to physical pin 11 on the pi
	pinOut = rpio.Pin(17)
	pinInButton = rpio.Pin(26)
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode

	pinOut.Output()
	pinOut.Low()

        pinInButton.PullUp()
	pinInButton.Input()

	for true {
		time.Sleep(time.Second / 20)
		buttonState:=pinInButton.Read();
	
		if buttonState == rpio.Low {
			pinOut.High();
		} else {
			pinOut.Low();
		}
	}
}


