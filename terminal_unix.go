// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import (
	"os"
	"unsafe"

	"github.com/mattn/go-isatty"
	"golang.org/x/sys/unix"
)

func defaultWidth() uint {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return terminalWidth()
	}

	return fallbackWidth
}

// http://stackoverflow.com/a/16576712/3204569
func terminalWidth() uint {
	ws := winsize{}

	ret, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(unix.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)

	if errno != 0 || int(ret) != 0 {
		return fallbackWidth
	}

	return uint(ws.Col)
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
