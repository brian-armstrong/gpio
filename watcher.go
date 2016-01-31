package gpio

import (
	"container/heap"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

type watcherAction int

const (
	watcherAdd watcherAction = iota
	watcherRemove
	watcherClose
)

type watcherCmd struct {
	pin    Pin
	action watcherAction
}

type watcherNotify struct {
	pin   Pin
	value uint
}

type fdHeap []uintptr

func (h fdHeap) Len() int { return len(h) }

// Less is actually greater (we want a max heap)
func (h fdHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h fdHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *fdHeap) Push(x interface{}) {
	*h = append(*h, x.(uintptr))
}

func (h *fdHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h fdHeap) FdSet() *syscall.FdSet {
	fdset := &syscall.FdSet{}
	for _, val := range h {
		fdset.Bits[val/64] |= 1 << uint(val) % 64
	}
	return fdset
}

const watcherCmdChanLen = 32
const notifyChanLen = 32

// Watcher provides asynchronous notifications on input changes
// The user should supply it pins to watch with AddPin and then wait for changes with Watch
type Watcher struct {
	pins       map[uintptr]Pin
	fds        fdHeap
	cmdChan    chan watcherCmd
	notifyChan chan watcherNotify
}

// NewWatcher creates a new Watcher instance for asynchronous inputs
func NewWatcher() *Watcher {
	w := &Watcher{
		pins:       make(map[uintptr]Pin),
		fds:        fdHeap{},
		cmdChan:    make(chan watcherCmd, watcherCmdChanLen),
		notifyChan: make(chan watcherNotify, notifyChanLen),
	}
	heap.Init(&w.fds)
	go w.watch()
	return w
}

func (w *Watcher) notify(fdset *syscall.FdSet) {
	for _, fd := range w.fds {
		if (fdset.Bits[fd/64] & (1 << uint(fd) % 64)) != 0 {
			pin := w.pins[fd]
			val, err := pin.Read()
			if err != nil {
				if err == io.EOF {
					w.removeFd(fd)
					continue
				}
				fmt.Printf("failed to read pinfile, %s", err)
				os.Exit(1)
			}
			msg := watcherNotify{
				pin:   pin,
				value: val,
			}
			select {
			case w.notifyChan <- msg:
			default:
			}
		}
	}
}

func (w *Watcher) fdSelect() {
	timeval := &syscall.Timeval{
		Sec:  1,
		Usec: 0,
	}
	fdset := w.fds.FdSet()
	changed, err := doSelect(int(w.fds[0]+1), nil, nil, fdset, timeval)
	if err != nil {
		fmt.Printf("failed to call syscall.Select, %s", err)
		os.Exit(1)
	}
	if changed {
		w.notify(fdset)
	}
}

func (w *Watcher) addPin(p Pin) {
	fd := p.f.Fd()
	w.pins[fd] = p
	heap.Push(&w.fds, fd)
}

func (w *Watcher) removeFd(fd uintptr) {
	// heap operates on an array index, so search heap for fd
	for index, v := range w.fds {
		if v == fd {
			heap.Remove(&w.fds, index)
			break
		}
	}
	pin := w.pins[fd]
	pin.f.Close()
	delete(w.pins, fd)
}

// removePin is only a wrapper around removeFd
// it finds fd given pin and then calls removeFd
func (w *Watcher) removePin(p Pin) {
	// we don't index by pin, so go looking
	for fd, pin := range w.pins {
		if pin.Number == p.Number {
			// found pin
			w.removeFd(fd)
			return
		}
	}
}

func (w *Watcher) doCmd(cmd watcherCmd) (shouldContinue bool) {
	shouldContinue = true
	switch cmd.action {
	case watcherAdd:
		w.addPin(cmd.pin)
	case watcherRemove:
		w.removePin(cmd.pin)
	case watcherClose:
		shouldContinue = false
	}
	return shouldContinue
}

func (w *Watcher) recv() (shouldContinue bool) {
	for {
		select {
		case cmd := <-w.cmdChan:
			shouldContinue = w.doCmd(cmd)
			if !shouldContinue {
				return
			}
		default:
			shouldContinue = true
			return
		}
	}
}

func (w *Watcher) watch() {
	for {
		// first we do a syscall.select with timeout if we have any fds to check
		if len(w.fds) != 0 {
			w.fdSelect()
		} else {
			// so that we don't churn when the fdset is empty, sleep as if in select call
			time.Sleep(1 * time.Second)
		}
		if w.recv() == false {
			return
		}
	}
}

// AddPin adds a new pin to be watched for changes
// The pin provided should be the pin known by the kernel
func (w *Watcher) AddPin(p uint) {
	pin := NewInput(p)
	setEdgeTrigger(pin, edgeBoth)
	w.cmdChan <- watcherCmd{
		pin:    pin,
		action: watcherAdd,
	}
}

// RemovePin stops the watcher from watching the specified pin
func (w *Watcher) RemovePin(p uint) {
	pin := Pin{
		Number: p,
	}
	w.cmdChan <- watcherCmd{
		pin:    pin,
		action: watcherRemove,
	}
}

// Watch blocks until one change occurs on one of the watched pins
// It returns the pin which changed and its new value
// Because the Watcher is not perfectly realtime it may miss very high frequency changes
// If that happens, it's possible to see consecutive changes with the same value
// Also, if the input is connected to a mechanical switch, the user of this library must deal with debouncing
func (w *Watcher) Watch() (p uint, v uint) {
	notification := <-w.notifyChan
	return notification.pin.Number, notification.value
}

// Close stops the watcher and releases all resources
func (w *Watcher) Close() {
	w.cmdChan <- watcherCmd{
		pin:    Pin{},
		action: watcherClose,
	}
}
