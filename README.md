GPIO
================

This is a Go library for asynchronous inputs on Linux devices with /sys/class/gpio. This implementation conforms to the [spec](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt). I have only tested it so far on Raspberry Pi but it should also work on similar systems like the Beaglebone. If you test this library on another system please tell me so that I can confirm it -- I'll give you credit here for the confirmation.

Watcher
---------------

The Watcher is a type which listens on the GPIO pins you specify and then notifies you when the values of those pins change. It uses a `select()` call so that it does not need to actively poll, which saves CPU time and gives you better latencies from your inputs.

Here is an example of how to use the Watcher. Note that the GPIO numbers we want are as the CPU/kernel knows them, not as they may be marked on any external hardware headers.

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
