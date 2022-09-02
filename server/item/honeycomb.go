package item

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// Honeycomb is an item obtained from bee nests and beehives.
type Honeycomb struct{}

// waxable represents a block that can be waxed.
type waxable interface {
	// Strip returns a block that is the result of waxing it. Alternatively, the bool returned may be false to
	// indicate the block couldn't be waxed.
	Wax() (world.Block, bool)
}

// EncodeItem ...
func (Honeycomb) EncodeItem() (name string, meta int16) {
	return "minecraft:honeycomb", 0
}

// UseOnBlock handles the waxing copper when a player clicks a log with honeycomb.
func (h Honeycomb) UseOnBlock(pos cube.Pos, _ cube.Face, _ mgl64.Vec3, w *world.World, _ User, ctx *UseContext) bool {
	if s, ok := w.Block(pos).(waxable); ok {
		if res, ok := s.Wax(); ok {
			w.SetBlock(pos, res, nil)
			w.PlaySound(pos.Vec3(), sound.ItemUseOn{Block: res})

			ctx.SubtractFromCount(1)
			return true
		}
	}
	return false
}
