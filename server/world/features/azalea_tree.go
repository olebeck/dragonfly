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
	height := 4 + rand.Intn(2)
	bend := rand.Intn(2)
	top := growBendyTrunk(pos, w, height, bend, block.Log{Wood: block.OakWood()})
	randomSpreadFoliage(top, w, block.AzaleaLeaves{}, 3, 0, 2, 50)
	return true
}

func growBendyTrunk(pos cube.Pos, w *world.World, height int, bend int, trunk world.Block) (top cube.Pos) {
	direction := randomHorizontalFace()
	i := height - 1

	for j := 0; j <= i; j++ {
		if j+1 >= i+rand.Intn(2) {
			pos = pos.Side(direction)
		}
		b := w.Block(pos)
		_, azalea := b.(block.Azalea)
		if canGrowInto(b) || azalea {
			w.SetBlock(pos, trunk, nil)
		}
		pos = pos.Side(cube.FaceUp)
	}

	for k := 0; k <= bend; k++ {
		if canGrowInto(w.Block(pos)) {
			w.SetBlock(pos, trunk, nil)
		}
		pos = pos.Side(direction)
	}
	return pos.Side(direction.Opposite())
}

func randomSpreadFoliage(pos cube.Pos, w *world.World, leaf world.Block, radius, offset, foliageHeight, attempts int) {
	for i := 0; i < attempts; i++ {
		p := pos.Add(cube.Pos{rand.Intn(radius) - rand.Intn(radius), rand.Intn(foliageHeight) - rand.Intn(foliageHeight), rand.Intn(radius) - rand.Intn(radius)})
		if canGrowInto(w.Block(p)) {
			w.SetBlock(p, leaf, nil)
		}
	}
}
