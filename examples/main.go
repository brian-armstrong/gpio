package main

import (
    "fmt"

    "github.com/everactive/gpio"
)

func main() {
    watcher := gpio.NewWatcher()
    watcher.AddPin(22)
    watcher.AddPin(27)
    defer watcher.Close()

    led1, err := gpio.NewOutput(666, true, true)
    if err != nil {
        panic(err)
    }

    led1.High()

    led2, err := gpio.NewInput(777, true)
    if err != nil {
        panic(err)
    }

    val, err := led2.Read()
    if err != nil {
        panic(err)
    }

    fmt.Printf("value: %d\n", val)

    go func() {
        for {
            pin, value := watcher.Watch()
            fmt.Printf("read %d from gpio %d\n", value, pin)
        }
    }()
}
