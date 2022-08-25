package block

import "github.com/df-mc/dragonfly/server/world"

// Moss is a natural block that can be spread to some other blocks by using bone meal.
type Moss struct {
	solid
}

// SoilFor ...
func (Moss) SoilFor(block world.Block) bool {
	switch block.(type) {
	case TallGrass, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, Azalea:
		return true
	}
	return false
}

// BreakInfo ...
func (m Moss) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(m))
}

// EncodeItem ...
func (Moss) EncodeItem() (name string, meta int16) {
	return "minecraft:moss_block", 0
}

// EncodeBlock ...
func (Moss) EncodeBlock() (string, map[string]any) {
	return "minecraft:moss_block", nil
}
