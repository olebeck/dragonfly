package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type OakTree struct{}

func (OakTree) Name() string { return "minecraft:oak_tree" }

func (t *OakTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 5, 3, pos, w)
}

func (t *OakTree) Place(pos cube.Pos, w *world.World) bool {
	height := 4 + rand.Intn(4)
	if !t.CanPlace(pos, w) {
		return false
	}
	growRegularLeaves(pos, w, height, block.Leaves{Wood: block.OakWood()})
	growStraightTrunk(pos, w, height-1, block.Log{Wood: block.OakWood()})
	return true
}
