package block

// WeatheringType represents a weathering stage of a block.
type WeatheringType struct {
	weathering
}

// NotWeathered returns the not weathered variant
func NotWeathered() WeatheringType {
	return WeatheringType{0}
}

// Exposed returns the exposed variant
func Exposed() WeatheringType {
	return WeatheringType{1}
}

// Weathered returns the weathered variant
func Weathered() WeatheringType {
	return WeatheringType{2}
}

// Oxidized returns the oxidized variant
func Oxidized() WeatheringType {
	return WeatheringType{3}
}

// WeatheringTypes returns all weathering types.
func WeatheringTypes() []WeatheringType {
	return []WeatheringType{NotWeathered(), Exposed(), Weathered(), Oxidized()}
}

type weathering uint8

// Uint8 returns the coral as a uint8.
func (w weathering) Uint8() uint8 {
	return uint8(w)
}

// Name ...
func (w weathering) Name() string {
	switch w {
	case 0:
		return ""
	case 1:
		return "Exposed"
	case 2:
		return "Weathered"
	case 3:
		return "Oxidized"
	}
	panic("unknown weathering type")
}

// String ...
func (w weathering) String() string {
	switch w {
	case 0:
		return ""
	case 1:
		return "exposed"
	case 2:
		return "weathered"
	case 3:
		return "oxidized"
	}
	panic("unknown weathering type")
}
