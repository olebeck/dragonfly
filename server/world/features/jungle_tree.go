package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// JungleTree is a type of tree generated in the jungle biome.
type JungleTree struct{}

// Name ...
func (JungleTree) Name() string { return "minecraft:jungle_tree" }

// CanPlace ...
func (t *JungleTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 5, 3, pos, w)
}

// Place ...
func (t *JungleTree) Place(pos cube.Pos, w *world.World) bool {
	height := 4 + rand.Intn(6)

	growRegularLeaves(pos, w, height, block.Leaves{Wood: block.JungleWood()})
	growStraightTrunk(pos, w, height-1, block.Log{Wood: block.JungleWood()})
	return true
}
