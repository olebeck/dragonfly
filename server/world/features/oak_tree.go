package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type OakTree struct {
	Trunk  world.Block
	Leaves world.Block
}

func (t *OakTree) Height() int {
	return 7
}

func (t *OakTree) TrunkBlock() world.Block {
	return t.Trunk
}

func (t *OakTree) LeafBlock() world.Block {
	return t.Leaves
}

func (t *OakTree) GrowTree(pos cube.Pos, w *world.World) bool {
	t.growTrunk(pos, w)
	t.growLeaves(pos, w)
	return true
}

func (t *OakTree) growTrunk(pos cube.Pos, w *world.World) {
	for i := 0; i < t.Height()-1; i++ {
		p := pos.Add(cube.Pos{0, i, 0})
		w.SetBlock(p, t.TrunkBlock(), nil)
	}
}

func (t *OakTree) growLeaves(pos cube.Pos, w *world.World) {
	for y := pos.Y() - 3 + t.Height(); y <= pos.Y()+t.Height(); y++ {
		yOff := y - (pos.Y() + t.Height())
		mid := int(1 - yOff/2)
		for x := pos.X() - mid; x <= pos.X()+mid; x++ {
			xOff := abs(x - pos.X())
			for z := pos.Z() - mid; z <= pos.Z()+mid; z++ {
				zOff := abs(z - pos.Z())
				if xOff == mid && zOff == mid && (yOff == 0 || rand.Intn(2) == 0) {
					continue
				}

				p := cube.Pos{x, y, z}
				if true /* !w.Block(p).(Solid) */ {
					w.SetBlock(p, t.LeafBlock(), nil)
				}
			}
		}
	}
}
