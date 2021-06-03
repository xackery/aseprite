package aseprite

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

// Cell represents an image
type Cell struct {
	PositionX   int16
	PositionY   int16
	Opacity     int8
	Image       *image.RGBA
	frameIndex  uint16
	boundsFixed image.Rectangle
	Duration    uint16
	UserData    *UserData
	EbitenImage *ebiten.Image
}

func readCellChunk(f io.ReadSeeker, layers []*Layer, frameIndex uint16, chunkSize uint32, pal *palette, duration uint16) (*Cell, error) {
	// log := log.New()
	var err error
	c := new(Cell)
	c.UserData = &UserData{}
	c.Duration = duration
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
			return nil, fmt.Errorf("raw_cell w: %w", err)
		}
		var h int16
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("raw_cell h: %w", err)
		}

		if w > 0 && h > 0 {
			img, err = readRawImage(f, pixelFormatIMAGERGB, w, h, pal)
			if err != nil {
				return nil, fmt.Errorf("raw_cell readImage: %w", err)
			}
		}
		c.PositionX = x
		c.PositionY = y
		c.frameIndex = frameIndex
		c.Opacity = opacity
		c.Image = img
	case 1: //ASE_FILE_LINK_CEL
		// log.Debug().Msg("link cell")
		var linkFrame int16
		err = binary.Read(f, binary.LittleEndian, &linkFrame)
		if err != nil {
			return nil, fmt.Errorf("link_cell linkFrame: %w", err)
		}
		if len(layer.Cells) <= int(linkFrame) {
			return nil, fmt.Errorf("link_cell linkFrame %d out of bounds (%d)", linkFrame, len(layer.Cells))
		}

		link := layer.Cells[int(linkFrame)]
		c.PositionX = link.PositionX
		c.PositionY = link.PositionY
		c.Image = link.Image
		c.Opacity = link.Opacity
		c.frameIndex = frameIndex
		fmt.Println("link", c)
	case 2: //ASE_FILE_COMPRESSED_CEL
		// log.Debug().Msg("compressed cell")
		var w int16
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return nil, fmt.Errorf("compressed_cell w: %w", err)
		}
		var h int16
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("compressed_cell h: %w", err)
		}
		if w <= 0 || h <= 0 {
			return nil, fmt.Errorf("compressed_cell %dx%d is invalid", w, h)
		}

		img, err = readCompressedImage(f, pixelFormatIMAGERGB, w, h, chunkSize, pal)
		if err != nil {
			return nil, fmt.Errorf("raw_cell readImage: %w", err)
		}
		c.PositionX = x
		c.PositionY = y
		c.frameIndex = frameIndex
		c.Opacity = opacity
		c.Image = img
	case 3: //ASE_FILE_COMPRESSED_TILEMAP
		var w int16
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return nil, fmt.Errorf("compressed_cell w: %w", err)
		}
		var h int16
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("compressed_cell h: %w", err)
		}
		if w <= 0 || h <= 0 {
			return nil, fmt.Errorf("compressed_cell %dx%d is invalid", w, h)
		}
		var bitsPerTile int16 //at the moment it's always 32-bit per tile
		err = binary.Read(f, binary.LittleEndian, &bitsPerTile)
		if err != nil {
			return nil, fmt.Errorf("bitsPerTile: %w", err)
		}
		if bitsPerTile != 32 {
			return nil, fmt.Errorf("bitsPerTile expected 32, got %d", bitsPerTile)
		}
		var bitMaskTileID int32 //(e.g. 0x1fffffff for 32-bit tiles)
		err = binary.Read(f, binary.LittleEndian, &bitMaskTileID)
		if err != nil {
			return nil, fmt.Errorf("bitMaskTileID: %w", err)
		}
		var bitMaskXFlip int32
		err = binary.Read(f, binary.LittleEndian, &bitMaskXFlip)
		if err != nil {
			return nil, fmt.Errorf("bitMaskXFlip: %w", err)
		}
		var bitMaskYFlip int32
		err = binary.Read(f, binary.LittleEndian, &bitMaskYFlip)
		if err != nil {
			return nil, fmt.Errorf("bitMaskYFlip: %w", err)
		}
		var bitMask90CWRot int32
		err = binary.Read(f, binary.LittleEndian, &bitMask90CWRot)
		if err != nil {
			return nil, fmt.Errorf("bitMask90CWRot: %w", err)
		}
		//10 bytes reserved
		_, err = f.Seek(10, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("seek 10: %w", err)
		}
		log.Debug().Msgf("tile %dx%d bitsPerTile: %d", w, h, bitsPerTile)
		return nil, fmt.Errorf("wut")
		img, err = readCompressedImage(f, pixelFormatIMAGERGB, w, h, chunkSize, pal)
		if err != nil {
			return nil, fmt.Errorf("raw_cell readImage: %w", err)
		}
		c.PositionX = x
		c.PositionY = y
		c.frameIndex = frameIndex
		c.Opacity = opacity
		c.Image = img
	default:
		return nil, fmt.Errorf("unknown cellType %d", celType)
	}

	layer.Cells = append(layer.Cells, c)
	return c, nil
}
