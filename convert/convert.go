// Package convert provides the core ASCII art conversion function.
package convert

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

// CharacterAspectRatio gives the percentage by which input images' heights
// will be scaled relative to their native aspect ratio.
// Characters in most fonts are taller than they are wide, and are
// even taller when the space between lines is taken into account,
// resulting in output ASCII "pixels" that are not square.
const CharacterAspectRatio = 1.0 / 2.0

// palette from http://paulbourke.net/dataformats/asciiart/
var palette = []byte(" .:-=+*#%@")
var paletteSize = uint8(len(palette))

// FromImage converts the given image to ASCII art, where each line
// is `width` characters wide, and the height of the output is chosen
// to preserve the input aspect ratio (taking into account that ASCII
// "pixels" are not square; see CharacterAspectRatio).
func FromImage(img image.Image, width uint) ASCIIArt {
	height := computeHeight(img.Bounds(), width)
	scaled := resize.Resize(width, height, img, resize.Lanczos3)

	return toASCII(scaled)
}

func computeHeight(bounds image.Rectangle, desiredWidth uint) uint {
	sourceWidth := bounds.Max.X - bounds.Min.X
	sourceHeight := bounds.Max.Y - bounds.Min.Y
	sourceAspectRatio := float64(sourceWidth) / float64(sourceHeight)

	return uint(float64(desiredWidth) / sourceAspectRatio * CharacterAspectRatio)
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

const maxLuminance = 255

func colorToByte(c color.Color) byte {
	luminance := color.GrayModel.Convert(c).(color.Gray).Y

	var index uint8
	if luminance == maxLuminance {
		index = paletteSize - 1
	} else {
		index = uint8(float32(luminance) / maxLuminance * float32(paletteSize))
	}

	return palette[index]
}
