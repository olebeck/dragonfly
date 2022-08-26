package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type BirchTree struct{}

func (BirchTree) Name() string { return "minecraft:birch_tree" }

func (t *BirchTree) CanPlace(pos cube.Pos, w *world.World) bool {
	if w.Light(pos) < 9 {
		return false
	}
	if !checkTreebox(3, 6, 3, pos, w) {
		return false
	}
	return true
}

func (t *BirchTree) Place(pos cube.Pos, w *world.World) bool {
	height := 5 + rand.Intn(2)
	if !t.CanPlace(pos, w) {
		return false
	}
	growRegularLeaves(pos, w, height, block.Leaves{Wood: block.BirchWood()})
	growStraightTrunk(pos, w, height-1, block.Log{Wood: block.BirchWood()})
	return true
}
