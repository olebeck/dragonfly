package features

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type BirchTree struct {
	Height int
	Trunk  world.Block
	Leaves world.Block
}

func (t *BirchTree) GrowTree(pos cube.Pos, w *world.World) bool {
	growRegularLeaves(pos, w, t.Height, t.Leaves)
	growStraightTrunk(pos, w, t.Height-1, t.Trunk)
	return true
}
