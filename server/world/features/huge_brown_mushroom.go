package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// HugeBrownMushroom is a large mushroom grown from a mushroom with bone meal
type HugeBrownMushroom struct{}

// Name ...
func (HugeBrownMushroom) Name() string { return "minecraft:huge_brown_mushroom" }

// CanPlace ...
func (m *HugeBrownMushroom) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 6, 3, pos, w)
}

// Place ...
func (m *HugeBrownMushroom) Place(pos cube.Pos, w *world.World) bool {
	height := 4 + rand.Intn(9)
	if rand.Intn(12) == 0 {
		height *= 2
	}

	growMushroomStem(pos, w, height)
	m.growTop(pos, w, height)
	return true
}

func (m *HugeBrownMushroom) growTop(pos cube.Pos, w *world.World, height int) {
	const radius = 3
	for x := -radius; x <= radius; x++ {
		for z := -radius; z <= radius; z++ {
			// what edge of the top its on
			edgeNorth := x == -radius
			edgeSouth := x == radius
			edgeWest := z == -radius
			edgeEast := z == radius

			corner1 := (edgeNorth || edgeSouth)
			corner2 := (edgeWest || edgeEast)

			// what sides the block should a cap on
			withWest := edgeNorth || corner2 && x == 1-radius
			withEast := edgeSouth || corner2 && x == radius-1
			withNorth := edgeWest || corner1 && z == 1-radius
			withSouth := edgeEast || corner1 && z == radius-1

			hugeBits := mushroomHugeBits(withWest, withNorth, withSouth, withEast)

			// dont place any on the corners
			if !corner1 || !corner2 {
				pos2 := pos.Add(cube.Pos{x, height, z})
				if canGrowInto(w.Block(pos2)) {
					w.SetBlock(pos2, block.MushroomBlock{Type: block.Brown(), HugeBits: hugeBits}, nil)
				}
			}
		}
	}
}

// growMushroomStem grows a mushroom stem upto the height
func growMushroomStem(pos cube.Pos, w *world.World, height int) {
	for i := 0; i < height; i++ {
		pos2 := pos.Add(cube.Pos{0, i, 0})
		if canGrowInto(w.Block(pos2)) {
			w.SetBlock(pos2, block.MushroomBlock{Type: block.Brown(), HugeBits: 10}, nil)
		}
	}
}
