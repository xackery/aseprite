package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
)

func readColorChunk(f *os.File) (*palette, error) {
	var err error
	var packetCount int16
	err = binary.Read(f, binary.LittleEndian, &packetCount)
	if err != nil {
		return nil, fmt.Errorf("packetCount: %w", err)
	}
	skip := int8(0)

	p := new(palette)
	for i := int16(0); i < packetCount; i++ {
		var skipAdd int8
		err = binary.Read(f, binary.LittleEndian, &skipAdd)
		if err != nil {
			return nil, fmt.Errorf("skipAdd: %w", err)
		}
		skip += skipAdd

		var sizeBuf int8
		err = binary.Read(f, binary.LittleEndian, &sizeBuf)
		if err != nil {
			return nil, fmt.Errorf("sizeBuf %d: %w", i, err)
		}
		var size int
		if sizeBuf == 0 {
			size = 256
		} else {
			size = int(sizeBuf)
		}
		for c := int(skip); c < int(skip)+size; c++ {
			var r uint8
			err = binary.Read(f, binary.LittleEndian, &r)
			if err != nil {
				return nil, fmt.Errorf("r %d: %w", i, err)
			}
			var g uint8
			err = binary.Read(f, binary.LittleEndian, &g)
			if err != nil {
				return nil, fmt.Errorf("g %d: %w", i, err)
			}
			var b uint8
			err = binary.Read(f, binary.LittleEndian, &b)
			if err != nil {
				return nil, fmt.Errorf("b %d: %w", i, err)
			}
			p.colors = append(p.colors, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return p, nil
}
