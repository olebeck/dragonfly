// Code generated by cmd/blockhash; DO NOT EDIT.

package block

const (
	hashAir = iota
	hashAmethyst
	hashAncientDebris
	hashAndesite
	hashAnvil
	hashAzalea
	hashAzaleaLeaves
	hashBanner
	hashBarrel
	hashBarrier
	hashBasalt
	hashBeacon
	hashBedrock
	hashBeetrootSeeds
	hashBlackstone
	hashBlastFurnace
	hashBlueIce
	hashBone
	hashBookshelf
	hashBricks
	hashCactus
	hashCake
	hashCalcite
	hashCarpet
	hashCarrot
	hashChain
	hashChest
	hashChiseledQuartz
	hashClay
	hashCoal
	hashCoalOre
	hashCobblestone
	hashCocoaBean
	hashComposter
	hashConcrete
	hashConcretePowder
	hashCopperOre
	hashCoral
	hashCoralBlock
	hashCraftingTable
	hashDeadBush
	hashDeepslate
	hashDeepslateBricks
	hashDeepslateTiles
	hashDiamond
	hashDiamondOre
	hashDiorite
	hashDirt
	hashDirtPath
	hashDoubleFlower
	hashDoubleTallGrass
	hashDragonEgg
	hashDriedKelp
	hashDripstone
	hashEmerald
	hashEmeraldOre
	hashEnchantingTable
	hashEndBricks
	hashEndStone
	hashEnderChest
	hashFarmland
	hashFire
	hashFletchingTable
	hashFlower
	hashFroglight
	hashFurnace
	hashGlass
	hashGlassPane
	hashGlazedTerracotta
	hashGlowstone
	hashGold
	hashGoldOre
	hashGranite
	hashGrass
	hashGravel
	hashGrindstone
	hashHayBale
	hashHoneycomb
	hashInvisibleBedrock
	hashIron
	hashIronBars
	hashIronOre
	hashItemFrame
	hashJukebox
	hashKelp
	hashLadder
	hashLantern
	hashLapis
	hashLapisOre
	hashLava
	hashLeaves
	hashLight
	hashLitPumpkin
	hashLog
	hashLoom
	hashMelon
	hashMelonSeeds
	hashMoss
	hashMossCarpet
	hashMud
	hashMudBricks
	hashMuddyMangroveRoots
	hashMushroomBlock
	hashNetherBrickFence
	hashNetherBricks
	hashNetherGoldOre
	hashNetherQuartzOre
	hashNetherSprouts
	hashNetherWart
	hashNetherWartBlock
	hashNetherite
	hashNetherrack
	hashNote
	hashObsidian
	hashPackedIce
	hashPackedMud
	hashPlanks
	hashPodzol
	hashPolishedBlackstoneBrick
	hashPotato
	hashPrismarine
	hashPumpkin
	hashPumpkinSeeds
	hashPurpur
	hashPurpurPillar
	hashQuartz
	hashQuartzBricks
	hashQuartzPillar
	hashRawCopper
	hashRawGold
	hashRawIron
	hashReinforcedDeepslate
	hashSand
	hashSandstone
	hashSeaLantern
	hashSeaPickle
	hashSeagrass
	hashShroomlight
	hashSign
	hashSkull
	hashSlab
	hashSmithingTable
	hashSmoker
	hashSnow
	hashSoulSand
	hashSoulSoil
	hashSponge
	hashSporeBlossom
	hashStainedGlass
	hashStainedGlassPane
	hashStainedTerracotta
	hashStairs
	hashStone
	hashStoneBricks
	hashStonecutter
	hashSugarCane
	hashTNT
	hashTallGrass
	hashTerracotta
	hashTorch
	hashTuff
	hashWall
	hashWater
	hashWaterlily
	hashWheatSeeds
	hashWood
	hashWoodDoor
	hashWoodFence
	hashWoodFenceGate
	hashWoodTrapdoor
	hashWool
)

func (Air) Hash() uint64 {
	return hashAir
}

func (Amethyst) Hash() uint64 {
	return hashAmethyst
}

func (AncientDebris) Hash() uint64 {
	return hashAncientDebris
}

func (a Andesite) Hash() uint64 {
	return hashAndesite | uint64(boolByte(a.Polished))<<8
}

func (a Anvil) Hash() uint64 {
	return hashAnvil | uint64(a.Type.Uint8())<<8 | uint64(a.Facing)<<10
}

func (a Azalea) Hash() uint64 {
	return hashAzalea | uint64(boolByte(a.Flowering))<<8
}

func (l AzaleaLeaves) Hash() uint64 {
	return hashAzaleaLeaves | uint64(boolByte(l.Flowering))<<8 | uint64(boolByte(l.Persistent))<<9 | uint64(boolByte(l.ShouldUpdate))<<10
}

func (b Banner) Hash() uint64 {
	return hashBanner | uint64(b.Attach.Uint8())<<8
}

func (b Barrel) Hash() uint64 {
	return hashBarrel | uint64(b.Facing)<<8 | uint64(boolByte(b.Open))<<11
}

func (Barrier) Hash() uint64 {
	return hashBarrier
}

func (b Basalt) Hash() uint64 {
	return hashBasalt | uint64(boolByte(b.Polished))<<8 | uint64(b.Axis)<<9
}

func (Beacon) Hash() uint64 {
	return hashBeacon
}

func (b Bedrock) Hash() uint64 {
	return hashBedrock | uint64(boolByte(b.InfiniteBurning))<<8
}

func (b BeetrootSeeds) Hash() uint64 {
	return hashBeetrootSeeds | uint64(b.Growth)<<8
}

func (b Blackstone) Hash() uint64 {
	return hashBlackstone | uint64(b.Type.Uint8())<<8
}

func (b BlastFurnace) Hash() uint64 {
	return hashBlastFurnace | uint64(b.Facing)<<8 | uint64(boolByte(b.Lit))<<11
}

func (BlueIce) Hash() uint64 {
	return hashBlueIce
}

func (b Bone) Hash() uint64 {
	return hashBone | uint64(b.Axis)<<8
}

func (Bookshelf) Hash() uint64 {
	return hashBookshelf
}

func (Bricks) Hash() uint64 {
	return hashBricks
}

func (c Cactus) Hash() uint64 {
	return hashCactus | uint64(c.Age)<<8
}

func (c Cake) Hash() uint64 {
	return hashCake | uint64(c.Bites)<<8
}

func (Calcite) Hash() uint64 {
	return hashCalcite
}

func (c Carpet) Hash() uint64 {
	return hashCarpet | uint64(c.Colour.Uint8())<<8
}

func (c Carrot) Hash() uint64 {
	return hashCarrot | uint64(c.Growth)<<8
}

func (c Chain) Hash() uint64 {
	return hashChain | uint64(c.Axis)<<8
}

func (c Chest) Hash() uint64 {
	return hashChest | uint64(c.Facing)<<8
}

func (ChiseledQuartz) Hash() uint64 {
	return hashChiseledQuartz
}

func (Clay) Hash() uint64 {
	return hashClay
}

func (Coal) Hash() uint64 {
	return hashCoal
}

func (c CoalOre) Hash() uint64 {
	return hashCoalOre | uint64(c.Type.Uint8())<<8
}

func (c Cobblestone) Hash() uint64 {
	return hashCobblestone | uint64(boolByte(c.Mossy))<<8
}

func (c CocoaBean) Hash() uint64 {
	return hashCocoaBean | uint64(c.Facing)<<8 | uint64(c.Age)<<10
}

func (c Composter) Hash() uint64 {
	return hashComposter | uint64(c.Level)<<8
}

func (c Concrete) Hash() uint64 {
	return hashConcrete | uint64(c.Colour.Uint8())<<8
}

func (c ConcretePowder) Hash() uint64 {
	return hashConcretePowder | uint64(c.Colour.Uint8())<<8
}

func (c CopperOre) Hash() uint64 {
	return hashCopperOre | uint64(c.Type.Uint8())<<8
}

func (c Coral) Hash() uint64 {
	return hashCoral | uint64(c.Type.Uint8())<<8 | uint64(boolByte(c.Dead))<<11
}

func (c CoralBlock) Hash() uint64 {
	return hashCoralBlock | uint64(c.Type.Uint8())<<8 | uint64(boolByte(c.Dead))<<11
}

func (CraftingTable) Hash() uint64 {
	return hashCraftingTable
}

func (DeadBush) Hash() uint64 {
	return hashDeadBush
}

func (d Deepslate) Hash() uint64 {
	return hashDeepslate | uint64(d.Type.Uint8())<<8 | uint64(d.Axis)<<10
}

func (d DeepslateBricks) Hash() uint64 {
	return hashDeepslateBricks | uint64(boolByte(d.Cracked))<<8
}

func (d DeepslateTiles) Hash() uint64 {
	return hashDeepslateTiles | uint64(boolByte(d.Cracked))<<8
}

func (Diamond) Hash() uint64 {
	return hashDiamond
}

func (d DiamondOre) Hash() uint64 {
	return hashDiamondOre | uint64(d.Type.Uint8())<<8
}

func (d Diorite) Hash() uint64 {
	return hashDiorite | uint64(boolByte(d.Polished))<<8
}

func (d Dirt) Hash() uint64 {
	return hashDirt | uint64(boolByte(d.Coarse))<<8
}

func (DirtPath) Hash() uint64 {
	return hashDirtPath
}

func (d DoubleFlower) Hash() uint64 {
	return hashDoubleFlower | uint64(boolByte(d.UpperPart))<<8 | uint64(d.Type.Uint8())<<9
}

func (d DoubleTallGrass) Hash() uint64 {
	return hashDoubleTallGrass | uint64(boolByte(d.UpperPart))<<8 | uint64(d.Type.Uint8())<<9
}

func (DragonEgg) Hash() uint64 {
	return hashDragonEgg
}

func (DriedKelp) Hash() uint64 {
	return hashDriedKelp
}

func (Dripstone) Hash() uint64 {
	return hashDripstone
}

func (Emerald) Hash() uint64 {
	return hashEmerald
}

func (e EmeraldOre) Hash() uint64 {
	return hashEmeraldOre | uint64(e.Type.Uint8())<<8
}

func (EnchantingTable) Hash() uint64 {
	return hashEnchantingTable
}

func (EndBricks) Hash() uint64 {
	return hashEndBricks
}

func (EndStone) Hash() uint64 {
	return hashEndStone
}

func (c EnderChest) Hash() uint64 {
	return hashEnderChest | uint64(c.Facing)<<8
}

func (f Farmland) Hash() uint64 {
	return hashFarmland | uint64(f.Hydration)<<8
}

func (f Fire) Hash() uint64 {
	return hashFire | uint64(f.Type.Uint8())<<8 | uint64(f.Age)<<9
}

func (FletchingTable) Hash() uint64 {
	return hashFletchingTable
}

func (f Flower) Hash() uint64 {
	return hashFlower | uint64(f.Type.Uint8())<<8
}

func (f Froglight) Hash() uint64 {
	return hashFroglight | uint64(f.Type.Uint8())<<8 | uint64(f.Axis)<<10
}

func (f Furnace) Hash() uint64 {
	return hashFurnace | uint64(f.Facing)<<8 | uint64(boolByte(f.Lit))<<11
}

func (Glass) Hash() uint64 {
	return hashGlass
}

func (GlassPane) Hash() uint64 {
	return hashGlassPane
}

func (t GlazedTerracotta) Hash() uint64 {
	return hashGlazedTerracotta | uint64(t.Colour.Uint8())<<8 | uint64(t.Facing)<<12
}

func (Glowstone) Hash() uint64 {
	return hashGlowstone
}

func (Gold) Hash() uint64 {
	return hashGold
}

func (g GoldOre) Hash() uint64 {
	return hashGoldOre | uint64(g.Type.Uint8())<<8
}

func (g Granite) Hash() uint64 {
	return hashGranite | uint64(boolByte(g.Polished))<<8
}

func (Grass) Hash() uint64 {
	return hashGrass
}

func (Gravel) Hash() uint64 {
	return hashGravel
}

func (g Grindstone) Hash() uint64 {
	return hashGrindstone | uint64(g.Attach.Uint8())<<8 | uint64(g.Facing)<<10
}

func (h HayBale) Hash() uint64 {
	return hashHayBale | uint64(h.Axis)<<8
}

func (Honeycomb) Hash() uint64 {
	return hashHoneycomb
}

func (InvisibleBedrock) Hash() uint64 {
	return hashInvisibleBedrock
}

func (Iron) Hash() uint64 {
	return hashIron
}

func (IronBars) Hash() uint64 {
	return hashIronBars
}

func (i IronOre) Hash() uint64 {
	return hashIronOre | uint64(i.Type.Uint8())<<8
}

func (i ItemFrame) Hash() uint64 {
	return hashItemFrame | uint64(i.Facing)<<8 | uint64(boolByte(i.Glowing))<<11
}

func (Jukebox) Hash() uint64 {
	return hashJukebox
}

func (k Kelp) Hash() uint64 {
	return hashKelp | uint64(k.Age)<<8
}

func (l Ladder) Hash() uint64 {
	return hashLadder | uint64(l.Facing)<<8
}

func (l Lantern) Hash() uint64 {
	return hashLantern | uint64(boolByte(l.Hanging))<<8 | uint64(l.Type.Uint8())<<9
}

func (Lapis) Hash() uint64 {
	return hashLapis
}

func (l LapisOre) Hash() uint64 {
	return hashLapisOre | uint64(l.Type.Uint8())<<8
}

func (l Lava) Hash() uint64 {
	return hashLava | uint64(boolByte(l.Still))<<8 | uint64(l.Depth)<<9 | uint64(boolByte(l.Falling))<<17
}

func (l Leaves) Hash() uint64 {
	return hashLeaves | uint64(l.Wood.Uint8())<<8 | uint64(boolByte(l.Persistent))<<12 | uint64(boolByte(l.ShouldUpdate))<<13
}

func (l Light) Hash() uint64 {
	return hashLight | uint64(l.Level)<<8
}

func (l LitPumpkin) Hash() uint64 {
	return hashLitPumpkin | uint64(l.Facing)<<8
}

func (l Log) Hash() uint64 {
	return hashLog | uint64(l.Wood.Uint8())<<8 | uint64(boolByte(l.Stripped))<<12 | uint64(l.Axis)<<13
}

func (l Loom) Hash() uint64 {
	return hashLoom | uint64(l.Facing)<<8
}

func (Melon) Hash() uint64 {
	return hashMelon
}

func (m MelonSeeds) Hash() uint64 {
	return hashMelonSeeds | uint64(m.Growth)<<8 | uint64(m.Direction)<<16
}

func (Moss) Hash() uint64 {
	return hashMoss
}

func (MossCarpet) Hash() uint64 {
	return hashMossCarpet
}

func (Mud) Hash() uint64 {
	return hashMud
}

func (MudBricks) Hash() uint64 {
	return hashMudBricks
}

func (m MuddyMangroveRoots) Hash() uint64 {
	return hashMuddyMangroveRoots | uint64(m.Axis)<<8
}

func (m MushroomBlock) Hash() uint64 {
	return hashMushroomBlock | uint64(m.Type.Uint8())<<8 | uint64(m.HugeBits)<<12
}

func (NetherBrickFence) Hash() uint64 {
	return hashNetherBrickFence
}

func (n NetherBricks) Hash() uint64 {
	return hashNetherBricks | uint64(n.Type.Uint8())<<8
}

func (NetherGoldOre) Hash() uint64 {
	return hashNetherGoldOre
}

func (NetherQuartzOre) Hash() uint64 {
	return hashNetherQuartzOre
}

func (NetherSprouts) Hash() uint64 {
	return hashNetherSprouts
}

func (n NetherWart) Hash() uint64 {
	return hashNetherWart | uint64(n.Age)<<8
}

func (n NetherWartBlock) Hash() uint64 {
	return hashNetherWartBlock | uint64(boolByte(n.Warped))<<8
}

func (Netherite) Hash() uint64 {
	return hashNetherite
}

func (Netherrack) Hash() uint64 {
	return hashNetherrack
}

func (Note) Hash() uint64 {
	return hashNote
}

func (o Obsidian) Hash() uint64 {
	return hashObsidian | uint64(boolByte(o.Crying))<<8
}

func (PackedIce) Hash() uint64 {
	return hashPackedIce
}

func (PackedMud) Hash() uint64 {
	return hashPackedMud
}

func (p Planks) Hash() uint64 {
	return hashPlanks | uint64(p.Wood.Uint8())<<8
}

func (Podzol) Hash() uint64 {
	return hashPodzol
}

func (b PolishedBlackstoneBrick) Hash() uint64 {
	return hashPolishedBlackstoneBrick | uint64(boolByte(b.Cracked))<<8
}

func (p Potato) Hash() uint64 {
	return hashPotato | uint64(p.Growth)<<8
}

func (p Prismarine) Hash() uint64 {
	return hashPrismarine | uint64(p.Type.Uint8())<<8
}

func (p Pumpkin) Hash() uint64 {
	return hashPumpkin | uint64(boolByte(p.Carved))<<8 | uint64(p.Facing)<<9
}

func (p PumpkinSeeds) Hash() uint64 {
	return hashPumpkinSeeds | uint64(p.Growth)<<8 | uint64(p.Direction)<<16
}

func (Purpur) Hash() uint64 {
	return hashPurpur
}

func (p PurpurPillar) Hash() uint64 {
	return hashPurpurPillar | uint64(p.Axis)<<8
}

func (q Quartz) Hash() uint64 {
	return hashQuartz | uint64(boolByte(q.Smooth))<<8
}

func (QuartzBricks) Hash() uint64 {
	return hashQuartzBricks
}

func (q QuartzPillar) Hash() uint64 {
	return hashQuartzPillar | uint64(q.Axis)<<8
}

func (RawCopper) Hash() uint64 {
	return hashRawCopper
}

func (RawGold) Hash() uint64 {
	return hashRawGold
}

func (RawIron) Hash() uint64 {
	return hashRawIron
}

func (ReinforcedDeepslate) Hash() uint64 {
	return hashReinforcedDeepslate
}

func (s Sand) Hash() uint64 {
	return hashSand | uint64(boolByte(s.Red))<<8
}

func (s Sandstone) Hash() uint64 {
	return hashSandstone | uint64(s.Type.Uint8())<<8 | uint64(boolByte(s.Red))<<10
}

func (SeaLantern) Hash() uint64 {
	return hashSeaLantern
}

func (s SeaPickle) Hash() uint64 {
	return hashSeaPickle | uint64(s.AdditionalCount)<<8 | uint64(boolByte(s.Dead))<<16
}

func (s Seagrass) Hash() uint64 {
	return hashSeagrass | uint64(s.Type.Uint8())<<8
}

func (Shroomlight) Hash() uint64 {
	return hashShroomlight
}

func (s Sign) Hash() uint64 {
	return hashSign | uint64(s.Wood.Uint8())<<8 | uint64(s.Attach.Uint8())<<12
}

func (s Skull) Hash() uint64 {
	return hashSkull | uint64(s.Attach.FaceUint8())<<8
}

func (s Slab) Hash() uint64 {
	return hashSlab | s.Block.Hash()<<8 | uint64(boolByte(s.Top))<<24 | uint64(boolByte(s.Double))<<25
}

func (SmithingTable) Hash() uint64 {
	return hashSmithingTable
}

func (s Smoker) Hash() uint64 {
	return hashSmoker | uint64(s.Facing)<<8 | uint64(boolByte(s.Lit))<<11
}

func (Snow) Hash() uint64 {
	return hashSnow
}

func (SoulSand) Hash() uint64 {
	return hashSoulSand
}

func (SoulSoil) Hash() uint64 {
	return hashSoulSoil
}

func (s Sponge) Hash() uint64 {
	return hashSponge | uint64(boolByte(s.Wet))<<8
}

func (SporeBlossom) Hash() uint64 {
	return hashSporeBlossom
}

func (g StainedGlass) Hash() uint64 {
	return hashStainedGlass | uint64(g.Colour.Uint8())<<8
}

func (p StainedGlassPane) Hash() uint64 {
	return hashStainedGlassPane | uint64(p.Colour.Uint8())<<8
}

func (t StainedTerracotta) Hash() uint64 {
	return hashStainedTerracotta | uint64(t.Colour.Uint8())<<8
}

func (s Stairs) Hash() uint64 {
	return hashStairs | s.Block.Hash()<<8 | uint64(boolByte(s.UpsideDown))<<24 | uint64(s.Facing)<<25
}

func (s Stone) Hash() uint64 {
	return hashStone | uint64(boolByte(s.Smooth))<<8
}

func (s StoneBricks) Hash() uint64 {
	return hashStoneBricks | uint64(s.Type.Uint8())<<8
}

func (s Stonecutter) Hash() uint64 {
	return hashStonecutter | uint64(s.Facing)<<8
}

func (c SugarCane) Hash() uint64 {
	return hashSugarCane | uint64(c.Age)<<8
}

func (TNT) Hash() uint64 {
	return hashTNT
}

func (g TallGrass) Hash() uint64 {
	return hashTallGrass | uint64(g.Type.Uint8())<<8
}

func (Terracotta) Hash() uint64 {
	return hashTerracotta
}

func (t Torch) Hash() uint64 {
	return hashTorch | uint64(t.Facing)<<8 | uint64(t.Type.Uint8())<<11
}

func (Tuff) Hash() uint64 {
	return hashTuff
}

func (w Wall) Hash() uint64 {
	return hashWall | w.Block.Hash()<<8 | uint64(w.NorthConnection.Uint8())<<24 | uint64(w.EastConnection.Uint8())<<26 | uint64(w.SouthConnection.Uint8())<<28 | uint64(w.WestConnection.Uint8())<<30 | uint64(boolByte(w.Post))<<32
}

func (w Water) Hash() uint64 {
	return hashWater | uint64(boolByte(w.Still))<<8 | uint64(w.Depth)<<9 | uint64(boolByte(w.Falling))<<17
}

func (Waterlily) Hash() uint64 {
	return hashWaterlily
}

func (s WheatSeeds) Hash() uint64 {
	return hashWheatSeeds | uint64(s.Growth)<<8
}

func (w Wood) Hash() uint64 {
	return hashWood | uint64(w.Wood.Uint8())<<8 | uint64(boolByte(w.Stripped))<<12 | uint64(w.Axis)<<13
}

func (d WoodDoor) Hash() uint64 {
	return hashWoodDoor | uint64(d.Wood.Uint8())<<8 | uint64(d.Facing)<<12 | uint64(boolByte(d.Open))<<14 | uint64(boolByte(d.Top))<<15 | uint64(boolByte(d.Right))<<16
}

func (w WoodFence) Hash() uint64 {
	return hashWoodFence | uint64(w.Wood.Uint8())<<8
}

func (f WoodFenceGate) Hash() uint64 {
	return hashWoodFenceGate | uint64(f.Wood.Uint8())<<8 | uint64(f.Facing)<<12 | uint64(boolByte(f.Open))<<14 | uint64(boolByte(f.Lowered))<<15
}

func (t WoodTrapdoor) Hash() uint64 {
	return hashWoodTrapdoor | uint64(t.Wood.Uint8())<<8 | uint64(t.Facing)<<12 | uint64(boolByte(t.Open))<<14 | uint64(boolByte(t.Top))<<15
}

func (w Wool) Hash() uint64 {
	return hashWool | uint64(w.Colour.Uint8())<<8
}
