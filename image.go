package aseprite

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

func readRawImage(f io.ReadSeeker, pixelFormat int, width int16, height int16, pal *palette) (*image.RGBA, error) {
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

func readCompressedImage(f io.ReadSeeker, pixelFormat int, width int16, height int16, chunkSize uint32, pal *palette) (*image.RGBA, error) {
	// log := log.New().With().Int16("width", width).Int16("height", height).Logger()
	var err error
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	zr, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("zlib: %w", err)
	}
	defer zr.Close()

	data := make([]byte, chunkSize)

	buf := bytes.NewBuffer(data)

	// log.Debug().Msg("parsing")
	_, err = io.Copy(buf, zr)
	if err != nil {
		return nil, fmt.Errorf("copy: %w", err)
	}

	//16x12 = 192
	//raw 0's: 101
	//16x16 = 256
	//raw 0's: 79

	data = buf.Bytes()

	br := bytes.NewReader(data[chunkSize:])
	//buf = bytes.NewBuffer(data)
	//fmt.Println(fmt.Sprintf("%d %dx%d %x", pixelFormat, width, height, buf))

	var r uint8
	var g uint8
	var b uint8
	var a uint8

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			switch pixelFormat {
			case pixelFormatIMAGERGB:
				err = binary.Read(br, binary.LittleEndian, &r)
				if err != nil {
					return nil, fmt.Errorf("%dx%d, r: %w", x, y, err)
				}
				err = binary.Read(br, binary.LittleEndian, &g)
				if err != nil {
					return nil, fmt.Errorf("g: %w", err)
				}
				err = binary.Read(br, binary.LittleEndian, &b)
				if err != nil {
					return nil, fmt.Errorf("b: %w", err)
				}
				err = binary.Read(br, binary.LittleEndian, &a)
				if err != nil {
					return nil, fmt.Errorf("a: %w", err)
				}
				img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
			case pixelFormatIMAGEGRAYSCALE:
				err = binary.Read(br, binary.LittleEndian, &r)
				if err != nil {
					return nil, fmt.Errorf("k: %w", err)
				}
				err = binary.Read(br, binary.LittleEndian, &a)
				if err != nil {
					return nil, fmt.Errorf("a: %w", err)
				}
				img.Set(x, y, color.RGBA{R: r, G: r, B: r, A: a})
			case pixelFormatIMAGEINDEXED:
				err = binary.Read(br, binary.LittleEndian, &r)
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

func convertImage(src *image.RGBA, width uint16, height uint16, positionX int16, positionY int16) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for y := int16(0); y < int16(height); y++ {
		for x := int16(0); x < int16(width); x++ {
			img.SetRGBA(int(x+positionX), int(y+positionY), src.RGBAAt(int(x), int(y)))
		}
	}
	return img
}
