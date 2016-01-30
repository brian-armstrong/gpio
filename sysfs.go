package gpio

import (
	"fmt"
	"os"
	"strconv"
)

type direction uint

const (
	inDirection direction = iota
	outDirection
)

type edge uint

const (
	edgeNone edge = iota
	edgeRising
	edgeFalling
	edgeBoth
)

func exportGPIO(p Pin) {
	export, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio export file for writing\n")
		os.Exit(1)
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
}

func unexportGPIO(p Pin) {
	export, err := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio unexport file for writing\n")
		os.Exit(1)
	}
	defer export.Close()
	export.Write([]byte(strconv.Itoa(int(p.Number))))
}

func setDirection(p Pin, d direction, initialValue uint) {
	dir, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Number), os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio %d direction file for writing\n", p.Number)
		os.Exit(1)
	}
	defer dir.Close()

	switch {
	case d == inDirection:
		dir.Write([]byte("in"))
	case d == outDirection && initialValue == 0:
		dir.Write([]byte("low"))
	case d == outDirection && initialValue == 1:
		dir.Write([]byte("high"))
	default:
		panic(fmt.Sprintf("setDirection called with invalid direction or initialValue, %d, %d", d, initialValue))
	}
}

func setEdgeTrigger(p Pin, e edge) {
	edge, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/edge", p.Number), os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("failed to open gpio %d edge file for writing\n", p.Number)
		os.Exit(1)
	}
	defer edge.Close()

	switch e {
	case edgeNone:
		edge.Write([]byte("none"))
	case edgeRising:
		edge.Write([]byte("rising"))
	case edgeFalling:
		edge.Write([]byte("falling"))
	case edgeBoth:
		edge.Write([]byte("both"))
	default:
		panic(fmt.Sprintf("setEdgeTrigger called with invalid edge %d", e))
	}
}

func openPin(p Pin) {
	f, err := os.Open(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Number))
	if err != nil {
		fmt.Printf("failed to open gpio %d value file for reading\n", p.Number)
		os.Exit(1)
	}
	p.f = f
}
