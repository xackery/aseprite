package main

import (
	"encoding/binary"
	"fmt"
	"io"
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

	lastLayer := s.rootLayer
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
	var lastCel *cel
	var lastSlice *slice
	var pal *palette
	layers := []*layer{}

	log.Debug().Msgf("processing %d chunks for frame %d", h.chunkCount, frameIndex)
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
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("seek default: %w", err)
		}
		switch chunkType {
		case 11, 4: //ASE_FILE_CHUNK_FLI_COLOR, ASE_FILE_CHUNK_FLI_COLOR2 legacy

			if isIgnoreOldColorChunks {
				log.Debug().Msgf("ignoreOldChunks enabled, skipping %d", chunkType)
				continue
			}
			log.Debug().Msgf("readColorChunk 0x%x", pos)
			pal, err = readColorChunk(f)
			if err != nil {
				return nil, fmt.Errorf("readColorChunk %d: %w", chunkIndex, err)
			}
			fmt.Println("colorChunk palette", pal)
		case 0x2019: //ASE_FILE_CHUNK_PALETTE
			log.Debug().Msgf("readPaletteChunk 0x%x", pos)
			pal, err := readPaletteChunk(f, frameIndex, flags)
			if err != nil {
				return nil, fmt.Errorf("readPalleteChunk %d: %w", chunkIndex, err)
			}
			fmt.Println("palette", pal)
		case 0x2004: //ASE_FILE_CHUNK_LAYER
			log.Debug().Msgf("readLayerChunk 0x%x", pos)
			layer, err := readLayerChunk(f, flags, prevLayer, currentLevel)
			if err != nil {
				return nil, fmt.Errorf("readLayerChunk %d: %w", chunkIndex, err)
			}
			layers = append(layers, layer)
			if layer != nil {
				s.rootLayer.layers = append(s.rootLayer.layers, layer)
				lastLayer = layer
				lastSlice = nil
				lastCel = nil
			}
		case 0x2005: //ASE_FILE_CHUNK_CEL
			log.Debug().Msgf("readCelChunk 0x%x", pos)
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
			if lastCel == nil {
				log.Debug().Msg("skipping readCelExtraChunk, no layer set")
				continue
			}
			log.Debug().Msgf("readCelExtraChunk 0x%x", pos)
			err = readCelExtraChunk(f, lastCel)
			if err != nil {
				return nil, fmt.Errorf("readCelExtraChunk %d: %w", chunkIndex, err)
			}
		case 0x2007: //ASE_FILE_CHUNK_COLOR_PROFILE
			log.Debug().Msgf("readColorProfile 0x%x", pos)
			err = readColorProfile(f, s)
			if err != nil {
				return nil, fmt.Errorf("readColorProfile %d: %w", chunkIndex, err)
			}
		case 0x2016: //ASE_FILE_CHUNK_MASK
			log.Debug().Msgf("readMaskChunk 0x%x", pos)
			mask, err := readMaskChunk(f)
			if err != nil {
				return nil, fmt.Errorf("readMaskChunk %d: %w", chunkIndex, err)
			}
			if mask != nil {
				fmt.Println(mask)
			}
		case 0x2017: //ASE_FILE_CHUNK_PATH
			log.Debug().Msgf("ignoring chunk path 0x%x", pos)
			//ignore
		case 0x2018: //ASE_FILE_CHUNK_TAGS
			log.Debug().Msgf("readTagsChunk 0x%x", pos)
			err = readTagsChunk(f, s)
			if err != nil {
				return nil, fmt.Errorf("readTagsChunk %d: %w", chunkIndex, err)
			}
		case 0x2021: //ASE_FILE_CHUNK_SLICES
			log.Debug().Msgf("readSlicesChunk 0x%x", pos)
			err = readSlicesChunk(f, frameIndex, s)
			if err != nil {
				return nil, fmt.Errorf("readSlicesChunk %d: %w", chunkIndex, err)
			}
		case 0x2022: //ASE_FILE_CHUNK_SLICE
			log.Debug().Msgf("readSliceChunk 0x%x", pos)
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
			log.Debug().Msgf("readUserDataChunk 0x%x", pos)
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
			log.Debug().Msgf("readFrameHeader: ignoring chunk tileset 0x%x", pos)
			//ignore
		default:
			log.Warn().Uint32("chunkSize", chunkSize).Msgf("readFrameHeader: unhandled chunk type %d at index %d 0x%x", chunkType, chunkIndex, pos)
		}
	}

	return h, nil
}
