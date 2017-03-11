package convert

import (
	"strings"
)

type ASCIIArt [][]byte

func (img ASCIIArt) String() string {
	rows := make([]string, 0, len(img))

	for _, byteRow := range img {
		rows = append(rows, string(byteRow))
	}

	return strings.Join(rows, "\n")
}
