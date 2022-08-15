package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// Seagrass is a non-solid plant block that generates in all oceans, except frozen oceans.
type Seagrass struct {
	empty
	transparent
	sourceWaterDisplacer

	// Type is the type of seagrass
	Type SeagrassType
}

// HasLiquidDrops ...
func (Seagrass) HasLiquidDrops() bool {
	return true
}

// SideClosed ...
func (Seagrass) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// BreakInfo ...
func (s Seagrass) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, simpleDrops(item.NewStack(s, 1)))
}

// FlammabilityInfo ...
func (Seagrass) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, false)
}

// SmeltInfo ...
func (Seagrass) SmeltInfo() item.SmeltInfo {
	return newSmeltInfo(item.NewStack(item.Dye{Colour: item.ColourLime()}, 1), 0.1)
}

// EncodeItem ...
func (Seagrass) EncodeItem() (name string, meta int16) {
	return "minecraft:seagrass", 0
}

// EncodeBlock ...
func (s Seagrass) EncodeBlock() (string, map[string]any) {
	return "minecraft:seagrass", map[string]any{
		"sea_grass_type": s.Type.String(),
	}
}

func allSeagrass() (c []world.Block) {
	for _, t := range SeagrassTypes() {
		c = append(c, Seagrass{Type: t})
	}
	return
}
