package gpio

import (
	"syscall"
)

func doSelect(nfd int, r *syscall.FdSet, w *syscall.FdSet, e *syscall.FdSet, timeout *syscall.Timeval) (changed bool, err error) {
	n, err := syscall.Select(nfd, r, w, e, timeout)
	if err != nil {
		return false, err
	}
	if n != 0 {
		return true, nil
	}
	return false, nil
}
