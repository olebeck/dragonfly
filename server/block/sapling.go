package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Sapling is a non-solid block that can be grown into a tree.
type Sapling struct {
	empty
	transparent

	Wood WoodType

	AgeBit bool
}

// findSaplings finds the same sapling type in a 2x2 area, returns the lowest xz coordinates.
func (s Sapling) findSaplings(pos cube.Pos, w *world.World) (*cube.Pos, bool) {
	validPositions := [][]cube.Pos{
		[]cube.Pos{
			pos, pos.Add(cube.Pos{1, 0, 0}), pos.Add(cube.Pos{0, 0, 1}), pos.Add(cube.Pos{1, 0, 1}),
		},
		[]cube.Pos{
			pos, pos.Add(cube.Pos{-1, 0, 0}), pos.Add(cube.Pos{0, 0, -1}), pos.Add(cube.Pos{-1, 0, -1}),
		},
		[]cube.Pos{
			pos, pos.Add(cube.Pos{1, 0, 0}), pos.Add(cube.Pos{0, 0, -1}), pos.Add(cube.Pos{1, 0, -1}),
		},
		[]cube.Pos{
			pos, pos.Add(cube.Pos{-1, 0, 0}), pos.Add(cube.Pos{-1, 0, 0}), pos.Add(cube.Pos{-1, 0, 1}),
		},
	}
	for _, v := range validPositions {
		var correct = true
		for _, p := range v {
			if sapling, ok := w.Block(p).(Sapling); ok {
				if sapling.Wood != s.Wood {
					correct = false
				}
			} else {
				correct = false
			}
		}
		if correct {
			var lowestX = 0
			var lowestZ = 0
			for _, p := range v {
				if p.X() < lowestX {
					lowestX = p.X()
				}
				if p.Z() < lowestZ {
					lowestZ = p.Z()
				}
			}
			return &cube.Pos{lowestX, pos.Y(), lowestZ}, true
		}
	}
	return nil, false
}

// Grow grows this sapling into a tree
func (s Sapling) Grow(pos cube.Pos, w *world.World) (success bool) {

	var treeName = "minecraft:"
	pos2, correct := s.findSaplings(pos, w)
	if correct {
		treeName += "large_"
		pos = *pos2
	}
	treeName += s.Wood.String()
	treeName += "_tree"

	if tree := world.GetFeature(treeName); tree != nil {
		tree.Place(pos, w)
		return true
	}
	return false
}

// RandomTick ...
func (s Sapling) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	if rand.Intn(16) == 1 {
		s.Grow(pos, w)
	}
}

// BoneMeal ...
func (s Sapling) BoneMeal(pos cube.Pos, w *world.World) (success bool) {
	s.Grow(pos, w)
	return success
}

// UseOnBlock ...
func (s Sapling) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, s)
	if !used {
		return false
	}
	if !supportsVegetation(s, w.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(w, pos, s, user, ctx)
	return placed(ctx)
}

// HasLiquidDrops ...
func (Sapling) HasLiquidDrops() bool {
	return true
}

// FlammabilityInfo ...
func (Sapling) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, false)
}

// BreakInfo ...
func (s Sapling) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(s))
}

// CompostChance ...
func (Sapling) CompostChance() float64 {
	return 0.3
}

// EncodeItem ...
func (s Sapling) EncodeItem() (name string, meta int16) {
	return "minecraft:sapling", int16(s.Wood.Uint8())
}

// EncodeBlock ...
func (s Sapling) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:sapling", map[string]any{
		"sapling_type": s.Wood.String(),
		"age_bit":      s.AgeBit,
	}
}

func allSapling() (saplings []world.Block) {
	for _, w := range WoodTypes() {
		if w == CrimsonWood() || w == WarpedWood() || w == Mangrove() {
			continue
		}
		saplings = append(saplings, Sapling{Wood: w})
	}
	return
}
