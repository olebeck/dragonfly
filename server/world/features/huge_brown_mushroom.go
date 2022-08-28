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
	radius := 3
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			// what edge of the top its on
			edgeNorth := x == -radius
			edgeSouth := x == radius
			edgeWest := y == -radius
			edgeEast := y == radius

			corner1 := (edgeNorth || edgeSouth)
			corner2 := (edgeWest || edgeEast)

			// what sides the block should a cap on
			withWest := edgeNorth || corner2 && x == 1-radius
			withEast := edgeSouth || corner2 && x == radius-1
			withNorth := edgeWest || corner1 && y == 1-radius
			withSouth := edgeEast || corner1 && y == radius-1

			hugeBits := 0
			if withWest {
				if withNorth {
					hugeBits = 1 // top + west + north
				} else if withSouth {
					hugeBits = 7 // top + west + south
				} else {
					hugeBits = 4 // top + west
				}
			} else if withEast {
				if withNorth {
					hugeBits = 3 // top + north + east
				} else if withSouth {
					hugeBits = 9 // top + south + east
				} else {
					hugeBits = 6 // top + south
				}
			} else if withNorth {
				hugeBits = 2 // top + north
			} else if withSouth {
				hugeBits = 8 // top + south
			} else {
				hugeBits = 5 // only top
			}

			// dont place any on the corners
			if !corner1 || !corner2 {
				pos2 := pos.Add(cube.Pos{x, height, y})
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
