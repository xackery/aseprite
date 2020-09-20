package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/xackery/log"
)

type frameHeader struct {
	size       uint32
	magic      uint16
	chunkCount uint16
	duration   uint16
}

func readFrameHeader(f *os.File, frameIndex uint16, flags uint32, isIgnoreOldColorChunks bool, s *sprite) (*frameHeader, error) {

	log := log.New()
	var err error
	h := &frameHeader{}
	if f == nil {
		return nil, fmt.Errorf("file must not be nil")
	}
	err = binary.Read(f, binary.LittleEndian, &h.size)
	if err != nil {
		return nil, fmt.Errorf("size: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.magic)
	if err != nil {
		return nil, fmt.Errorf("magic: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.chunkCount)
	if err != nil {
		return nil, fmt.Errorf("chunkCount: %w", err)
	}
	err = binary.Read(f, binary.LittleEndian, &h.duration)
	if err != nil {
		return nil, fmt.Errorf("duration: %w", err)
	}
	_, err = f.Seek(2, 1)
	if err != nil {
		return nil, fmt.Errorf("nchunks offset: %w", err)
	}

	var nchunks uint32
	err = binary.Read(f, binary.LittleEndian, &nchunks)
	if err != nil {
		return nil, fmt.Errorf("nchunks: %w", err)
	}
	if h.chunkCount == 0xFFFF && h.chunkCount < uint16(nchunks) {
		h.chunkCount = uint16(nchunks)
	}

	/*if h.magic != 0xF1FA {
		_, err = f.Seek(int64(h.size), 1)
		if err != nil {
			return nil, fmt.Errorf("seek frame magic: %w", err)
		}
		return h, nil
	}*/

	var prevLayer *layer
	var currentLevel int16
	var lastLayer *layer
	var lastCel *cel
	var lastSlice *slice
	var pal *palette
	layers := []*layer{}

	for chunkIndex := uint16(0); chunkIndex < h.chunkCount; chunkIndex++ {
		var chunkSize uint32
		err = binary.Read(f, binary.LittleEndian, &chunkSize)
		if err != nil {
			return nil, fmt.Errorf("chunkSize %d: %w", chunkIndex, err)
		}
		var chunkType uint16
		err = binary.Read(f, binary.LittleEndian, &chunkType)
		if err != nil {
			return nil, fmt.Errorf("chunkType %d: %w", chunkIndex, err)
		}
		switch chunkType {
		case 11, 4: //ASE_FILE_CHUNK_FLI_COLOR, ASE_FILE_CHUNK_FLI_COLOR2 legacy
			if isIgnoreOldColorChunks {
				continue
			}
			pal, err = readColorChunk(f)
			if err != nil {
				return nil, fmt.Errorf("readColorChunk %d: %w", chunkIndex, err)
			}
			fmt.Println("colorChunk palette", pal)
		case 0x2019: //ASE_FILE_CHUNK_PALETTE
			pal, err := readPaletteChunk(f, frameIndex, flags)
			if err != nil {
				return nil, fmt.Errorf("readPalleteChunk %d: %w", chunkIndex, err)
			}
			fmt.Println("palette", pal)
		case 0x2004: //ASE_FILE_CHUNK_LAYER
			layer, err := readLayerChunk(f, flags, prevLayer, currentLevel)
			if err != nil {
				return nil, fmt.Errorf("readLayerChunk %d: %w", chunkIndex, err)
			}
			if layer != nil {
				layers = append(layers, layer)
				lastLayer = layer
				lastSlice = nil
				lastCel = nil
			}
		case 0x2005: //ASE_FILE_CHUNK_CEL
			cel, err := readCelChunk(f, layers, frameIndex, chunkSize, pal)
			if err != nil {
				return nil, fmt.Errorf("readCelChunk %d: %w", chunkIndex, err)
			}
			if cel != nil {
				lastCel = cel
				lastLayer = nil
				lastSlice = nil
			}
		case 0x2006: //ASE_FILE_CHUNK_CEL_EXTRA
			if lastCel != nil {
				err = readCelExtraChunk(f, lastCel)
				if err != nil {
					return nil, fmt.Errorf("readCelExtraChunk %d: %w", chunkIndex, err)
				}
			}
		case 0x2007: //ASE_FILE_CHUNK_COLOR_PROFILE
			err = readColorProfile(f, s)
			if err != nil {
				return nil, fmt.Errorf("readColorProfile %d: %w", chunkIndex, err)
			}
		case 0x2016: //ASE_FILE_CHUNK_MASK
			mask, err := readMaskChunk(f)
			if err != nil {
				return nil, fmt.Errorf("readMaskChunk %d: %w", chunkIndex, err)
			}
			if mask != nil {
				fmt.Println(mask)
			}
		case 0x2017: //ASE_FILE_CHUNK_PATH
			//ignore
		case 0x2018: //ASE_FILE_CHUNK_TAGS
			err = readTagsChunk(f, s)
			if err != nil {
				return nil, fmt.Errorf("readTagsChunk %d: %w", chunkIndex, err)
			}
		case 0x2021: //ASE_FILE_CHUNK_SLICES
			err = readSlicesChunk(f, frameIndex, s)
			if err != nil {
				return nil, fmt.Errorf("readSlicesChunk %d: %w", chunkIndex, err)
			}
		case 0x2022: //ASE_FILE_CHUNK_SLICE
			sl, err := readSliceChunk(f, frameIndex, s)
			if err != nil {
				return nil, fmt.Errorf("readSliceChunk %d: %w", chunkIndex, err)
			}
			if sl != nil {
				lastCel = nil
				lastLayer = nil
				lastSlice = sl
			}
		case 0x2020: //ASE_FILE_CHUNK_USER_DATA
			ud, err := readUserDataChunk(f, s)
			if err != nil {
				return nil, fmt.Errorf("readUserDataChunk %d: %w", chunkIndex, err)
			}
			if lastCel != nil {
				lastCel.userData.set(ud)
			}
			if lastLayer != nil {
				lastLayer.userData.set(ud)
			}
			if lastSlice != nil {
				lastSlice.userData.set(ud)
			}
		case 0x2023: //ASE_FILE_CHUNK_TILESET
			//ignore
		default:
			log.Warn().Msgf("unhandled chunk type %d at index %d", chunkType, chunkIndex)
		}
	}

	return h, nil
}
