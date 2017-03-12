// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package main

func defaultWidth() uint {
	return fallbackWidth
}
