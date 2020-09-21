package aseprite

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
)

type tag struct {
	from               int16
	to                 int16
	name               string
	color              color.RGBA
	animationDirection int8
}

func readTagChunk(f *os.File, s *Sprite) error {
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
		t := new(tag)
		err = binary.Read(f, binary.LittleEndian, &t.from)
		if err != nil {
			return fmt.Errorf("from: %w", err)
		}
		err = binary.Read(f, binary.LittleEndian, &t.to)
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
		t.animationDirection = aniDir
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
		t.color = color.RGBA{R: r, G: g, B: b, A: 255}
		_, err = f.Seek(1, 1)
		if err != nil {
			return fmt.Errorf("seek name: %w", err)
		}
		t.name, err = readString(f)
		if err != nil {
			return fmt.Errorf("name: %w", err)
		}
		s.tags = append(s.tags, t)
	}
	return nil
}
