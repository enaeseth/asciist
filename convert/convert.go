package convert

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

const characterAspectRatio = 1.0 / 2.0

// http://paulbourke.net/dataformats/asciiart/
var palette = []byte(" .:-=+*#%@")
var paletteSize = uint8(len(palette))

func FromImage(img image.Image, width uint) ASCIIArt {
	height := computeHeight(img.Bounds(), width)
	scaled := resize.Resize(width, height, img, resize.Lanczos3)

	return toASCII(scaled)
}

func computeHeight(bounds image.Rectangle, desiredWidth uint) uint {
	sourceWidth := bounds.Max.X - bounds.Min.X
	sourceHeight := bounds.Max.Y - bounds.Min.Y
	sourceAspectRatio := float64(sourceWidth) / float64(sourceHeight)

	return uint(float64(desiredWidth) / sourceAspectRatio * characterAspectRatio)
}

func toASCII(img image.Image) ASCIIArt {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	chars := make(ASCIIArt, 0, height)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]byte, 0, width)

		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			row = append(row, colorToByte(img.At(x, y)))
		}

		chars = append(chars, row)
	}

	return chars
}

func colorToByte(c color.Color) byte {
	gray := color.GrayModel.Convert(c).(color.Gray).Y

	var index uint8
	if gray == 255 {
		index = paletteSize - 1
	} else {
		index = uint8(float32(gray) / 255.0 * float32(paletteSize))
	}

	return palette[index]
}
