package block

import (
	"math/rand"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)

// CoralFan is a non-solid block that comes in 5 variants.
type CoralFan struct {
	empty
	transparent
	sourceWaterDisplacer

	// Type is the type of coral of the block.
	Type CoralType
	// Dead is whether the coral is dead.
	Dead bool
	// Hanging is whether this coral is hanging on a wall
	Hanging bool
	// Facing is the direction that the coral fan is facing.
	Facing cube.Direction
}

// UseOnBlock ...
func (c CoralFan) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) bool {
	pos, _, used := firstReplaceable(w, pos, face, c)
	if !used {
		return false
	}
	if face == cube.FaceDown {
		return false
	}
	if face != cube.FaceUp {
		c.Hanging = true
		c.Facing = cube.Direction(face - 2)
	}
	attach := pos.Side(face)
	if !w.Block(attach).Model().FaceSolid(attach, face.Opposite(), w) {
		return false
	}
	if liquid, ok := w.Liquid(pos); ok {
		if water, ok := liquid.(Water); ok {
			if water.Depth != 8 {
				return false
			}
		} else {
			return false
		}
	} else {
		c.Dead = true
	}

	place(w, pos, c, user, ctx)
	return placed(ctx)
}

// HasLiquidDrops ...
func (c CoralFan) HasLiquidDrops() bool {
	return false
}

// SideClosed ...
func (c CoralFan) SideClosed(cube.Pos, cube.Pos, *world.World) bool {
	return false
}

// NeighbourUpdateTick ...
func (c CoralFan) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
	var face cube.Face
	if c.Hanging {
		face = c.Facing.Face()
	} else {
		face = cube.FaceDown
	}
	attach := pos.Side(face)
	if !w.Block(attach).Model().FaceSolid(attach, face.Opposite(), w) {
		w.SetBlock(pos, nil, nil)
		w.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: c})
		return
	}
	if c.Dead {
		return
	}
	w.ScheduleBlockUpdate(pos, time.Second*5/2)
}

// ScheduledTick ...
func (c CoralFan) ScheduledTick(pos cube.Pos, w *world.World, _ *rand.Rand) {
	if c.Dead {
		return
	}

	water := false
	if liquid, ok := w.Liquid(pos); ok {
		if _, ok := liquid.(Water); ok {
			water = true
		}
	}
	if !water {
		c.Dead = true
		w.SetBlock(pos, c, nil)
	}
}

// BreakInfo ...
func (c CoralFan) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(c))
}

func (c CoralFan) name() (name string, coralHangTypeBit bool) {
	name = "minecraft:coral_fan"
	if c.Hanging {
		name += "_hang"
		switch c.Type {
		case TubeCoral():
			coralHangTypeBit = false
		case BrainCoral():
			coralHangTypeBit = true
		case BubbleCoral():
			name += "2"
			coralHangTypeBit = false
		case FireCoral():
			name += "2"
			coralHangTypeBit = true
		case HornCoral():
			name += "3"
			coralHangTypeBit = false
		}
	} else {
		if c.Dead {
			name += "_dead"
		}
	}
	return name, coralHangTypeBit
}

// EncodeBlock ...
func (c CoralFan) EncodeBlock() (name string, properties map[string]any) {
	properties = map[string]any{}
	name, hangBit := c.name()
	if c.Hanging {
		properties["coral_hang_type_bit"] = hangBit
		properties["coral_direction"] = int32(c.Facing)
		properties["dead_bit"] = c.Dead
	} else {
		properties["coral_color"] = c.Type.Colour().String()
		switch c.Facing {
		case cube.East, cube.West:
			properties["coral_fan_direction"] = int32(0)
		case cube.North, cube.South:
			properties["coral_fan_direction"] = int32(1)
		}

	}

	return name, properties
}

// EncodeItem ...
func (c CoralFan) EncodeItem() (name string, meta int16) {
	name, hangBit := c.name()
	if c.Hanging {
		if hangBit {
			meta |= 0x1
		}
		if c.Dead {
			meta |= int16(0x1 << 2)
		}
		meta |= int16(c.Facing) << 4
	} else {
		meta |= int16(c.Type.Uint8())
		switch c.Facing {
		case cube.North, cube.South:
			meta |= int16(1) << 8
		}
	}
	return name, meta
}

// allCoral returns a list of all coral block variants
func allCoralFan() (c []world.Block) {
	f := func(dead bool) {
		for _, t := range CoralTypes() {
			c = append(c, CoralFan{Type: t, Dead: dead})
			for _, d := range cube.Directions() {
				c = append(c, CoralFan{
					Type:    t,
					Dead:    dead,
					Hanging: true,
					Facing:  d,
				})
			}
		}
	}
	f(true)
	f(false)
	return
}
