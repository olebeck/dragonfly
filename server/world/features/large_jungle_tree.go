package features

import (
	"math"
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// LargeJungleTree is the large type of the jungle tree, grown using 4 saplings.
type LargeJungleTree struct{}

// Name ...
func (LargeJungleTree) Name() string { return "minecraft:large_jungle_tree" }

// CanPlace ...
func (t *LargeJungleTree) CanPlace(pos cube.Pos, w *world.World) bool {
	return checkTreebox(3, 5, 3, pos, w)
}

// Place ...
func (t *LargeJungleTree) Place(pos cube.Pos, w *world.World) bool {
	const baseHeight = 10
	const extraRandomHeight = 20

	height := rand.Intn(3) + baseHeight + rand.Intn(extraRandomHeight)

	placeLog := func(blockPos cube.Pos) {
		w.SetBlock(blockPos, block.Log{Wood: block.JungleWood()}, nil)
	}

	// builds the leaves on the top
	t.createCrown(pos.Add(cube.Pos{0, height, 0}), w, 2)

	// adds small branches off of the tree
	for j := pos.Y() + height - 2 - rand.Intn(4); j > pos.Y()+height/2; j -= 2 + rand.Intn(4) {
		f := rand.Float64() * math.Pi * 2
		var k, l int
		for i1 := 0; i1 < 5; i1++ {
			k = pos.X() + int(1.5+math.Cos(f)*float64(i1))
			l = pos.Z() + int(1.5+math.Sin(f)*float64(i1))
			placeLog(cube.Pos{k, j - 3 + i1/2, l})
		}

		j2 := 1 + rand.Intn(2)
		for k1 := j - j2; k1 <= j; k1++ {
			t.growLeavesLayer(cube.Pos{k, k1, l}, w, 1-(k1-j))
		}
	}

	// builds the trunk, with vines
	for i2 := 0; i2 < height; i2++ {
		blockPos := pos.Add(cube.Pos{0, i2, 0})

		if canGrowInto(w.Block(blockPos)) {
			placeLog(blockPos)

			if i2 > 0 {
				t.placeVine(w, blockPos.Side(cube.FaceWest), 8)
				t.placeVine(w, blockPos.Side(cube.FaceNorth), 1)
			}
		}

		if i2 < height-1 {
			blockPos1 := blockPos.Side(cube.FaceEast)

			if canGrowInto(w.Block(blockPos1)) {
				placeLog(blockPos1)

				if i2 > 0 {
					t.placeVine(w, blockPos1.Side(cube.FaceEast), 2)
					t.placeVine(w, blockPos1.Side(cube.FaceNorth), 1)
				}
			}

			blockPos2 := blockPos.Side(cube.FaceSouth).Side(cube.FaceEast)

			if canGrowInto(w.Block(blockPos2)) {
				placeLog(blockPos2)

				if i2 > 0 {
					t.placeVine(w, blockPos2.Side(cube.FaceEast), 2)
					t.placeVine(w, blockPos2.Side(cube.FaceSouth), 4)
				}
			}

			blockPos3 := blockPos.Side(cube.FaceSouth)

			if canGrowInto(w.Block(blockPos3)) {
				placeLog(blockPos3)

				if i2 > 0 {
					t.placeVine(w, blockPos3.Side(cube.FaceWest), 8)
					t.placeVine(w, blockPos3.Side(cube.FaceSouth), 4)
				}
			}
		}
	}
	return true
}

// createCrown grows the leaves at the top of the large jungle tree
func (t *LargeJungleTree) createCrown(pos cube.Pos, w *world.World, width int) {
	for j := -2; j <= 2; j++ {
		t.growLeavesLayerStrict(pos.Add(cube.Pos{0, j, 0}), w, width+1-j)
	}
}

// growLeavesLayerStrict generates a circular leaf layer on the tree
func (t *LargeJungleTree) growLeavesLayerStrict(pos cube.Pos, w *world.World, width int) {
	widthSquared := width * width
	for j := -width; j <= width+1; j++ {
		for k := -width; k <= width+1; k++ {
			l := j - 1
			i1 := k - 1
			if j*j+k*k <= widthSquared || l*l+i1*i1 <= widthSquared || j*j+i1*i1 <= widthSquared || l*l+k*k <= widthSquared {
				blockPos := pos.Add(cube.Pos{j, 0, k})
				if canGrowInto(w.Block(blockPos)) {
					w.SetBlock(blockPos, block.Leaves{Wood: block.JungleWood()}, nil)
				}
			}
		}
	}
}

// placeVine adds a vine to the tree at the position, with a certain chance
func (t *LargeJungleTree) placeVine(w *world.World, pos cube.Pos, i1 int) {
	_, isAir := w.Block(pos).(block.Air)
	if rand.Intn(3) > 0 && isAir {
		//w.SetBlock(pos, Vine)
	}
}

// growLeavesLayer grows a lower leaf layer on the tree
func (t *LargeJungleTree) growLeavesLayer(pos cube.Pos, w *world.World, width int) {
	widthSquared := width * width
	for j := -width; j <= width+1; j++ {
		for k := -width; k <= width+1; k++ {
			if j*j+k*k <= widthSquared {
				blockPos := pos.Add(cube.Pos{j, 0, k})
				if canGrowInto(w.Block(blockPos)) {
					w.SetBlock(blockPos, block.Leaves{Wood: block.JungleWood()}, nil)
				}
			}
		}
	}
}
