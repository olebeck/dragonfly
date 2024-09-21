package world

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/customblock"
)

// Block is a block that may be placed or found in a world. In addition, the block may also be added to an
// inventory: It is also an item.
// Every Block implementation must be able to be hashed as key in a map.
type Block interface {
	// EncodeBlock encodes the block to a string ID such as 'minecraft:grass' and properties associated
	// with the block.
	EncodeBlock() (string, map[string]any)
	// Hash returns two different identifiers for the block. The first is the base hash which is unique for
	// each type of block at runtime. For vanilla blocks, this is an auto-incrementing constant and for custom
	// blocks, you can call block.NextHash() to get a unique identifier. The second is the hash of the block's
	// own state and does not need to worry about colliding with other types of blocks. This is later combined
	// with the base hash to create a unique identifier for the full block.
	Hash() (uint64, uint64)
	// Model returns the BlockModel of the Block.
	Model() BlockModel
	// Color returns an RGBA color used to represent this block on a map
	Color() color.RGBA
}

// CustomBlock represents a block that is non-vanilla and requires a resource pack and extra steps to show it to the
// client.
type CustomBlock interface {
	Block
	Properties() customblock.Properties
}

type CustomBlockBuildable interface {
	CustomBlock
	// Name is the name displayed to clients using the block.
	Name() string
	// Geometries is the geometries for the block that define the shape of the block. If false is returned, no custom
	// geometry will be applied. Permutation-specific geometry can be defined by returning a map of permutations to
	// geometry.
	Geometry() []byte
	// Textures is a map of images indexed by their target, used to map textures on to the block. Permutation-specific
	// textures can be defined by returning a map of permutations to textures.
	Textures() map[string]image.Image
}

// Liquid represents a block that can be moved through and which can flow in the world after placement. There
// are two liquids in vanilla, which are lava and water.
type Liquid interface {
	Block
	// LiquidDepth returns the current depth of the liquid.
	LiquidDepth() int
	// SpreadDecay returns the amount of depth that is subtracted from the liquid's depth when it spreads to
	// a next block.
	SpreadDecay() int
	// WithDepth returns the liquid with the depth passed.
	WithDepth(depth int, falling bool) Liquid
	// LiquidFalling checks if the liquid is currently considered falling down.
	LiquidFalling() bool
	// BlastResistance is the blast resistance of the liquid, which influences the liquid's ability to withstand an
	// explosive blast.
	BlastResistance() float64
	// LiquidType returns an int unique for the liquid, used to check if two liquids are considered to be
	// of the same type.
	LiquidType() string
	// Harden checks if the block should harden when looking at the surrounding blocks and sets the position
	// to the hardened block when adequate. If the block was hardened, the method returns true.
	Harden(pos cube.Pos, w *World, flownIntoBy *cube.Pos) bool
}

// RandomTicker represents a block that executes an action when it is ticked randomly. Every 20th of a second,
// one random block in each sub chunk are picked to receive a random tick.
type RandomTicker interface {
	// RandomTick handles a random tick of the block at the position passed. Additionally, a rand.Rand
	// instance is passed which may be used to generate values randomly without locking.
	RandomTick(pos cube.Pos, w *World, r *rand.Rand)
}

// ScheduledTicker represents a block that executes an action when it has a block update scheduled, such as
// when a block adjacent to it is broken.
type ScheduledTicker interface {
	// ScheduledTick handles a scheduled tick initiated by an event in one of the neighbouring blocks, such as
	// when a block is placed or broken. Additionally, a rand.Rand instance is passed which may be used to
	// generate values randomly without locking.
	ScheduledTick(pos cube.Pos, w *World, r *rand.Rand)
}

// TickerBlock is an implementation of NBTer with an additional Tick method that is called on every world
// tick for loaded blocks that implement this interface.
type TickerBlock interface {
	NBTer
	Tick(currentTick int64, pos cube.Pos, w *World)
}

// NeighbourUpdateTicker represents a block that is updated when a block adjacent to it is updated, either
// through placement or being broken.
type NeighbourUpdateTicker interface {
	// NeighbourUpdateTick handles a neighbouring block being updated. The position of that block and the
	// position of this block is passed.
	NeighbourUpdateTick(pos, changedNeighbour cube.Pos, w *World)
}

// NBTer represents either an item or a block which may decode NBT data and encode to NBT data. Typically,
// this is done to store additional data.
type NBTer interface {
	// DecodeNBT returns the (new) item, block or entity, depending on which of those the NBTer was, with the NBT data
	// decoded into it.
	DecodeNBT(data map[string]any) any
	// EncodeNBT encodes the entity into a map which can then be encoded as NBT to be written.
	EncodeNBT() map[string]any
}

// LiquidDisplacer represents a block that is able to displace a liquid to a different world layer, without
// fully removing the liquid.
type LiquidDisplacer interface {
	// CanDisplace specifies if the block is able to displace the liquid passed.
	CanDisplace(b Liquid) bool
	// SideClosed checks if a position on the side of the block placed in the world at a specific position is
	// closed. When this returns true (for example, when the side is below the position and the block is a
	// slab), liquid inside the displacer won't flow from pos into side.
	SideClosed(pos, side cube.Pos, w *World) bool
}

// lightEmitter is identical to a block.LightEmitter.
type lightEmitter interface {
	LightEmissionLevel() uint8
}

// lightDiffuser is identical to a block.LightDiffuser.
type lightDiffuser interface {
	LightDiffusionLevel() uint8
}

// replaceableBlock represents a block that may be replaced by another block automatically. An example is
// grass, which may be replaced by clicking it with another block.
type replaceableBlock interface {
	// ReplaceableBy returns a bool which indicates if the block is replaceable by another block.
	ReplaceableBy(b Block) bool
}

// replaceable checks if the block at the position passed is replaceable with the block passed.
func replaceable(w *World, c *Column, pos cube.Pos, with Block) bool {
	if r, ok := w.blockInChunk(c, pos).(replaceableBlock); ok {
		return r.ReplaceableBy(with)
	}
	return false
}

// BlockAction represents an action that may be performed by a block. Typically, these actions are sent to
// viewers in a world so that they can see these actions.
type BlockAction interface {
	BlockAction()
}

// UnknownBlock represents a block that has not yet been implemented. It is used for registering block
// states that haven't yet been added.
type UnknownBlock struct {
	BlockState
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
func (b UnknownBlock) Hash() (uint64, uint64) {
	return 0, math.MaxUint64
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

func BlockHash(b Block) uint64 {
	return DefaultBlockRegistry.BlockHash(b)
}
