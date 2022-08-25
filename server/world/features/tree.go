package features

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type Tree interface {
	// GrowTree grows a tree at pos
	GrowTree(pos cube.Pos, w *world.World) bool
	// Height is the total height of the tree
	Height() int
	// TrunkBlock is the block used for the trunk
	TrunkBlock() world.Block
	// LeafBlock is the block used for the leaves of the tree
	LeafBlock() world.Block
}
