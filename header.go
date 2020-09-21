package aseprite

import (
	"encoding/binary"
	"fmt"
	"os"
)

type header struct {
	size             uint32
	magic            uint16
	frameCount       uint16
	width            uint16
	height           uint16
	depth            uint16
	flags            uint32
	speed            uint16 // Deprecated, use "duration" of AsepriteFrameHeader
	next             uint32
	frit             uint32
	transparentIndex uint8
	ignore           [3]uint8 //3 bytes to ignore uint8
	ncolors          uint16
	pixelWidth       uint8
	pixelHeight      uint8
	gridX            int16
	gridY            int16
	gridWidth        int16
	gridHeight       int16
}

func (h *header) String() string {
	return fmt.Sprintf("header &{size: %d, frameCount: %d, dimensions: %dx%d, depth: %d, flags: %d, speed: %d, next: %d, frit: %d, transparentIndex: %d, ncolors: %d, pixelDimensions: %dx%d, gridPos: %dx%d, gridDimensions: %dx%d}", h.size, h.frameCount, h.width, h.height, h.depth, h.flags, h.speed, h.next, h.frit, h.transparentIndex, h.ncolors, h.pixelWidth, h.pixelHeight, h.gridX, h.gridY, h.gridWidth, h.gridHeight)
}

func readHeader(f *os.File) (*header, error) {
	var err error
	h := &header{}
	if f == nil {
		return nil, fmt.Errorf("file must not be nil")
	}
	err = binary.Read(f, binary.LittleEndian, &h.size)
	if err != nil {
		return nil, fmt.Errorf("size: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.magic)
	if err != nil {
		return nil, fmt.Errorf("magic: %w", err)
	}
	if h.magic != 0xA5E0 {
		return nil, fmt.Errorf("magic should be %x, got %x", 0xA5E0, h.magic)
	}
	err = binary.Read(f, binary.LittleEndian, &h.frameCount)
	if err != nil {
		return nil, fmt.Errorf("frameCount: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.width)
	if err != nil {
		return nil, fmt.Errorf("width: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.height)
	if err != nil {
		return nil, fmt.Errorf("height: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.depth)
	if err != nil {
		return nil, fmt.Errorf("depth: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.speed)
	if err != nil {
		return nil, fmt.Errorf("speed: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.next)
	if err != nil {
		return nil, fmt.Errorf("next: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.frit)
	if err != nil {
		return nil, fmt.Errorf("frit: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.transparentIndex)
	if err != nil {
		return nil, fmt.Errorf("transparentIndex: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.ignore)
	if err != nil {
		return nil, fmt.Errorf("ignore: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.ncolors)
	if err != nil {
		return nil, fmt.Errorf("ncolors: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.pixelWidth)
	if err != nil {
		return nil, fmt.Errorf("pixelWidth: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.pixelHeight)
	if err != nil {
		return nil, fmt.Errorf("pixelHeight: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.gridX)
	if err != nil {
		return nil, fmt.Errorf("gridX: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.gridY)
	if err != nil {
		return nil, fmt.Errorf("gridY: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.gridWidth)
	if err != nil {
		return nil, fmt.Errorf("gridWidth: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.gridHeight)
	if err != nil {
		return nil, fmt.Errorf("gridHeight: %w", err)
	}

	if h.ncolors == 0 {
		h.ncolors = 256
	}
	if h.pixelWidth == 0 || h.pixelHeight == 0 {
		h.pixelHeight = 1
		h.pixelWidth = 1
	}
	_, err = f.Seek(128, 0)
	if err != nil {
		return nil, fmt.Errorf("header offset: %w", err)
	}

	return h, nil
}
