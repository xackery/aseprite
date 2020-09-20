package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"os"
)

func readCelExtraChunk(f *os.File, c *cel) error {
	var err error
	var flags int32
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("flags: %w", err)
	}
	if flags&1 == 1 { //ASE_CEL_EXTRA_FLAG_PRECISE_BOUNDS
		var x int32
		err = binary.Read(f, binary.LittleEndian, &x)
		if err != nil {
			return fmt.Errorf("x: %w", err)
		}
		var y int32
		err = binary.Read(f, binary.LittleEndian, &y)
		if err != nil {
			return fmt.Errorf("y: %w", err)
		}
		var w int32
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return fmt.Errorf("w: %w", err)
		}
		var h int32
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return fmt.Errorf("h: %w", err)
		}
		if w > 0 && h > 0 {
			c.boundsFixed = image.Rect(int(x), int(y), int(w), int(h))
		}
	}
	return nil
}
