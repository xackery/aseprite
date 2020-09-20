package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"
)

func readRawImage(f *os.File, pixelFormat int, width int16, height int16, pal *palette) (*image.RGBA, error) {
	var err error
	img := new(image.RGBA)
	var r uint8
	var g uint8
	var b uint8
	var a uint8
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			switch pixelFormat {
			case pixelFormatIMAGERGB:
				err = binary.Read(f, binary.LittleEndian, &r)
				if err != nil {
					return nil, fmt.Errorf("r: %w", err)
				}
				err = binary.Read(f, binary.LittleEndian, &g)
				if err != nil {
					return nil, fmt.Errorf("g: %w", err)
				}
				err = binary.Read(f, binary.LittleEndian, &b)
				if err != nil {
					return nil, fmt.Errorf("b: %w", err)
				}
				err = binary.Read(f, binary.LittleEndian, &a)
				if err != nil {
					return nil, fmt.Errorf("a: %w", err)
				}
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
			case pixelFormatIMAGEGRAYSCALE:
				err = binary.Read(f, binary.LittleEndian, &r)
				if err != nil {
					return nil, fmt.Errorf("k: %w", err)
				}
				err = binary.Read(f, binary.LittleEndian, &a)
				if err != nil {
					return nil, fmt.Errorf("a: %w", err)
				}
				img.Set(x, y, color.RGBA{R: r, G: r, B: r, A: a})
			case pixelFormatIMAGEINDEXED:
				err = binary.Read(f, binary.LittleEndian, &r)
				if err != nil {
					return nil, fmt.Errorf("index: %w", err)
				}
				if int(r) >= len(pal.colors) {
					return nil, fmt.Errorf("index %d out of range for palette (%d)", r, len(pal.colors))
				}
				color := pal.colors[int(r)]
				img.Set(x, y, color)
			default:
				return nil, fmt.Errorf("unknown pixel format %d", pixelFormat)
			}
		}
	}
	return img, nil
}
