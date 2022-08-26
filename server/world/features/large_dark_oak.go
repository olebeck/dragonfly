package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type LargeDarkOakTree struct {
	height int
}

// Name ...
func (LargeDarkOakTree) Name() string { return "minecraft:large_dark_oak_tree" }

// CanPlace ...
func (t *LargeDarkOakTree) CanPlace(pos cube.Pos, w *world.World) bool {
	if t.height == 0 {
		t.height = rand.Intn(5) + 6
	}
	if !(pos.Y() >= w.Range().Min() && pos.Y()+t.height+1 < w.Range().Max()) {
		return false
	}
	return checkTreebox(4, 6, 4, pos, w)
}

// Place ...
func (t *LargeDarkOakTree) Place(pos cube.Pos, w *world.World) bool {
	placeLog := func(blockPos cube.Pos) {
		w.SetBlock(blockPos, block.Log{Wood: block.DarkOakWood()}, nil)
	}
	placeLeaf := func(blockPos cube.Pos) {
		w.SetBlock(blockPos, block.Leaves{Wood: block.DarkOakWood()}, nil)
	}

	if t.height == 0 {
		t.height = rand.Intn(5) + 6
	}

	blockPos := pos.Sub(cube.Pos{0, 1, 0})
	w.SetBlock(blockPos, block.Dirt{}, nil)
	w.SetBlock(blockPos.Side(cube.FaceEast), block.Dirt{}, nil)
	w.SetBlock(blockPos.Side(cube.FaceSouth), block.Dirt{}, nil)
	w.SetBlock(blockPos.Side(cube.FaceSouth).Side(cube.FaceEast), block.Dirt{}, nil)

	direction := randomHorizontalFace()
	i1 := t.height - rand.Intn(4)
	j1 := 2 - rand.Intn(3)
	xCenter := pos.X()
	zCenter := pos.Z()
	yLeaves := pos.Y() + t.height - 1

	// main trunk
	for y := 0; y < t.height; y++ {
		if y >= i1 && j1 > 0 {
			xCenter += pos.X() - pos.Side(direction).X()
			zCenter += pos.Z() - pos.Side(direction).Z()
			j1--
		}

		blockPos1 := cube.Pos{xCenter, pos.Y() + y, zCenter}
		if canGrowInto(w.Block(blockPos1)) {
			placeLog(blockPos1)
			placeLog(blockPos1.Side(cube.FaceEast))
			placeLog(blockPos1.Side(cube.FaceSouth))
			placeLog(blockPos1.Side(cube.FaceEast).Side(cube.FaceSouth))
		}
	}

	// top leaves
	for x := -2; x <= 0; x++ {
		for z := -3; z <= 0; z++ {
			yOff := -1
			placeLeaf(cube.Pos{0 + xCenter + x, yLeaves + yOff, 0 + zCenter + z})
			placeLeaf(cube.Pos{1 + xCenter - x, yLeaves + yOff, 0 + zCenter + z})
			placeLeaf(cube.Pos{0 + xCenter + x, yLeaves + yOff, 1 + zCenter - z})
			placeLeaf(cube.Pos{1 + xCenter - x, yLeaves + yOff, 1 + zCenter - z})

			if (x > -2 || z > -1) && (x != -1 || z != -2) {
				yOff = 1
				placeLeaf(cube.Pos{0 + xCenter + x, yLeaves + yOff, 0 + zCenter + z})
				placeLeaf(cube.Pos{1 + xCenter - x, yLeaves + yOff, 0 + zCenter + z})
				placeLeaf(cube.Pos{0 + xCenter + x, yLeaves + yOff, 1 + zCenter - z})
				placeLeaf(cube.Pos{1 + xCenter - x, yLeaves + yOff, 1 + zCenter - z})
			}
		}
	}

	// add a bit more on top
	if rand.Float64() > 0.5 {
		placeLeaf(cube.Pos{xCenter + 0, yLeaves + 2, zCenter + 0})
		placeLeaf(cube.Pos{xCenter + 1, yLeaves + 2, zCenter + 0})
		placeLeaf(cube.Pos{xCenter + 1, yLeaves + 2, zCenter + 1})
		placeLeaf(cube.Pos{xCenter + 0, yLeaves + 2, zCenter + 1})
	}

	for x := -3; x <= 4; x++ {
		for z := -3; z <= 4; z++ {
			if (x != -3 || z != -3) && (x != -3 || z != 4) && (x != 4 || z != -3) && (x != 4 || z != 4) && (abs(x) < 3 || abs(z) < 3) {
				placeLeaf(cube.Pos{xCenter + x, yLeaves, zCenter + z})
			}
		}
	}

	// add small branches near the top
	for x := -1; x <= 2; x++ {
		for z := -1; z <= 2; z++ {
			if (x < 0 || x > 1 || z < 0 || z > 1) && rand.Intn(3) <= 0 {
				branchHeight := rand.Intn(3) + 2

				for yy := 0; yy < branchHeight; yy++ {
					placeLog(cube.Pos{pos.X() + x, yLeaves - yy - 1, pos.Z() + z})
				}

				for xx := -1; xx <= 1; xx++ {
					for zz := -1; zz <= 1; zz++ {
						placeLeaf(cube.Pos{xCenter + x + xx, yLeaves, zCenter + z + zz})
					}
				}

				for xx := -2; xx <= 2; xx++ {
					for zz := -2; zz <= 2; zz++ {
						if abs(xx) != 2 || abs(zz) != 2 {
							placeLeaf(cube.Pos{xCenter + x + xx, yLeaves - 1, zCenter + z + zz})
						}
					}
				}
			}
		}
	}
	return true
}
