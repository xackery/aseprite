package aseprite

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"

	"github.com/xackery/log"
)

type cel struct {
	positionX   int16
	positionY   int16
	opacity     int8
	img         *image.RGBA
	frameIndex  uint16
	boundsFixed image.Rectangle
	userData    *userData
}

func readCelChunk(f *os.File, layers []*Layer, frameIndex uint16, chunkSize uint32, pal *palette) (*cel, error) {
	log := log.New()
	var err error
	c := new(cel)
	var layerIndex int16

	err = binary.Read(f, binary.LittleEndian, &layerIndex)
	if err != nil {
		return nil, fmt.Errorf("layerIndex: %w", err)
	}
	var x int16
	err = binary.Read(f, binary.LittleEndian, &x)
	if err != nil {
		return nil, fmt.Errorf("x: %w", err)
	}
	var y int16
	err = binary.Read(f, binary.LittleEndian, &y)
	if err != nil {
		return nil, fmt.Errorf("y: %w", err)
	}
	var opacity int8
	err = binary.Read(f, binary.LittleEndian, &opacity)
	if err != nil {
		return nil, fmt.Errorf("opacity: %w", err)
	}
	var celType int16
	err = binary.Read(f, binary.LittleEndian, &celType)
	if err != nil {
		return nil, fmt.Errorf("celType: %w", err)
	}
	_, err = f.Seek(7, 1)
	if err != nil {
		return nil, fmt.Errorf("seek celType: %w", err)
	}

	if layerIndex < 0 {
		return nil, fmt.Errorf("invalid layer index %d", layerIndex)
	}
	if len(layers) <= int(layerIndex) {
		return nil, fmt.Errorf("layerIndex %d out of bound of layers (%d)", layerIndex, len(layers))
	}
	layer := layers[int(layerIndex)]
	if !layer.isImage {
		return nil, fmt.Errorf("layer %d does not contain image", layerIndex)
	}
	var img *image.RGBA
	switch celType {
	case 0: //ASE_FILE_RAW_CEL
		var w int16
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return nil, fmt.Errorf("raw_cel w: %w", err)
		}
		var h int16
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("raw_cel h: %w", err)
		}

		if w > 0 && h > 0 {
			img, err = readRawImage(f, pixelFormatIMAGERGB, w, h, pal)
			if err != nil {
				return nil, fmt.Errorf("raw_cel readImage: %w", err)
			}
		}
		c.positionX = x
		c.positionY = y
		c.frameIndex = frameIndex
		c.opacity = opacity
		c.img = img
	case 1: //ASE_FILE_LINK_CEL
		log.Debug().Msg("link cel")
		var linkFrame int16
		err = binary.Read(f, binary.LittleEndian, &linkFrame)
		if err != nil {
			return nil, fmt.Errorf("link_cel linkFrame: %w", err)
		}
		if len(layer.cels) <= int(linkFrame) {
			return nil, fmt.Errorf("link_cel linkFrame %d out of bounds (%d)", linkFrame, len(layer.cels))
		}

		link := layer.cels[int(linkFrame)]
		c.positionX = link.positionX
		c.positionY = link.positionY
		c.img = link.img
		c.opacity = link.opacity
		c.frameIndex = frameIndex
		fmt.Println("link", c)
	case 2: //ASE_FILE_COMPRESSED_CEL
		log.Debug().Msg("compressed cel")
		var w int16
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return nil, fmt.Errorf("compressed_cel w: %w", err)
		}
		var h int16
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("compressed_cel h: %w", err)
		}
		if w <= 0 || h <= 0 {
			return nil, fmt.Errorf("compressed_cel %dx%d is invalid", w, h)
		}

		img, err = readCompressedImage(f, pixelFormatIMAGERGB, w, h, chunkSize, pal)
		if err != nil {
			return nil, fmt.Errorf("raw_cel readImage: %w", err)
		}
		c.positionX = x
		c.positionY = y
		c.frameIndex = frameIndex
		c.opacity = opacity
		c.img = img
	default:
		return nil, fmt.Errorf("unknown celType %d", celType)
	}

	layer.cels = append(layer.cels, c)
	return c, nil
}
