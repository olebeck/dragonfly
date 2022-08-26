package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type JungleTree struct{}

func (JungleTree) Name() string { return "minecraft:jungle_tree" }

func (t *JungleTree) CanPlace(pos cube.Pos, w *world.World) bool {
	if w.Light(pos) < 9 {
		return false
	}
	if !checkTreebox(3, 5, 3, pos, w) {
		return false
	}
	return true
}

func (t *JungleTree) Place(pos cube.Pos, w *world.World) bool {
	height := 4 + rand.Intn(6)
	if !t.CanPlace(pos, w) {
		return false
	}
	growRegularLeaves(pos, w, height, block.Leaves{Wood: block.JungleWood()})
	growStraightTrunk(pos, w, height-1, block.Log{Wood: block.JungleWood()})
	return true
}
