package convert

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func TestConvert(t *testing.T) {
	fixtures := []string{"diag-ramp.gif", "bmo.png", "forest.jpg"}

	for _, fixture := range fixtures {
		img, width, expectedArt := loadFixture(fixture)
		actualArt := FromImage(img, width).String()

		if actualArt != expectedArt {
			t.Errorf("%s: unexpected art:\n%s", fixture, actualArt)
		}
	}
}

func loadFixture(imgFilename string) (img image.Image, width uint, art string) {
	basePath := filepath.Join("..", "fixture")
	imgPath := filepath.Join(basePath, imgFilename)
	artFilename := fmt.Sprintf("%s.txt", strings.TrimSuffix(imgFilename, filepath.Ext(imgFilename)))
	artPath := filepath.Join(basePath, artFilename)

	imgFile, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

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

	img, _, err = image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	return
}
