package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Vines are climbable non-solid vegetation blocks that grow on walls.
type Vines struct {
	replaceable
	transparent
	empty

	// NorthDirection is true if the vines are attached towards north.
	NorthDirection bool
	// EastDirection is true if the vines are attached towards east.
	EastDirection bool
	// SouthDirection is true if the vines are attached towards south.
	SouthDirection bool
	// WestDirection is true if the vines are attached towards west.
	WestDirection bool
}

// FlammabilityInfo ...
func (Vines) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(15, 100, true)
}

// BreakInfo ...
func (v Vines) BreakInfo() BreakInfo {
	return newBreakInfo(0.2, alwaysHarvestable, func(t item.Tool) bool {
		return t.ToolType() == item.TypeShears || t.ToolType() == item.TypeAxe
	}, func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if t.ToolType() == item.TypeShears || hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(v, 1)}
		}
		return nil
	})
}

// EntityInside ...
func (Vines) EntityInside(_ cube.Pos, _ *world.World, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
}

// SetAttachment sets an attachment on the given cube.Direction.
func (v Vines) SetAttachment(direction cube.Direction, attached bool) Vines {
	switch direction {
	case cube.North:
		v.NorthDirection = attached
		return v
	case cube.East:
		v.EastDirection = attached
		return v
	case cube.South:
		v.SouthDirection = attached
		return v
	case cube.West:
		v.WestDirection = attached
		return v
	}
	panic("should never happen")
}

// Attachment returns the attachment of the vines at the given direction.
func (v Vines) Attachment(direction cube.Direction) bool {
	switch direction {
	case cube.North:
		return v.NorthDirection
	case cube.East:
		return v.EastDirection
	case cube.South:
		return v.SouthDirection
	case cube.West:
		return v.WestDirection
	}
	panic("should never happen")
}

// Attachments returns all attachments of the vines.
func (v Vines) Attachments() (attachments []cube.Direction) {
	for _, d := range cube.Directions() {
		if v.Attachment(d) {
			attachments = append(attachments, d)
		}
	}
	return
}

// UseOnBlock ...
func (v Vines) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	if _, ok := tx.Block(pos).Model().(model.Solid); !ok || face.Axis() == cube.Y {
		return false
	}
	pos, face, used := firstReplaceable(tx, pos, face, v)
	if !used {
		return false
	}
	//noinspection GoAssignmentToReceiver
	v = v.SetAttachment(face.Direction().Opposite(), true)

	place(tx, pos, v, user, ctx)
	return placed(ctx)
}

// NeighbourUpdateTick ...
func (v Vines) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	above, updated := tx.Block(pos.Side(cube.FaceUp)), false
	for _, d := range v.Attachments() {
		if _, ok := tx.Block(pos.Side(d.Face())).Model().(model.Solid); !ok {
			if o, ok := above.(Vines); !ok || ok && !o.Attachment(d) {
				//noinspection GoAssignmentToReceiver
				v = v.SetAttachment(d, false)
				updated = true
			}
		}
	}
	if !updated {
		return
	}
	if _, ok := above.Model().(model.Solid); !ok && len(v.Attachments()) == 0 {
		tx.SetBlock(pos, nil, nil)
		return
	}
	tx.SetBlock(pos, v, nil)
}

// RandomTick ...
func (v Vines) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	if r.Float64() > 0.25 {
		// Vines have a 25% chance of spreading.
		return
	}

	face := cube.Face(r.Intn(len(cube.Faces())))
	selectedPos := pos.Side(face)
	if face == cube.FaceUp || face == cube.FaceDown {
		if selectedPos.OutOfBounds(tx.Range()) {
			// Vines can't spread outside world bounds.
			return
		}

		_, air := tx.Block(selectedPos).(Air)
		newVines, vines := tx.Block(selectedPos).(Vines)
		if face == cube.FaceUp {
			if !air {
				// Vines can't spread upwards unless there is air.
				return
			}

			for _, f := range cube.HorizontalFaces() {
				if r.Intn(2) == 0 && v.acceptableNeighbour(tx, selectedPos, f.Opposite()) {
					newVines = newVines.SetAttachment(f.Direction(), true)
				}
			}
		} else if air {
			for _, f := range cube.HorizontalFaces() {
				if r.Intn(2) == 0 && v.Attachment(f.Direction()) {
					newVines = newVines.SetAttachment(f.Direction(), true)
				}
			}
		} else if vines {
			var changed bool
			for _, f := range cube.HorizontalFaces() {
				if r.Intn(2) == 0 && v.Attachment(f.Direction()) {
					newVines, changed = newVines.SetAttachment(f.Direction(), true), true
				}
			}
			if changed {
				tx.SetBlock(selectedPos, newVines, nil)
			}
			return
		}

		if len(newVines.Attachments()) > 0 {
			tx.SetBlock(selectedPos, newVines, nil)
		}
	} else if !v.Attachment(face.Direction()) && v.canSpread(tx, pos) {
		if _, ok := tx.Block(selectedPos).(Air); ok {
			rightRotatedFace := face.RotateRight()
			leftRotatedFace := face.RotateLeft()

			attachedOnRight := v.Attachment(rightRotatedFace.Direction())
			attachedOnLeft := v.Attachment(leftRotatedFace.Direction())

			rightSelectedPos := selectedPos.Side(rightRotatedFace)
			leftSelectedPos := selectedPos.Side(leftRotatedFace)

			if attachedOnRight && v.acceptableNeighbour(tx, rightSelectedPos.Side(rightRotatedFace), rightRotatedFace) {
				tx.SetBlock(selectedPos, (Vines{}).SetAttachment(rightRotatedFace.Direction(), true), nil)
			} else if attachedOnLeft && v.acceptableNeighbour(tx, leftSelectedPos.Side(leftRotatedFace), leftRotatedFace) {
				tx.SetBlock(selectedPos, (Vines{}).SetAttachment(leftRotatedFace.Direction(), true), nil)
			} else if _, ok = tx.Block(rightSelectedPos).(Air); ok && attachedOnRight && v.acceptableNeighbour(tx, rightSelectedPos, face) {
				tx.SetBlock(rightSelectedPos, (Vines{}).SetAttachment(face.Opposite().Direction(), true), nil)
			} else if _, ok = tx.Block(leftSelectedPos).(Air); ok && attachedOnLeft && v.acceptableNeighbour(tx, leftSelectedPos, face) {
				tx.SetBlock(leftSelectedPos, (Vines{}).SetAttachment(face.Opposite().Direction(), true), nil)
			}
		} else if _, ok = tx.Block(selectedPos).Model().(model.Solid); ok {
			tx.SetBlock(pos, v.SetAttachment(face.Direction(), true), nil)
		}
	}
}

// EncodeItem ...
func (Vines) EncodeItem() (name string, meta int16) {
	return "minecraft:vine", 0
}

// EncodeBlock ...
func (v Vines) EncodeBlock() (string, map[string]any) {
	var bits int
	for i, ok := range []bool{v.SouthDirection, v.WestDirection, v.NorthDirection, v.EastDirection} {
		if ok {
			bits |= 1 << i
		}
	}
	return "minecraft:vine", map[string]any{"vine_direction_bits": int32(bits)}
}

// acceptableNeighbour returns true if the block at the given position is an acceptable neighbour.
func (v Vines) acceptableNeighbour(tx *world.Tx, pos cube.Pos, face cube.Face) bool {
	above := tx.Block(pos.Side(cube.FaceUp))
	_, air := above.(Air)
	_, vines := above.(Vines)
	return v.canSpreadTo(tx.Block(pos.Side(face.Opposite()))) && (air || vines || v.canSpreadTo(above))
}

// canSpreadTo returns true if the vines can spread to the given position.
func (Vines) canSpreadTo(b world.Block) bool {
	// TODO: Account for pistons, cauldrons, and shulker boxes.
	switch b.(type) {
	case Glass, StainedGlass:
		return false
	case Beacon:
		return false
	case WoodTrapdoor:
		return false
	}

	_, ok := b.Model().(model.Solid)
	return ok
}

// canSpread returns true if the vines can spread from the given position.
func (v Vines) canSpread(tx *world.Tx, pos cube.Pos) bool {
	var count int
	for x := -4; x <= 4; x++ {
		for z := -4; z <= 4; z++ {
			for y := -1; y <= 1; y++ {
				if _, ok := tx.Block(pos.Add(cube.Pos{x, y, z})).(Vines); ok {
					count++
					if count >= 5 {
						return false
					}
				}
			}
		}
	}
	return true
}

// allVines ...
func allVines() (b []world.Block) {
	for _, north := range []bool{true, false} {
		for _, east := range []bool{true, false} {
			for _, south := range []bool{true, false} {
				for _, west := range []bool{true, false} {
					b = append(b, Vines{
						NorthDirection: north,
						EastDirection:  east,
						SouthDirection: south,
						WestDirection:  west,
					})
				}
			}
		}
	}
	return
}
