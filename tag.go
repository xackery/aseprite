package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
)

// Tag represents animation groupings
type Tag struct {
	From               int16
	To                 int16
	Name               string
	Color              color.RGBA
	AnimationDirection int8
}

func readTagChunk(f io.ReadSeeker, s *Sprite) error {
	var err error
	var tagCount int16

	err = binary.Read(f, binary.LittleEndian, &tagCount)
	if err != nil {
		return fmt.Errorf("tagCount: %w", err)
	}

	_, err = f.Seek(8, 1)
	if err != nil {
		return fmt.Errorf("seek tags: %w", err)
	}
	for c := int16(0); c < tagCount; c++ {
		t := new(Tag)
		err = binary.Read(f, binary.LittleEndian, &t.From)
		if err != nil {
			return fmt.Errorf("from: %w", err)
		}
		err = binary.Read(f, binary.LittleEndian, &t.To)
		if err != nil {
			return fmt.Errorf("to: %w", err)
		}
		var aniDir int8
		err = binary.Read(f, binary.LittleEndian, &aniDir)
		if err != nil {
			return fmt.Errorf("aniDir: %w", err)
		}
		if aniDir != 0 && //forward
			aniDir != 1 && //reverse
			aniDir != 2 { //ping pong
			aniDir = 0
		}
		t.AnimationDirection = aniDir
		_, err = f.Seek(8, 1)
		if err != nil {
			return fmt.Errorf("seek rgb: %w", err)
		}
		var r uint8
		err = binary.Read(f, binary.LittleEndian, &r)
		if err != nil {
			return fmt.Errorf("r: %w", err)
		}
		var g uint8
		err = binary.Read(f, binary.LittleEndian, &g)
		if err != nil {
			return fmt.Errorf("g: %w", err)
		}
		var b uint8
		err = binary.Read(f, binary.LittleEndian, &b)
		if err != nil {
			return fmt.Errorf("b: %w", err)
		}
		t.Color = color.RGBA{R: r, G: g, B: b, A: 255}
		_, err = f.Seek(1, 1)
		if err != nil {
			return fmt.Errorf("seek name: %w", err)
		}
		t.Name, err = readString(f)
		if err != nil {
			return fmt.Errorf("name: %w", err)
		}
		s.Tags = append(s.Tags, t)
	}
	return nil
}
