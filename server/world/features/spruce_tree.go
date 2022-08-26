package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type SpruceTree struct {
}

func (SpruceTree) Name() string { return "minecraft:spruce_tree" }

func (t *SpruceTree) CanPlace(pos cube.Pos, w *world.World) bool {
	if !checkTreebox(4, 6, 4, pos, w) {
		return false
	}
	return true
}

func (t *SpruceTree) Place(pos cube.Pos, w *world.World) bool {
	height := 6 + rand.Intn(4)
	growSpruceLeaves(pos, w, height, height-(1+rand.Intn(2)), 3+rand.Intn(2), block.Leaves{Wood: block.SpruceWood()})
	growStraightTrunk(pos, w, height-rand.Intn(3), block.Log{Wood: block.SpruceWood()})
	return true
}

func growSpruceLeaves(pos cube.Pos, w *world.World, height, top, lRadius int, leaves world.Block) {
	radius := rand.Intn(2)
	maxR := 1
	minR := 0

	for y := 0; y <= top; y++ {
		yy := pos.Y() + height - y

		for x := pos.X() - radius; x <= pos.X()+radius; x++ {
			xOff := abs(x - pos.X())
			for z := pos.Z() - radius; z <= pos.Z()+radius; z++ {
				zOff := abs(z - pos.Z())
				if xOff == radius && zOff == radius && radius > 0 {
					continue
				}

				p := cube.Pos{x, yy, z}
				if true /* !w.Block(p).(Solid) */ {
					w.SetBlock(p, leaves, nil)
				}
			}
		}

		if radius >= maxR {
			radius = minR
			minR = 1
			maxR++
			if maxR > lRadius {
				maxR = lRadius
			}
		} else {
			radius++
		}
	}
}
