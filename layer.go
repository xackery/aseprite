package aseprite

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Layer represents layers of a sprite
type Layer struct {
	isImage   bool
	BlendMode int16
	Name      string
	Opacity   int8
	Flags     int16
	parents   []*Layer
	layers    []*Layer
	Cells     []*Cell
	UserData  *UserData
}

const (
	// Special internal/undocumented alpha compositing and blend modes
	blendModeUnspecified int16 = -1
	blendModeSrc         int16 = -2
	blendModeMerge       int16 = -3
	blendModeNegBw       int16 = -4 // Negative Black & White
	blendModeRedTint     int16 = -5
	blendModeBlueTint    int16 = -6
	blendModeDstOver     int16 = -7
	// Aseprite (.ase files) blend modes
	blendModeNormal        int16 = 0
	blendModeMultiply      int16 = 1
	blendModeScreen        int16 = 2
	blendModeOverlay       int16 = 3
	blendModeDarken        int16 = 4
	blendModeLighten       int16 = 5
	blendModeColorDodge    int16 = 6
	blendModeColorBurn     int16 = 7
	blendModeHardLight     int16 = 8
	blendModeSoftLight     int16 = 9
	blendModeDifference    int16 = 10
	blendModeExclusion     int16 = 11
	blendModeHslHue        int16 = 12
	blendModeHslSaturation int16 = 13
	blendModeHslColor      int16 = 14
	blendModeHslLuminosity int16 = 15
	blendModeAddition      int16 = 16
	blendModeSubtract      int16 = 17
	blendModeDivide        int16 = 18
)

func readLayerChunk(f io.ReadSeeker, headerFlags uint32, prevLayer *Layer, currentLevel int16) (*Layer, error) {
	// log := log.New()
	var err error
	layer := &Layer{
		UserData: &UserData{},
	}

	var flags int16
	err = binary.Read(f, binary.LittleEndian, &flags)
	if err != nil {
		return nil, fmt.Errorf("flags: %w", err)
	}
	var layerType int16
	err = binary.Read(f, binary.LittleEndian, &layerType)
	if err != nil {
		return nil, fmt.Errorf("layerType: %w", err)
	}
	var childLevel int16
	err = binary.Read(f, binary.LittleEndian, &childLevel)
	if err != nil {
		return nil, fmt.Errorf("childLevel: %w", err)
	}
	_, err = f.Seek(4, 1)
	if err != nil {
		return nil, fmt.Errorf("seek blendMode: %w", err)
	}
	var blendMode int16
	err = binary.Read(f, binary.LittleEndian, &blendMode)
	if err != nil {
		return nil, fmt.Errorf("blendMode: %w", err)
	}
	var opacity int8
	err = binary.Read(f, binary.LittleEndian, &opacity)
	if err != nil {
		return nil, fmt.Errorf("opacity: %w", err)
	}
	_, err = f.Seek(3, 1)
	if err != nil {
		return nil, fmt.Errorf("seek name: %w", err)
	}

	name, err := readString(f)
	if err != nil {
		return nil, fmt.Errorf("name: %w", err)
	}

	switch layerType {
	case 0: //ASE_FILE_LAYER_IMAGE
		layer.isImage = true
		if flags&8 != 8 { //8 = background
			layer.BlendMode = blendMode
			if headerFlags&1 == 1 { //ASE_FILE_FLAG_LAYER_WITH_OPACITY
				layer.Opacity = opacity
			}
		}
	case 1: //ASE_FILE_LAYER_GROUP
	default:
		return nil, nil
	}

	layer.Flags = flags
	layer.Name = name
	if prevLayer != nil {
		if childLevel == currentLevel {
			prevLayer.parents = append(prevLayer.parents, layer)
		} else if childLevel > currentLevel {
			prevLayer.layers = append(prevLayer.layers, layer)
		} else { //if childLevel < currentLevel {
			layer.parents = append(layer.parents, prevLayer)
		}
	}
	currentLevel = childLevel
	// log.Debug().Msgf("layer: %v", layer)
	return layer, nil
}
