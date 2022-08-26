package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// BirchTree is a tree that appears similar to a common oak in terms of height and leaves, but with light bark and pale wood.
type BirchTree struct{}

// Name ...
func (BirchTree) Name() string { return "minecraft:birch_tree" }

// CanPlace ...
func (t *BirchTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 6, 3, pos, w)
}

// Place ...
func (t *BirchTree) Place(pos cube.Pos, w *world.World) bool {
	height := 5 + rand.Intn(2)

	growRegularLeaves(pos, w, height, block.Leaves{Wood: block.BirchWood()})
	growStraightTrunk(pos, w, height-1, block.Log{Wood: block.BirchWood()})
	return true
}
