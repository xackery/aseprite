package aseprite

import "image"

const (
	pixelFormatNone = iota
	pixelFormatIMAGERGB
	pixelFormatIMAGEGRAYSCALE
	pixelFormatIMAGEINDEXED
	pixelFormatIMAGEBITMAP
)

// Sprite represents an aseprite sprite file
type Sprite struct {
	frameCount       uint16
	width            uint16
	height           uint16
	depth            uint16
	ncolors          uint16
	speed            uint16
	transparentIndex uint8
	colorSpace       int
	pixelRatio       float32
	gridBounds       image.Rectangle
	tags             []*tag
	slices           []*slice
	rootLayer        *layer
	layers           []*layer
}

func (s *Sprite) pixelFormat() {

}
