package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
)

type userData struct {
	text  string
	color color.RGBA
}

func (ud userData) set(val userData) {
	if len(val.text) > 0 {
		ud.text = val.text
	}
	if val.color.R != 0 &&
		val.color.G != 0 &&
		val.color.B != 0 {
		ud.color = val.color
	}
}

func readUserDataChunk(f io.ReadSeeker, s *Sprite) (userData, error) {
	var err error
	ud := userData{}

	var flags int32
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return ud, fmt.Errorf("flags: %w", err)
	}
	if flags&1 == 1 { //ASE_USER_DATA_FLAG_HAS_TEXT
		ud.text, err = readString(f)
		if err != nil {
			return ud, fmt.Errorf("text: %w", err)
		}
	}
	if flags&2 == 2 { //ASE_USER_DATA_FLAG_HAS_COLOR
		var r uint8
		err = binary.Read(f, binary.LittleEndian, &r)
		if err != nil {
			return ud, fmt.Errorf("r: %w", err)
		}
		var g uint8
		err = binary.Read(f, binary.LittleEndian, &g)
		if err != nil {
			return ud, fmt.Errorf("g: %w", err)
		}
		var b uint8
		err = binary.Read(f, binary.LittleEndian, &b)
		if err != nil {
			return ud, fmt.Errorf("b: %w", err)
		}
		var a uint8
		err = binary.Read(f, binary.LittleEndian, &a)
		if err != nil {
			return ud, fmt.Errorf("a: %w", err)
		}
		ud.color = color.RGBA{R: r, G: g, B: b, A: a}
	}
	return ud, nil
}
