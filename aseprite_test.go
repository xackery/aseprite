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
	log.Debug().Msgf("sprite: %v, layers: %v", s, s.rootLayer.layers)
	for lIndex, l := range s.rootLayer.layers {
		if !l.isImage {
			continue
		}
		for cIndex, c := range l.cels {
			f, err := os.Create(fmt.Sprintf("tmp/image%d-%d.png", lIndex, cIndex))
			if err != nil {
				t.Fatalf("create: %v", err)
			}
			defer f.Close()

			err = png.Encode(f, convertImage(c.img, s.width, s.height, c.positionX, c.positionY))
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
		}
	}
	for cTag, t := range s.tags {
		log.Debug().Msgf("tag %d: %v", cTag, t)
	}
}
