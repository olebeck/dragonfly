package world

import "image/color"

// Biome is a region in a world with distinct geographical features, flora, temperatures, humidity ratings,
// and sky, water, grass and foliage colours.
type Biome interface {
	// Temperature returns the temperature of the biome.
	Temperature() float64
	// Rainfall returns the rainfall of the biome.
	Rainfall() float64
	// Depth returns the depth of the biome.
	Depth() float64
	// Scale returns the scale of the biome.
	Scale() float64
	// WaterColour returns the water colour of the biome.
	WaterColour() color.RGBA
	// Tags returns the tags for the biome.
	Tags() []string
	// String returns the biome name as a string.
	String() string
	// EncodeBiome encodes the biome into an int value that is used to identify the biome over the network.
	EncodeBiome() int
}

type BiomeRegistry struct {
	IDToBiome   map[int]Biome
	NameToBiome map[string]Biome
}

var DefaultBiomes = &BiomeRegistry{
	IDToBiome:   make(map[int]Biome),
	NameToBiome: make(map[string]Biome),
}

func (br *BiomeRegistry) Clone() *BiomeRegistry {
	br2 := &BiomeRegistry{
		make(map[int]Biome),
		make(map[string]Biome),
	}
	for id, biome := range br.IDToBiome {
		br2.IDToBiome[id] = biome
	}
	for name, biome := range br.NameToBiome {
		br2.NameToBiome[name] = biome
	}
	return br2
}

func (br *BiomeRegistry) Register(b Biome) {
	id := b.EncodeBiome()
	if _, ok := br.IDToBiome[id]; ok {
		panic("cannot register the same biome (" + b.String() + ") twice")
	}
	br.IDToBiome[id] = b
	br.NameToBiome[b.String()] = b
}

// BiomeByID looks up a biome by the ID and returns it if found.
func (br *BiomeRegistry) BiomeByID(id int) (Biome, bool) {
	e, ok := br.IDToBiome[id]
	return e, ok
}

// BiomeByName looks up a biome by the name and returns it if found.
func (br *BiomeRegistry) BiomeByName(name string) (Biome, bool) {
	e, ok := br.NameToBiome[name]
	return e, ok
}

// Biomes returns a slice of all registered biomes.
func (br *BiomeRegistry) Biomes() []Biome {
	bs := make([]Biome, 0, len(br.IDToBiome))
	for _, b := range br.IDToBiome {
		bs = append(bs, b)
	}
	return bs
}

// ocean returns the ocean biome, or an unknown biome if no biome is registered.
func ocean() Biome {
	if o, ok := DefaultBiomes.BiomeByID(0); ok {
		return o
	}
	return unknownBiome{}
}

// unknownBiome is returned in place of a Biome that is not registered. It encodes back to the ID it was read from.
type unknownBiome struct {
	id int
}

func (unknownBiome) Temperature() float64    { return 0.5 }
func (unknownBiome) Rainfall() float64       { return 0 }
func (unknownBiome) Depth() float64          { return 0.1 }
func (unknownBiome) Scale() float64          { return 0.1 }
func (unknownBiome) WaterColour() color.RGBA { return color.RGBA{R: 0x44, G: 0xaf, B: 0xf5, A: 0xff} }
func (unknownBiome) Tags() []string          { return nil }
func (unknownBiome) String() string          { return "unknown" }
func (b unknownBiome) EncodeBiome() int      { return b.id }

func RegisterBiome(b Biome) {
	DefaultBiomes.Register(b)
}
