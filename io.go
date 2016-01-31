package gpio

import (
	"errors"
	"os"
	"time"
)

// Pin represents a single pin, which can be used either for reading or writing
type Pin struct {
	Number    uint
	direction direction
	f         *os.File
}

// NewInput opens the given pin number for reading. The number provided should be the pin number known by the kernel
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

// NewOutput opens the given pin number for writing. The number provided should be the pin number known by the kernel
// NewOutput also needs to know whether the pin should be initialized high (true) or low (false)
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

// Close releases the resources related to Pin
func (p Pin) Close() {
	p.f.Close()
}

// Read returns the value read at the pin as reported by the kernel. This should only be used for input pins
func (p Pin) Read() (value uint, err error) {
	if p.direction != inDirection {
		return 0, errors.New("pin is not configured for input")
	}
	return readPin(p)
}

// High sets the value of an output pin to logic high
func (p Pin) High() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 1)
}

// Low sets the value of an output pin to logic low
func (p Pin) Low() error {
	if p.direction != outDirection {
		return errors.New("pin is not configured for output")
	}
	return writePin(p, 0)
}
