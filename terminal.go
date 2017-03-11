package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/mattn/go-isatty"
)

const fallbackWidth = 80

func defaultWidth() uint {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		return terminalWidth()
	}

	return fallbackWidth
}

// http://stackoverflow.com/a/16576712/3204569
func terminalWidth() uint {
	ws := winsize{}

	ret, _, _ := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)

	if int(ret) != 0 {
		return uint(fallbackWidth)
	}

	return uint(ws.Col)
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}
