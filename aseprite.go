package aseprite

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"

	"github.com/xackery/log"
)

// Load loads a sprite
func Load(path string) (*Sprite, error) {
	log := log.New()
	//path := "examples/_default.aseprite"
	log.Debug().Msgf("parsing %s", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s, err := decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return s, nil
}

func decode(f *os.File) (*Sprite, error) {
	log := log.New()
	isIgnoreOldColorChunks := false
	header, err := readHeader(f)
	if err != nil {
		return nil, fmt.Errorf("readHeader: %w", err)
	}
	log.Debug().Msgf("%s", header)
	if header.depth != 32 &&
		header.depth != 16 &&
		header.depth != 8 {
		return nil, fmt.Errorf("invalid color depth %d", header.depth)
	}

	if header.width < 1 || header.height < 1 {
		return nil, fmt.Errorf("invalid sprite site %dx%d", header.width, header.height)
	}

	s := &Sprite{
		width:            header.width,
		height:           header.height,
		ncolors:          header.ncolors,
		depth:            header.depth,
		frameCount:       header.frameCount,
		speed:            header.speed,
		transparentIndex: header.transparentIndex,
		pixelRatio:       float32(header.pixelWidth / header.pixelHeight),
		gridBounds:       image.Rect(int(header.gridX), int(header.gridY), int(header.gridWidth), int(header.gridHeight)),
		rootLayer:        &layer{},
		layers:           []*layer{},
	}
	for frameIndex := uint16(0); frameIndex < header.frameCount; frameIndex++ {
		err := readFrameHeader(f, frameIndex, header.flags, isIgnoreOldColorChunks, s)
		if err != nil {
			return nil, fmt.Errorf("readFrameHeader %d: %w", frameIndex, err)
		}
	}

	return s, nil
}

func readString(f *os.File) (string, error) {
	value := ""
	var length int16
	err := binary.Read(f, binary.LittleEndian, &length)
	if err != nil {
		return "", fmt.Errorf("length: %w", err)
	}
	if length == -1 {
		return "", nil
	}
	var buf uint8
	for i := 0; i < int(length); i++ {
		err := binary.Read(f, binary.LittleEndian, &buf)
		if err != nil {
			return "", fmt.Errorf("pos %d: %w", i, err)
		}
		value += string(buf)
	}

	return value, nil
}
