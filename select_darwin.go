package gpio

import (
	"syscall"
)

func doSelect(nfd int, r *syscall.FdSet, w *syscall.FdSet, e *syscall.FdSet, timeout *syscall.Timeval) (changed bool, err error) {
	err = syscall.Select(nfd, r, w, e, timeout)
	if err != nil {
		return false, err
	}
	return true, nil
}
