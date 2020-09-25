package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
)

// UserData represents user defined data
type UserData struct {
	Text  string
	Color color.RGBA
}

func (ud UserData) set(val UserData) {
	if len(val.Text) > 0 {
		ud.Text = val.Text
	}
	if val.Color.R != 0 &&
		val.Color.G != 0 &&
		val.Color.B != 0 {
		ud.Color = val.Color
	}
}

func readUserDataChunk(f io.ReadSeeker, s *Sprite) (UserData, error) {
	var err error
	ud := UserData{}

	var flags int32
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return ud, fmt.Errorf("flags: %w", err)
	}
	if flags&1 == 1 { //ASE_USER_DATA_FLAG_HAS_TEXT
		ud.Text, err = readString(f)
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
		ud.Color = color.RGBA{R: r, G: g, B: b, A: a}
	}
	return ud, nil
}
