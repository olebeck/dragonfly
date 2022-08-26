package features

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type OakTree struct {
	Height int
	Trunk  world.Block
	Leaves world.Block
}

func (t *OakTree) GrowTree(pos cube.Pos, w *world.World) bool {
	growRegularLeaves(pos, w, t.Height, t.Leaves)
	growStraightTrunk(pos, w, t.Height-1, t.Trunk)
	return true
}
