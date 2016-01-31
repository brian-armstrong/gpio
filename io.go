package gpio

import (
	"errors"
	"os"
	"time"
)

type Pin struct {
	Number    uint
	direction direction
	f         *os.File
}

func NewInput(p uint) Pin {
	pin := Pin{
		Number: p,
	}
	exportGPIO(pin)
	time.Sleep(10 * time.Millisecond)
	pin.direction = inDirection
	setDirection(pin, inDirection, 0)
	pin = openPin(pin, false)
	return pin
}

func NewOutput(p uint, initHigh bool) Pin {
	pin := Pin{
		Number: p,
	}
	exportGPIO(pin)
	time.Sleep(10 * time.Millisecond)
	initVal := uint(0)
	if initHigh {
		initVal = uint(1)
	}
	pin.direction = outDirection
	setDirection(pin, outDirection, initVal)
	pin = openPin(pin, true)
	return pin
}

func (p Pin) Close() {
	p.f.Close()
}

func (p Pin) Read() (value uint, err error) {
	if p.direction != inDirection {
		return 0, errors.New("pin is not configured for input")
	}
	return readPin(p)
}

func (p Pin) High() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 1)
}

func (p Pin) Low() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 0)
}
