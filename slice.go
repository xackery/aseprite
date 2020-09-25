package aseprite

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
)

// Slice represents a slice
type Slice struct {
	name       string
	bounds     image.Rectangle
	center     image.Rectangle
	pivot      image.Point
	frameIndex uint16
	UserData   *UserData
}

type sliceKey struct {
}

func readSlicesChunk(f io.ReadSeeker, frameIndex uint16, s *Sprite) error {
	var err error
	var sliceCount int32
	err = binary.Read(f, binary.LittleEndian, &sliceCount)
	if err != nil {
		return fmt.Errorf("sliceCount: %w", err)
	}
	_, err = f.Seek(8, 1)
	if err != nil {
		return fmt.Errorf("seek sliceChunk: %w", err)
	}

	for i := int32(0); i < sliceCount; i++ {
		_, err = readSliceChunk(f, frameIndex, s)
		if err != nil {
			return fmt.Errorf("readSliceChunk %d: %w", i, err)
		}
	}
	return nil
}

func readSliceChunk(f io.ReadSeeker, frameIndex uint16, s *Sprite) (*Slice, error) {
	var err error
	sl := new(Slice)
	var keyCount int32
	err = binary.Read(f, binary.LittleEndian, &keyCount)
	if err != nil {
		return nil, fmt.Errorf("keyCount: %w", err)
	}
	var flags int32
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %w", err)
	}
	_, err = f.Seek(4, 1)
	if err != nil {
		return nil, fmt.Errorf("seek name: %w", err)
	}

	sl.name, err = readString(f)
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}

	for j := int32(0); j < keyCount; j++ {
		var x int32
		err = binary.Read(f, binary.LittleEndian, &x)
		if err != nil {
			return nil, fmt.Errorf("x: %w", err)
		}
		var y int32
		err = binary.Read(f, binary.LittleEndian, &y)
		if err != nil {
			return nil, fmt.Errorf("y: %w", err)
		}
		var w int32
		err = binary.Read(f, binary.LittleEndian, &w)
		if err != nil {
			return nil, fmt.Errorf("w: %w", err)
		}
		var h int32
		err = binary.Read(f, binary.LittleEndian, &h)
		if err != nil {
			return nil, fmt.Errorf("h: %w", err)
		}
		sl.bounds = image.Rect(int(x), int(y), int(w), int(h))
		if flags&1 == 1 { //ASE_SLICE_FLAG_HAS_CENTER_BOUNDS
			var x int32
			err = binary.Read(f, binary.LittleEndian, &x)
			if err != nil {
				return nil, fmt.Errorf("x: %w", err)
			}
			var y int32
			err = binary.Read(f, binary.LittleEndian, &y)
			if err != nil {
				return nil, fmt.Errorf("y: %w", err)
			}
			var w int32
			err = binary.Read(f, binary.LittleEndian, &w)
			if err != nil {
				return nil, fmt.Errorf("w: %w", err)
			}
			var h int32
			err = binary.Read(f, binary.LittleEndian, &h)
			if err != nil {
				return nil, fmt.Errorf("h: %w", err)
			}
			sl.center = image.Rect(int(x), int(y), int(w), int(h))
		}
		if flags&2 == 2 { //ASE_SLICE_FLAG_HAS_PIVOT_POINT
			var x int32
			err = binary.Read(f, binary.LittleEndian, &x)
			if err != nil {
				return nil, fmt.Errorf("x: %w", err)
			}
			var y int32
			err = binary.Read(f, binary.LittleEndian, &y)
			if err != nil {
				return nil, fmt.Errorf("y: %w", err)
			}
			sl.pivot = image.Pt(int(x), int(y))
		}
		sl.frameIndex = frameIndex
		s.slices = append(s.slices, sl)
	}
	return sl, nil
}
