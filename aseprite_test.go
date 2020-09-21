package aseprite

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/xackery/log"
)

func TestLoad(t *testing.T) {
	log := log.New()
	s, err := Load("examples/_default.aseprite")
	if err != nil {
		t.Fatalf("%v", err)
	}
	log.Debug().Msgf("sprite: %v, layers: %v", s, s.coreLayers)
	for lIndex, l := range s.coreLayers {
		if !l.isImage {
			continue
		}
		for cIndex, c := range l.Cells {
			f, err := os.Create(fmt.Sprintf("tmp/image%d-%d.png", lIndex, cIndex))
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			defer f.Close()

			err = png.Encode(f, convertImage(c.Image, s.Width, s.Height, c.PositionX, c.PositionY))
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
		}
	}
	for cTag, t := range s.Tags {
		log.Debug().Msgf("tag %d: %v", cTag, t)
	}
}
