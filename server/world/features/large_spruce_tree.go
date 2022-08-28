package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// LargeSpruceTree is the large variant of the spruce tree, it is 2x2 blocks wide and 20 to 35 blocks high
type LargeSpruceTree struct{}

// Name ...
func (LargeSpruceTree) Name() string { return "minecraft:large_spruce_tree" }

// CanPlace ...
func (LargeSpruceTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(4, 6, 4, pos, w)
}

// Place ...
func (LargeSpruceTree) Place(pos cube.Pos, w *world.World) bool {
	height := 20 + rand.Intn(15)

	growLargeTrunk(pos, w, height, block.Log{Wood: block.SpruceWood()})
	// TODO: add leaves
	return true
}
