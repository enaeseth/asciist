package convert

import (
	"strings"
)

// ASCIIArt is the output of asciist's image conversion.
// The top-level slice contains the lines of the ASCII art (y axis),
// and each inner slice contains the bytes of each line (x axis).
type ASCIIArt [][]byte

// String builds a string of the art by converting each line to a string
// and joining them with newline characters.
func (img ASCIIArt) String() string {
	rows := make([]string, 0, len(img))

	for _, byteRow := range img {
		rows = append(rows, string(byteRow))
	}

	return strings.Join(rows, "\n")
}
