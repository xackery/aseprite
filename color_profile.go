package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func readColorProfile(f *os.File, s *sprite) error {
	var err error

	var profileType int16
	err = binary.Read(f, binary.LittleEndian, &profileType)
	if err != nil {
		return fmt.Errorf("profileType: %w", err)
	}
	var flags int16
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return fmt.Errorf("flags: %w", err)
	}
	var gamma int32
	err = binary.Read(f, binary.LittleEndian, &gamma)
	if err != nil {
		return fmt.Errorf("gamma: %w", err)
	}

	switch profileType {
	case 0: //ASE_FILE_NO_COLOR_PROFILE
		if flags&1 == 1 { //ASE_COLOR_PROFILE_FLAG_GAMMA
			return fmt.Errorf("color profiles not supported")
		}
	default:
		return fmt.Errorf("profileType %d not supported", profileType)
	}
	s.colorSpace = 0
	return nil
}
