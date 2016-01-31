GPIO
================

This is a Go library for general purpose pins on Linux devices which support /sys/class/gpio. This implementation conforms to the [spec](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt). An advantage of using /sys/class/gpio is that we can receive interrupt-like notifications from the kernel when an input changes, rather than polling an input periodically. See the notes on Watcher for more info.

I have only tested it so far on Raspberry Pi but it should also work on similar systems like the Beaglebone. If you test this library on another system please tell me so that I can confirm it -- I'll give you credit here for the confirmation.

Note that the GPIO numbers we want here as the CPU/kernel knows them, not as they may be marked on any external hardware headers.

Input
---------------

Call `pin := gpio.NewInput(number)` to create a new input with the given pin numbering. You can then access the value of this pin with `pin.Read()`, which returns 0 when the pin's value is logic low and 1 when high.

If you are only concerned with when the pin's value changes, consider using `gpio.Watcher` instead.

Output
---------------

Call `pin := gpio.NewOutput(number, high)`, where `high` is a bool that describes the initial value of the pin -- set to false if you'd like the pin to initialize low, and true if you'd like it to initialize high.

Once you have a pin, you can change its value with `pin.Low()` and `pin.High()`.

Watcher
---------------

The Watcher is a type which listens on the GPIO pins you specify and then notifies you when the values of those pins change. It uses a `select()` call so that it does not need to actively poll, which saves CPU time and gives you better latencies from your inputs.

Here is an example of how to use the Watcher.

```
watcher := gpio.NewWatcher()
watcher.AddPin(22)
watcher.AddPin(27)
defer watcher.Close()

go func() {
    for {
        pin, value := watcher.Watch()
        fmt.Printf("read %d from gpio %d\n", value, pin)
    }
}()
```

This example would print once each time the value read on either pin 22 or 27 changes. It also prints each pin once when starting.

License
--------------
3-clause BSD
