package convert

import (
	"github.com/enaeseth/asciist/fixture"
	"testing"
)

func TestConvert(t *testing.T) {
	filenames := []string{"diag-ramp.gif", "bmo.png", "forest.jpg"}

	for _, f := range filenames {
		img, width, expectedArt := fixture.LoadFixture(f)
		actualArt := FromImage(img, width).String()

		if actualArt != expectedArt {
			t.Errorf("%s: unexpected art:\n%s", f, actualArt)
		}
	}
}
