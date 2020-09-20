package main

const (
	pixelFormatNone = iota
	pixelFormatIMAGERGB
	pixelFormatIMAGEGRAYSCALE
	pixelFormatIMAGEINDEXED
	pixelFormatIMAGEBITMAP
)

type sprite struct {
	frameCount       uint16
	width            uint16
	height           uint16
	depth            uint16
	ncolors          uint16
	speed            uint16
	transparentIndex uint8
	colorSpace       int
	pixelRatio       float32
	tags             []*tag
	slices           []*slice
}

func (s *sprite) pixelFormat() {

}
