package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/features"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// Azalea are a flower like
type Azalea struct {
	solid
	transparent

	// Flowering specifies if this block is flowering.
	Flowering bool
}

// BoneMeal ...
func (a Azalea) BoneMeal(pos cube.Pos, w *world.World) (success bool) {
	if _, ok := w.Block(pos.Side(cube.FaceDown)).(Moss); ok {
		for i := 0; i < 8; i++ {
			p := pos.Add(cube.Pos{rand.Intn(7) - 3, rand.Intn(3) - 1, rand.Intn(7) - 3})
			if _, ok := w.Block(p).(Air); !ok {
				continue
			}
			if _, ok := w.Block(p.Side(cube.FaceDown)).(Moss); !ok {
				continue
			}
			w.SetBlock(p, Azalea{}, nil)
			success = true
		}
	} else {
		t := features.AzaleaTree{
			Trunk:   Log{Wood: OakWood()},
			Leaves:  AzaleaLeaves{},
			Leaves2: AzaleaLeaves{Flowering: true},
		}
		return t.GrowTree(pos, w)
	}
	return success
}

// NeighbourUpdateTick ...
func (a Azalea) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !supportsVegetation(a, w.Block(pos.Side(cube.FaceDown))) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: a})
		dropItem(w, item.NewStack(a, 1), pos.Vec3Centre())
	}
}

// UseOnBlock ...
func (a Azalea) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, a)
	if !used {
		return false
	}
	down := w.Block(pos.Side(cube.FaceDown))
	if !supportsVegetation(a, down) {
		return false
	}

	place(w, pos, a, user, ctx)
	return placed(ctx)
}

// HasLiquidDrops ...
func (Azalea) HasLiquidDrops() bool {
	return true
}

// FlammabilityInfo ...
func (a Azalea) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(60, 100, false)
}

// BreakInfo ...
func (a Azalea) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(a))
}

// CompostChance ...
func (Azalea) CompostChance() float64 {
	return 0.65
}

// EncodeItem ...
func (a Azalea) EncodeItem() (name string, meta int16) {
	name = "minecraft:"
	if a.Flowering {
		name += "flowering_"
	}
	name += "azalea"
	return name, 0
}

// EncodeBlock ...
func (a Azalea) EncodeBlock() (name string, properties map[string]any) {
	name = "minecraft:"
	if a.Flowering {
		name += "flowering_"
	}
	name += "azalea"
	return name, map[string]any{}
}

// allLogs returns a list of all possible leaves states.
func allAzalea() (leaves []world.Block) {
	f := func(flowering bool) {
		leaves = append(leaves, Azalea{Flowering: flowering})
	}
	f(true)
	f(false)
	return leaves
}
