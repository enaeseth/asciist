// Package fixture is an internal package providing sample images
// and expected rendering results.
// Verifying that asciist rendered an image "correctly" is subjective
// and requires human eyes, so we use a "snapshot testing" approach
// where the output of a judged-good run is saved and used for regression
// testing.
// (see https://facebook.github.io/jest/docs/snapshot-testing.html)
package fixture

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"

	// Enable image decoders:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// ReadFixture loads a fixture image bundled with this package. The width used
// in rendering is determined automatically from the length of the first line
// of the saved (known-good) output.
// This must be called from a working directory adjacent to the fixture directory.
func ReadFixture(imgFilename string) (img *bytes.Reader, width uint, art string) {
	basePath := filepath.Join("..", "fixture")
	imgPath := filepath.Join(basePath, imgFilename)
	artFilename := fmt.Sprintf("%s.txt", strings.TrimSuffix(imgFilename, filepath.Ext(imgFilename)))
	artPath := filepath.Join(basePath, artFilename)

	imgBytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		panic(err)
	}
	img = bytes.NewReader(imgBytes)

	artBytes, err := ioutil.ReadFile(artPath)
	if err != nil {
		panic(err)
	}
	art = string(artBytes[:len(artBytes)-1]) // trim \n on last line

	index := strings.IndexByte(art, '\n')
	if index < 0 {
		panic(fmt.Errorf("missing \n in %s", imgFilename))
	}
	width = uint(index)

	return
}

// LoadFixture loads a fixture image bundled with this package. The width used
// in rendering is determined automatically from the length of the first line
// of the saved (known-good) output.
// This must be called from a working directory adjacent to the fixture directory.
func LoadFixture(imgFilename string) (img image.Image, width uint, art string) {
	imgReader, width, art := ReadFixture(imgFilename)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		panic(err)
	}

	return
}
