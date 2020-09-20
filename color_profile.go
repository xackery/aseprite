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

	_, err = f.Seek(8, 1)
	if err != nil {
		return fmt.Errorf("colorProfile padding: %w", err)
	}
	switch profileType {
	case 0: //ASE_FILE_NO_COLOR_PROFILE
		if flags&1 == 1 { //ASE_COLOR_PROFILE_FLAG_GAMMA
			return fmt.Errorf("color profiles not supported")
		}
	case 1: //ASE_FILE_SRGB_COLOR_PROFILE
	case 2: //ASE_FILE_ICC_COLOR_PROFILE
		_, err = f.Seek(4, 1)
		if err != nil {
			return fmt.Errorf("seek colorProfile: %w", err)
		}
	default:
		return fmt.Errorf("profileType %d not supported", profileType)
	}
	s.colorSpace = 0
	return nil
}
