package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
)

type palette struct {
	colors []color.RGBA
}

func readPaletteChunk(f *os.File, frameIndex uint16, flags uint32) (*palette, error) {
	var err error
	var newSize int32
	err = binary.Read(f, binary.LittleEndian, &newSize)
	if err != nil {
		return nil, fmt.Errorf("newSize: %w", err)
	}

	var from int32
	err = binary.Read(f, binary.LittleEndian, &from)
	if err != nil {
		return nil, fmt.Errorf("from: %w", err)
	}

	var to int32
	err = binary.Read(f, binary.LittleEndian, &to)
	if err != nil {
		return nil, fmt.Errorf("to: %w", err)
	}
	_, err = f.Seek(8, 1)
	if err != nil {
		return nil, fmt.Errorf("seek pallette: %w", err)
	}
	p := new(palette)
	for i := from; i <= to; i++ {
		var flags int16
		err = binary.Read(f, binary.LittleEndian, &flags)
		if err != nil {
			return nil, fmt.Errorf("flags %d: %w", i, err)
		}
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
		var a uint8
		err = binary.Read(f, binary.LittleEndian, &a)
		if err != nil {
			return nil, fmt.Errorf("a %d: %w", i, err)
		}
		p.colors = append(p.colors, color.RGBA{R: r, G: g, B: b, A: a})
		if flags&1 == 1 { //ASE_PALETTE_FLAG_HAS_NAME
			var name string
			err = binary.Read(f, binary.LittleEndian, &name)
			if err != nil {
				return nil, fmt.Errorf("name %d: %w", i, err)
			}
		}
	}

	return p, nil
}
