package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Sapling is a non-solid block that can be grown into a tree.
type Sapling struct {
	empty
	transparent

	Wood   WoodType
	AgeBit bool
}

// UseOnBlock ...
func (s Sapling) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(tx, pos, face, s)
	if !used {
		return false
	}
	if !supportsVegetation(s, tx.Block(pos.Side(cube.FaceDown))) {
		return false
	}

	place(tx, pos, s, user, ctx)
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
		if w == CrimsonWood() || w == WarpedWood() || w == MangroveWood() || w == CherryWood() {
			continue
		}
		saplings = append(saplings, Sapling{Wood: w})
	}
	return
}
