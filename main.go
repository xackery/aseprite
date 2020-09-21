package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/xackery/log"
)

func main() {
	log := log.New()
	err := run()
	if err != nil {
		log.Error().Err(err).Msg("failed to run")
	}
}

func run() error {
	log := log.New()
	path := "examples/_default.aseprite"
	log.Debug().Msgf("parsing %s", path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = decode(f)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}

func decode(f *os.File) error {
	log := log.New()
	isIgnoreOldColorChunks := false
	header, err := readHeader(f)
	if err != nil {
		return fmt.Errorf("readHeader: %w", err)
	}
	log.Debug().Msgf("%s", header)
	if header.depth != 32 &&
		header.depth != 16 &&
		header.depth != 8 {
		return fmt.Errorf("invalid color depth %d", header.depth)
	}

	if header.width < 1 || header.height < 1 {
		return fmt.Errorf("invalid sprite site %dx%d", header.width, header.height)
	}

	s := &sprite{
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
			return fmt.Errorf("readFrameHeader %d: %w", frameIndex, err)
		}
	}
	fmt.Println("sprite", s, "layers", s.rootLayer.layers)
	for lIndex, l := range s.rootLayer.layers {
		fmt.Println("layer index", lIndex, l)
		if l.isImage {
			for cIndex, c := range l.cels {
				f, err := os.Create(fmt.Sprintf("image%d-%d.png", lIndex, cIndex))
				if err != nil {
					return fmt.Errorf("create: %w", err)
				}
				defer f.Close()

				err = png.Encode(f, convertImage(c.img, s.width, s.height, c.positionX, c.positionY))
				if err != nil {
					return fmt.Errorf("encode: %w", err)
				}
			}
		}
	}
	return nil
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
	var buf int8
	for i := 0; i < int(length); i++ {
		err := binary.Read(f, binary.LittleEndian, &buf)
		if err != nil {
			return "", fmt.Errorf("pos %d: %w", i, err)
		}
		value += string(buf)
	}

	return value, nil
}
