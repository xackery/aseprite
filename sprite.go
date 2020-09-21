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
	Width            uint16
	Height           uint16
	depth            uint16
	ncolors          uint16
	speed            uint16
	transparentIndex uint8
	colorSpace       int
	pixelRatio       float32
	gridBounds       image.Rectangle
	Tags             []*Tag
	slices           []*slice
	coreLayers       []*Layer
	Layers           map[string]*Layer
}
