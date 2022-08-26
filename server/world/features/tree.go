package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Tree is a tree that can be grown in the world.
type Tree interface {
	// GrowTree grows a tree at pos
	GrowTree(pos cube.Pos, w *world.World) bool
}

// growRegularLeaves grows normal leaves for like an oak tree.
func growRegularLeaves(pos cube.Pos, w *world.World, height int, leaves world.Block) {
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
					w.SetBlock(p, leaves, nil)
				}
			}
		}
	}
}

func growStraightTrunk(pos cube.Pos, w *world.World, height int, trunk world.Block) {
	for i := 0; i < height; i++ {
		p := pos.Add(cube.Pos{0, i, 0})
		w.SetBlock(p, trunk, nil)
	}
}
