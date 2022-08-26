package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type AzaleaTree struct{}

func (AzaleaTree) Name() string { return "minecraft:azalea_tree" }

func (t *AzaleaTree) LeafBlock() world.Block {
	if rand.Intn(8) == 0 {
		return block.AzaleaLeaves{Flowering: true}
	}
	return block.AzaleaLeaves{}
}

func (AzaleaTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(5, 6, 5, pos, w)
}

func (a *AzaleaTree) Place(pos cube.Pos, w *world.World) bool {
	height := 7
	a.growLeaves(pos, w, height)
	a.growTrunk(pos, w, height)
	return true
}

func (a *AzaleaTree) growTrunk(pos cube.Pos, w *world.World, height int) {
	for i := 0; i < height-1; i++ {
		p := pos.Add(cube.Pos{0, i, 0})
		w.SetBlock(p, block.Log{Wood: block.OakWood()}, nil)
	}
}

func (a *AzaleaTree) growLeaves(pos cube.Pos, w *world.World, height int) {
	for y := pos.Y() - 3 + height; y <= pos.Y()+height; y++ {
		yOff := y - (pos.Y() + height)
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
