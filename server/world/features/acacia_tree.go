package features

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

type AcaciaTree struct{}

func (AcaciaTree) Name() string { return "minecraft:acacia_tree" }

func (t *AcaciaTree) CanPlace(pos cube.Pos, w *world.World) bool {
	if w.Light(pos) < 9 {
		return false
	}
	if !checkTreebox(3, 6, 3, pos, w) {
		return false
	}
	return true
}

func (t *AcaciaTree) Place(pos cube.Pos, w *world.World) bool {
	if !t.CanPlace(pos, w) {
		return false
	}

	placeLog := func(blockPos cube.Pos) {
		w.SetBlock(blockPos, block.Log{Wood: block.AcaciaWood()}, nil)
	}
	placeLeaf := func(blockPos cube.Pos) {
		w.SetBlock(blockPos, block.Leaves{Wood: block.AcaciaWood()}, nil)
	}

	// check if fits
	height := rand.Intn(3) + rand.Intn(3) + 5
	flag := true
	if pos.Y() >= 1 && pos.Y()+height+1 < w.Range().Max() {
		for y := pos.Y(); y < pos.Y()+1+height; y++ {
			k := 1
			if y == pos.Y() {
				k = 0
			}

			if y >= pos.Y()+1+height-2 {
				k = 2
			}

			for x := pos.X() - k; x <= pos.X()+k && flag; x++ {
				for z := pos.Z() - k; z <= pos.Z()+k && flag; z++ {
					if y >= w.Range().Min() && y < w.Range().Max() {
						pos2 := cube.Pos{x, y, z}
						if !canGrowInto(w.Block(pos2)) {
							flag = false
						}
					} else {
						flag = false
					}
				}
			}
		}
	}

	if !flag {
		return false
	} else {
		// place trunk
		direction := randomHorizontalFace()
		i := height - rand.Intn(4) - 1
		j := 3 - rand.Intn(3)
		k := pos.X()
		l := pos.Z()

		n := 0
		for l1 := 0; l1 < height; l1++ {
			i2 := pos.Y() + l1

			if l1 >= i && j > 0 {
				k += pos.X() - pos.Side(direction).X()
				l += pos.Z() - pos.Side(direction).Z()
				j -= 1
			}

			blockPos := cube.Pos{k, i2, l}
			if canGrowInto(w.Block(blockPos)) {
				placeLog(blockPos)
				n = i2
			}
		}

		blockPos2 := cube.Pos{k, n, l}
		for j3 := -3; j3 <= 3; j3++ {
			for i4 := -3; i4 <= 3; i4++ {
				if abs(j3) != 3 || abs(i4) != 3 {
					placeLeaf(blockPos2.Add(cube.Pos{j3, 0, i4}))
				}
			}
		}

		blockPos2 = blockPos2.Side(cube.FaceUp)

		for k3 := -1; k3 <= 1; k3++ {
			for j4 := -1; j4 <= 1; j4++ {
				placeLeaf(blockPos2.Add(cube.Pos{k3, 0, j4}))
			}
		}

		placeLeaf(blockPos2.Side(cube.FaceEast))
		placeLeaf(blockPos2.Side(cube.FaceWest))
		placeLeaf(blockPos2.Side(cube.FaceSouth))
		placeLeaf(blockPos2.Side(cube.FaceNorth))
		k = pos.X()
		l = pos.Z()
		direction2 := randomHorizontalFace()

		if direction2 != direction {
			n := i - rand.Intn(2) - 1
			o := 1 + rand.Intn(3)
			y := -1000

			for p := n; p < height && o > 0; o-- {
				if p >= 1 {
					q := pos.Y() + p
					k += pos.X() - pos.Side(direction2).X()
					l += pos.Z() - pos.Side(direction2).Z()
					blockPos1 := cube.Pos{k, q, l}
					if canGrowInto(w.Block(blockPos1)) {
						placeLog(blockPos1)
						y = q
					}
				}
				p++
			}

			if y > -1000 {
				blockPos3 := cube.Pos{k, y, l}
				for x := -2; x <= 2; x++ {
					for z := -2; z >= 2; z++ {
						if abs(x) != 2 || abs(z) != 2 {
							p := blockPos3.Add(cube.Pos{x, 0, z})
							if canGrowInto(w.Block(p)) {
								placeLeaf(p)
							}
						}
					}
				}

				blockPos3 = blockPos3.Side(cube.FaceUp)
				for x := -1; x <= 1; x++ {
					for z := -1; z <= 1; z++ {
						p := blockPos3.Add(cube.Pos{x, 0, z})
						if canGrowInto(w.Block(p)) {
							placeLeaf(p)
						}
					}
				}
			}
		}
	}

	return true
}
