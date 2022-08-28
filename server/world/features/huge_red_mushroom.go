package features

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// HugeRedMushroom is a large mushroom grown from a mushroom with bone meal
type HugeRedMushroom struct{}

// Name ...
func (HugeRedMushroom) Name() string { return "minecraft:huge_red_mushroom" }

// CanPlace ...
func (m *HugeRedMushroom) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 6, 3, pos, w)
}

// Place ...
func (m *HugeRedMushroom) Place(pos cube.Pos, w *world.World) bool {
	const height = 4

	growMushroomStem(pos, w, height)
	m.growTop(pos, w, height)
	return true
}

func mushroomHugeBits(withWest, withNorth, withSouth, withEast bool) (hugeBits int) {
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
	return hugeBits
}

func (m *HugeRedMushroom) growTop(pos cube.Pos, w *world.World, height int) {
	const foliageRadius = 2
	for i := height - 3; i <= height; i++ {
		radius := foliageRadius - 1
		if i < height {
			radius = foliageRadius
		}
		k := foliageRadius - 2

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
				withWest := edgeNorth || corner2 && x < -k
				withEast := edgeSouth || corner2 && x > k
				withNorth := edgeWest || corner1 && z < -k
				withSouth := edgeEast || corner1 && z > k

				hugeBits := mushroomHugeBits(withWest, withNorth, withSouth, withEast)

				// dont place any on the corners except top one
				if i >= height || (corner1 != corner2) {
					pos2 := pos.Add(cube.Pos{x, i, z})
					if canGrowInto(w.Block(pos2)) {
						w.SetBlock(pos2, block.MushroomBlock{Type: block.Red(), HugeBits: hugeBits}, nil)
					}
				}
			}
		}
	}
}
