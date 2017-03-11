package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/enaeseth/asciist/convert"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	width = kingpin.Flag("width", "Character width to output").Default("80").Short('w').Uint()
	file  = kingpin.Arg("file", "Input image file").File()
)

func main() {
	var input *os.File

	kingpin.Parse()

	if *file != nil {
		input = *file
		defer input.Close()
	} else {
		input = os.Stdin
	}

	img, _, err := image.Decode(input)
	if err != nil {
		log.Fatal(err)
	}

	art := convert.FromImage(img, *width)
	fmt.Println(art)
}
