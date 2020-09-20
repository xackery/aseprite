package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"
)

type mask struct {
	name   string
	bounds image.Rectangle
	img    image.RGBA
}

func readMaskChunk(f *os.File) (*mask, error) {
	var err error
	m := &mask{}
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
	var w int16
	err = binary.Read(f, binary.LittleEndian, &w)
	if err != nil {
		return nil, fmt.Errorf("w: %w", err)
	}
	var h int16
	err = binary.Read(f, binary.LittleEndian, &h)
	if err != nil {
		return nil, fmt.Errorf("h: %w", err)
	}

	_, err = f.Seek(8, 1)
	if err != nil {
		return nil, fmt.Errorf("seek name: %w", err)
	}
	m.name, err = readString(f)
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}
	m.bounds = image.Rect(int(x), int(y), int(w), int(h))
	for v := int16(0); v < h; v++ {
		for u := int16(0); u < (w+7)/8; u++ {
			var bData int8
			err = binary.Read(f, binary.LittleEndian, &bData)
			if err != nil {
				return nil, fmt.Errorf("bData: %w", err)
			}
			for c := 0; c < 8; c++ {
				fmt.Println(bData & (1 << (7 - c)))
				//m.img.SetRGBA(u*8+c, v,byte & (1<<(7-c)))
			}
		}
	}
	return m, nil

}
