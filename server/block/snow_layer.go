package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// SnowLayer is a layer of snow.
type SnowLayer struct {
	solid

	// Height is how many extra layers this block has.
	Height int
	// Covered is if this block is placed on a leaf block.
	Covered bool
}

// BreakInfo ...
func (s SnowLayer) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, shovelEffective, silkTouchDrop(item.NewStack(item.Snowball{}, 4), item.NewStack(s, 1)))
}

// UseOnBlock ...
func (s SnowLayer) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	if existing, ok := w.Block(pos).(SnowLayer); ok {
		if existing.Height >= 7 {
			return false
		}
		existing.Height++
		place(w, pos, existing, user, ctx)
		return placed(ctx)
	}

	pos, _, used := firstReplaceable(w, pos, face, s)
	if !used {
		return false
	}

	below := w.Block(pos.Side(cube.FaceDown))
	_, leaves := below.(Leaves)

	if !leaves && !below.Model().FaceSolid(pos, face, w) {
		return false
	}

	if leaves {
		s.Covered = true
	}

	place(w, pos, s, user, ctx)
	return placed(ctx)
}

func (SnowLayer) canSurvive(pos cube.Pos, w *world.World) bool {
	below := w.Block(pos.Side(cube.FaceDown))
	_, leaves := below.(Leaves)
	if !leaves && !below.Model().FaceSolid(pos.Side(cube.FaceDown), cube.FaceUp, w) {
		return false
	}
	return true
}

// NeighbourUpdateTick ...
func (s SnowLayer) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	if !s.canSurvive(pos, w) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: s})
		return
	}
}

func (s SnowLayer) melt(pos cube.Pos, w *world.World) {
	if s.Height == 0 {
		w.SetBlock(pos, nil, nil)
		return
	}
	s.Height -= 1
	w.SetBlock(pos, s, nil)
}

// RandomTick ...
func (s SnowLayer) RandomTick(pos cube.Pos, w *world.World, r *rand.Rand) {
	pos = pos.Side(cube.FaceUp)
	if w.Light(pos) > 12 {
		s.melt(pos, w)
		return
	}
	if w.Biome(pos).Temperature() >= 1 {
		s.melt(pos, w)
	}
}

// EncodeItem ...
func (SnowLayer) EncodeItem() (name string, meta int16) {
	return "minecraft:snow_layer", 0
}

// EncodeBlock ...
func (s SnowLayer) EncodeBlock() (string, map[string]any) {
	return "minecraft:snow_layer", map[string]any{
		"height":      int32(s.Height),
		"covered_bit": s.Covered,
	}
}

// allSnowLayers ...
func allSnowLayers() (b []world.Block) {
	for i := 0; i <= 7; i++ {
		b = append(b, SnowLayer{Height: i})
		b = append(b, SnowLayer{Height: i, Covered: true})
	}
	return
}
