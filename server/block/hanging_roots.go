package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// HangingRoots is a natural decorative block found underground in the lush caves biome.
type HangingRoots struct {
	empty
}

// BreakInfo ...
func (h HangingRoots) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, func(t item.Tool, e []item.Enchantment) (s []item.Stack) {
		if _, ok := t.(item.Shears); ok || hasSilkTouch(e) {
			s = append(s, item.NewStack(HangingRoots{}, 1))
		}
		return s
	})
}

func (h HangingRoots) canSurvive(pos cube.Pos, w *world.World) bool {
	above := pos.Side(cube.FaceUp)
	if !w.Block(above).Model().FaceSolid(above, cube.FaceDown, w) {
		return false
	}
	return true
}

// NeighbourUpdateTick ...
func (h HangingRoots) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !h.canSurvive(pos, w) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: h})
		dropItem(w, item.NewStack(h, 1), pos.Vec3Centre())
	}
}

// UseOnBlock ...
func (h HangingRoots) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, h)
	if !used {
		return false
	}
	if !h.canSurvive(pos, w) {
		return false
	}

	place(w, pos, h, user, ctx)
	return placed(ctx)
}

// EncodeItem ...
func (h HangingRoots) EncodeItem() (name string, meta int16) {
	return "minecraft:hanging_roots", 0
}

// EncodeBlock ...
func (h HangingRoots) EncodeBlock() (string, map[string]any) {
	return "minecraft:hanging_roots", nil
}
