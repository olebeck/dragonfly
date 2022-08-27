package block

// MushroomBlockType represents a type of Mushroom block.
type MushroomBlockType struct {
	mushroomblock
}

func Brown() MushroomBlockType {
	return MushroomBlockType{0}
}

func Red() MushroomBlockType {
	return MushroomBlockType{1}
}

// Uint8 ...
func (f mushroomblock) Uint8() uint8 {
	return uint8(f)
}

// Name ...
func (f mushroomblock) Name() string {
	switch f {
	case 0:
		return "Brown Mushroom Block"
	case 1:
		return "Red Mushroom Block"
	}
	panic("unknown mushroomblock type")
}

// String ...
func (f mushroomblock) String() string {
	switch f {
	case 0:
		return "brown"
	case 1:
		return "red"
	}
	panic("unknown mushroomblock type")
}

type mushroomblock uint8
