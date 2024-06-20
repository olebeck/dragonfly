package chunk

type BlockRegistry interface {
	AirRuntimeID() uint32
	RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool)
	StateToRuntimeID(name string, properties map[string]any) (runtimeID uint32, found bool)
	FilteringBlocks() []uint8
	LightBlocks() []uint8
	RandomTickBlock(rid uint32) bool
	IsWater(rid uint32) bool
	NBTBlock(rid uint32) bool
	LiquidDisplacingBlock(rid uint32) bool
	LiquidBlock(rid uint32) bool
	HashToRuntimeID(hash uint32) (rid uint32, ok bool)
}
