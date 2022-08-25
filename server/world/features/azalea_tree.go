package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type AzaleaTree struct {
	Trunk   world.Block
	Leaves  world.Block
	Leaves2 world.Block
}

func (t *AzaleaTree) Height() int {
	return 7
}

func (t *AzaleaTree) TrunkBlock() world.Block {
	return t.Trunk
}

func (t *AzaleaTree) LeafBlock() world.Block {
	if rand.Intn(16) == 0 {
		return t.Leaves2
	}
	return t.Leaves
}

func (a *AzaleaTree) GrowTree(pos cube.Pos, w *world.World) bool {
	a.growTrunk(pos, w)
	a.growLeaves(pos, w)
	return true
}

func (a *AzaleaTree) growTrunk(pos cube.Pos, w *world.World) {
	for i := 0; i < a.Height()-1; i++ {
		p := pos.Add(cube.Pos{0, i, 0})
		w.SetBlock(p, a.TrunkBlock(), nil)
	}
}

func (a *AzaleaTree) growLeaves(pos cube.Pos, w *world.World) {
	for y := pos.Y() - 3 + a.Height(); y <= pos.Y()+a.Height(); y++ {
		yOff := y - (pos.Y() + a.Height())
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
					w.SetBlock(p, a.LeafBlock(), nil)
				}
			}
		}
	}
}
