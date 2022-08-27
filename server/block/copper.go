package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// Copper is a block that can be waxed and can weather.
type Copper struct {
	solid
	// Waxed is if this copper block has been waxed and wont weather.
	Waxed bool
	// Cut is if this copper block is the cut variant.
	Cut bool
	// Weather is the current weathering state.
	Weather WeatheringType
}

// BreakInfo ...
func (c Copper) BreakInfo() BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, oneOf(c))
}

func (c Copper) name(isBlock bool) (name string) {
	if c.Waxed {
		name += "waxed_"
	}

	if c.Weather != NotWeathered() {
		name += c.Weather.String() + "_"
	}

	if c.Cut {
		name += "cut_"
	}
	name += "copper"

	if !c.Cut && c.Weather == NotWeathered() && !c.Waxed && isBlock {
		name += "_block"
	}
	return name
}

// EncodeItem ...
func (c Copper) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.name(true), 0
}

// EncodeBlock ...
func (c Copper) EncodeBlock() (name string, data map[string]any) {
	// to stop blockhash ignoring these:
	// Waxed
	// Cut
	// Weather
	return "minecraft:" + c.name(true), nil
}

// RandomTick ...
func (c Copper) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	// TODO: weathering
}

func allCopper() (c []world.Block) {
	f := func(cut, waxed bool) {
		for _, w := range WeatheringTypes() {
			c = append(c, Copper{Waxed: waxed, Cut: cut, Weather: w})
		}
	}
	f(true, true)
	f(true, false)
	f(false, true)
	f(false, false)
	return
}
