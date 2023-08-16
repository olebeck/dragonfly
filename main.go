package main

import (
	"fmt"
	"os"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

// Flat is the flat generator of World. It generates flat worlds (like those in vanilla) with no other
// decoration. It may be constructed by calling NewFlat.
type Flat struct {
	// biome is the encoded biome that the generator should use.
	biome uint32
	// layers is a list of block runtime ID layers placed by the Flat generator. The layers are ordered in a way where
	// the last element in the slice is placed as the bottom-most block of the chunk.
	layers []uint32
	// n is the amount of layers in the slice above.
	n int16
}

// NewFlat creates a new Flat generator. Chunks generated are completely filled with the world.Biome passed. layers is a
// list of block layers placed by the Flat generator. The layers are ordered in a way where the last element in the
// slice is placed as the bottom-most block of the chunk.
func NewFlat(biome world.Biome, layers []world.Block) Flat {
	f := Flat{
		biome:  uint32(biome.EncodeBiome()),
		layers: make([]uint32, len(layers)),
		n:      int16(len(layers)),
	}
	for i, b := range layers {
		f.layers[i] = world.BlockRuntimeID(b)
	}
	return f
}

// GenerateChunk ...
func (f Flat) GenerateChunk(cp world.ChunkPos, chunk *chunk.Chunk) {
	if !(cp.X() == 0 && cp.Z() == 0) {
		return
	}
	min, max := int16(chunk.Range().Min()), int16(chunk.Range().Max())

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			for y := int16(0); y <= max; y++ {
				if y < f.n {
					chunk.SetBlock(x, min+y, z, 0, f.layers[f.n-y-1])
				}
				chunk.SetBiome(x, min+y, z, f.biome)
			}
		}
	}
}

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(log)
	if err != nil {
		log.Fatalln(err)
	}

	conf.Generator = func(dim world.Dimension) world.Generator {
		return NewFlat(biome.Jungle{}, []world.Block{block.Dirt{}})
	}
	conf.ReadOnlyWorld = true

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for srv.Accept(nil) {
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log server.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0o644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
