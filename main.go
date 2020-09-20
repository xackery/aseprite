package main

import (
	"encoding/binary"
	"fmt"
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
	f, err := os.Open("examples/basic.aseprite")
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
	isIgnoreOldColorChunks := false
	header, err := readHeader(f)
	if err != nil {
		return fmt.Errorf("readHeader: %w", err)
	}
	fmt.Println("header", header)
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
	}

	for frameIndex := uint16(0); frameIndex < header.frameCount; frameIndex++ {
		fHeader, err := readFrameHeader(f, frameIndex, header.flags, isIgnoreOldColorChunks, s)
		if err != nil {
			return fmt.Errorf("readFrameHeader %d: %w", frameIndex, err)
		}
		fmt.Println("fheader", frameIndex, fHeader)
	}
	fmt.Println("sprite", s)
	return nil
}

func readString(f *os.File) (string, error) {
	value := ""
	for {
		var buf int8
		err := binary.Read(f, binary.LittleEndian, &buf)
		if err != nil {
			return "", fmt.Errorf("readString: %w", err)
		}
		if buf == 0 {
			break
		}
		value += string(buf)
	}
	return value, nil
}
