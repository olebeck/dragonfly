package world

import (
	"bytes"
	_ "embed"
	"fmt"
	"hash/fnv"
	"image/color"
	"math"
	"sort"
	"strings"
	"unsafe"

	"github.com/brentp/intintmap"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/segmentio/fasthash/fnv1"
	"github.com/zaataylor/cartesian/cartesian"
	"golang.org/x/exp/slices"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte

	blockProperties map[string]map[string]any
	// blocks holds a list of all registered Blocks indexed by their runtime ID. Blocks that were not explicitly
	// registered are of the type unknownBlock.
	blocks     []Block
	BlockCount int
	// stateRuntimeIDs holds a map for looking up the runtime ID of a block by the stateHash it produces.
	stateRuntimeIDs map[stateHash]uint32
	// nbtBlocks holds a list of NBTer implementations for blocks registered that implement the NBTer interface.
	// These are indexed by their runtime IDs. Blocks that do not implement NBTer have a false value in this slice.
	nbtBlocks []bool
	// randomTickBlocks holds a list of RandomTicker implementations for blocks registered that implement the RandomTicker interface.
	// These are indexed by their runtime IDs. Blocks that do not implement RandomTicker have a false value in this slice.
	randomTickBlocks []bool
	// liquidBlocks holds a list of Liquid implementations for blocks registered that implement the Liquid interface.
	// These are indexed by their runtime IDs. Blocks that do not implement Liquid have a false value in this slice.
	liquidBlocks []bool
	// liquidDisplacingBlocks holds a list of LiquidDisplacer implementations for blocks registered that implement the LiquidDisplacer interface.
	// These are indexed by their runtime IDs. Blocks that do not implement LiquidDisplacer have a false value in this slice.
	liquidDisplacingBlocks []bool
	// airRID is the runtime ID of an air block.
	airRID uint32

	customBlocks []protocol.BlockEntry
)

func AirRID() uint32 {
	return airRID
}

func Blocks() []Block {
	return blocks
}

func init() {
	ClearStates()
	LoadBlockStates()
	BlockCount = len(blocks)

	chunk.RuntimeIDToState = func(runtimeID uint32) (name string, properties map[string]any, found bool) {
		if runtimeID >= uint32(len(blocks)) {
			return "", nil, false
		}
		name, properties = blocks[runtimeID].EncodeBlock()
		return name, properties, true
	}
	chunk.StateToRuntimeID = func(name string, properties map[string]any) (runtimeID uint32, found bool) {
		if rid, ok := stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]; ok {
			return rid, true
		}
		rid, ok := stateRuntimeIDs[stateHash{name: name, properties: hashProperties(blockProperties[name])}]
		return rid, ok
	}
}

func ClearStates() {
	blockProperties = map[string]map[string]any{}
	stateRuntimeIDs = map[stateHash]uint32{}
	chunk.HashToRuntimeID = make(map[uint32]uint32)
	hashes = intintmap.New(7000, 0.999)

	customBlocks = nil
	blocks = nil
	nbtBlocks = nil
	randomTickBlocks = nil
	liquidBlocks = nil
	liquidDisplacingBlocks = nil
	chunk.FilteringBlocks = nil
	chunk.LightBlocks = nil
	chunk.WaterBlocks = nil
}

func LoadBlockStates() {
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	var states []blockState
	for {
		s := blockState{}
		if err := dec.Decode(&s); err != nil {
			break
		}
		states = append(states, s)
	}
	registerBlockStates(states)
}

func networkBlockHash(name string, properties map[string]any) uint32 {
	data, err := nbt.MarshalEncoding(map[string]any{
		"name":   name,
		"states": properties,
	}, nbt.LittleEndian)
	if err != nil {
		panic(err)
	}

	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

// registerBlockStates registers multiple new blockStates to the states slice. The function panics if the properties the
// blockState hold are invalid or if the blockState was already registered.
func registerBlockStates(ss []blockState) {
	newStates := map[stateHash]uint32{}

	// add blocks
	for _, s := range ss {
		blocks = append(blocks, UnknownBlock{s})
		newStates[stateHash{s.Name, hashProperties(s.Properties)}] = 0
	}
	BlockCount = len(blocks)

	// sort the new blocks
	sort.SliceStable(blocks, func(i, j int) bool {
		nameOne, _ := blocks[i].EncodeBlock()
		nameTwo, _ := blocks[j].EncodeBlock()
		return fnv1.HashString64(nameOne) < fnv1.HashString64(nameTwo)
	})

	for id, b := range blocks {
		name, properties := b.EncodeBlock()
		i := stateHash{name: name, properties: hashProperties(properties)}
		if _, ok := blockProperties[name]; !ok {
			blockProperties[name] = properties
		}
		rid := uint32(id)
		if name == "minecraft:air" {
			airRID = rid
		}

		isWater := name == "minecraft:water"

		// if its one of the added ones
		if _, ok := newStates[i]; ok {
			if _, ok := stateRuntimeIDs[i]; ok {
				panic(fmt.Sprintf("cannot register the same state twice (%+v)", b))
			}

			nbtBlocks = slices.Insert(nbtBlocks, int(rid), false)
			randomTickBlocks = slices.Insert(randomTickBlocks, int(rid), false)
			liquidBlocks = slices.Insert(liquidBlocks, int(rid), false)
			liquidDisplacingBlocks = slices.Insert(liquidDisplacingBlocks, int(rid), false)
			chunk.FilteringBlocks = slices.Insert(chunk.FilteringBlocks, int(rid), 15)
			chunk.LightBlocks = slices.Insert(chunk.LightBlocks, int(rid), 0)
			chunk.WaterBlocks = slices.Insert(chunk.WaterBlocks, int(rid), isWater)
		}
		stateRuntimeIDs[i] = rid
		chunk.HashToRuntimeID[networkBlockHash(name, properties)] = rid
		hashes.Put(int64(b.Hash()), int64(id))
	}
}

func ns_name_split(identifier string) (ns, name string) {
	ns_name := strings.Split(identifier, ":")
	return ns_name[0], ns_name[len(ns_name)-1]
}

type BlockState = blockState

var traitLookup = map[string][]any{
	"facing_direction": {
		"north", "east", "south", "west", "down", "up",
	},
}

func InsertCustomBlocks(entries []protocol.BlockEntry) []BlockState {
	customBlocks = entries
	var states []blockState
	for _, entry := range entries {
		ns, _ := ns_name_split(entry.Name)
		if ns == "minecraft" {
			continue
		}

		var propertyNames []string
		var propertyValues []any

		props, ok := entry.Properties["properties"].([]any)
		if ok {
			for _, v := range props {
				v := v.(map[string]any)
				name := v["name"].(string)
				enum := v["enum"]
				propertyNames = append(propertyNames, name)
				propertyValues = append(propertyValues, enum)
			}
		}

		traits, ok := entry.Properties["traits"].([]any)
		if ok {
			for _, trait := range traits {
				trait := trait.(map[string]any)
				enabled_states := trait["enabled_states"].(map[string]any)
				for k, enabled := range enabled_states {
					enabled := enabled.(uint8)
					if enabled == 0 {
						continue
					}
					v, ok := traitLookup[k]
					if !ok {
						panic("unresolved trait " + k)
					}

					propertyNames = append(propertyNames, "minecraft:"+k)
					propertyValues = append(propertyValues, v)
				}
			}
		}

		permutations := cartesian.NewCartesianProduct(propertyValues).Values()

		for _, values := range permutations {
			m := make(map[string]any)
			for i, value := range values {
				name := propertyNames[i]
				m[name] = value
			}
			states = append(states, blockState{
				Name:       entry.Name,
				Properties: m,
			})
		}
	}
	registerBlockStates(states)
	return states
}

func CustomBlocks() []protocol.BlockEntry {
	return customBlocks
}

// UnknownBlock represents a block that has not yet been implemented. It is used for registering block
// states that haven't yet been added.
type UnknownBlock struct {
	blockState
}

// EncodeBlock ...
func (b UnknownBlock) EncodeBlock() (string, map[string]any) {
	return b.Name, b.Properties
}

// Model ...
func (UnknownBlock) Model() BlockModel {
	return unknownModel{}
}

// Hash ...
func (b UnknownBlock) Hash() uint64 {
	return math.MaxUint64
}

func (b UnknownBlock) Color() color.RGBA {
	return color.RGBA{255, 0, 255, 255}
}

func (b UnknownBlock) DecodeNBT(data map[string]any) any {
	b.Properties = data
	return b
}

// EncodeNBT encodes the entity into a map which can then be encoded as NBT to be written.
func (b UnknownBlock) EncodeNBT() map[string]any {
	return b.Properties
}

// blockState holds a combination of a name and properties, together with a version.
type blockState struct {
	Name       string         `nbt:"name"`
	Properties map[string]any `nbt:"states"`
	Version    int32          `nbt:"version"`
}

// stateHash is a struct that may be used as a map key for block states. It contains the name of the block state
// and an encoded version of the properties.
type stateHash struct {
	name, properties string
}

// hashProperties produces a hash for the block properties held by the blockState.
func hashProperties(properties map[string]any) string {
	if properties == nil {
		return ""
	}
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	var b strings.Builder
	for _, k := range keys {
		switch v := properties[k].(type) {
		case bool:
			if v {
				b.WriteByte(1)
			} else {
				b.WriteByte(0)
			}
		case uint8:
			b.WriteByte(v)
		case int32:
			a := *(*[4]byte)(unsafe.Pointer(&v))
			b.Write(a[:])
		case string:
			b.WriteString(v)
		default:
			// If block encoding is broken, we want to find out as soon as possible. This saves a lot of time
			// debugging in-game.
			panic(fmt.Sprintf("invalid block property type %T for property %v", v, k))
		}
	}

	return b.String()
}
