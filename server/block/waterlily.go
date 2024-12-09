package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Waterlily is a block that can be placed on water.
type Waterlily struct {
	carpet
	transparent
}

// FlammabilityInfo ...
func (c Waterlily) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, false)
}

// SideClosed ...
func (Waterlily) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// BreakInfo ...
func (c Waterlily) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(c))
}

// EncodeItem ...
func (c Waterlily) EncodeItem() (name string, meta int16) {
	return "minecraft:waterlily", 0
}

// EncodeBlock ...
func (c Waterlily) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:waterlily", map[string]any{}
}

// HasLiquidDrops ...
func (Waterlily) HasLiquidDrops() bool {
	return true
}

// NeighbourUpdateTick ...
func (c Waterlily) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	under := tx.Block(pos.Side(cube.FaceDown))
	if water, ok := under.(Water); !ok || water.LiquidDepth() < 8 {
		tx.SetBlock(pos, nil, nil)
		dropItem(tx, item.NewStack(c, 1), pos.Vec3Centre())
	}
}

// UseOnBlock handles only placing waterlilies on water
func (c Waterlily) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(tx, pos, face, c)
	if !used {
		return
	}

	if water, ok := tx.Block(pos).(Water); !ok || water.LiquidDepth() < 8 {
		return
	}
	place(tx, pos.Side(cube.FaceUp), c, user, ctx)
	return placed(ctx)
}
