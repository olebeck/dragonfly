package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/internal/nbtconv"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// NewSplashPotion creates a splash potion. SplashPotion is an item that grants
// effects when thrown.
func NewSplashPotion(pos mgl64.Vec3, owner world.Entity, t potion.Potion) *Ent {
	colour, _ := effect.ResultingColour(t.Effects())

	conf := splashPotionConf
	conf.Potion = t
	conf.Particle = particle.Splash{Colour: colour}
	conf.Hit = potionSplash(1, t, false)
	return Config{Behaviour: conf.New(owner)}.New(SplashPotionType{}, pos)
}

var splashPotionConf = ProjectileBehaviourConfig{
	Gravity: 0.05,
	Drag:    0.01,
	Damage:  -1,
	Sound:   sound.GlassBreak{},
}

// SplashPotionType is a world.EntityType implementation for SplashPotion.
type SplashPotionType struct{}

func (SplashPotionType) EncodeEntity() string { return "minecraft:splash_potion" }
func (SplashPotionType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.125, 0, -0.125, 0.125, 0.25, 0.125)
}

func (SplashPotionType) DecodeNBT(m map[string]any) world.Entity {
	pot := NewSplashPotion(nbtconv.Vec3(m, "Pos"), nil, potion.From(nbtconv.Int32(m, "PotionId")))
	pot.vel = nbtconv.Vec3(m, "Motion")
	if uniqueID, ok := m["UniqueID"].(int64); ok {
		pot.uniqueID = uniqueID
	}
	return pot
}

func (SplashPotionType) EncodeNBT(e world.Entity) map[string]any {
	pot := e.(*Ent)
	return map[string]any{
		"UniqueID": pot.uniqueID,
		"Pos":      nbtconv.Vec3ToFloat32Slice(pot.Position()),
		"Motion":   nbtconv.Vec3ToFloat32Slice(pot.Velocity()),
		"PotionId": int32(pot.conf.Behaviour.(*ProjectileBehaviour).conf.Potion.Uint8()),
	}
}
