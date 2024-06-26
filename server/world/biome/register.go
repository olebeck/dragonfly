package biome

import "github.com/df-mc/dragonfly/server/world"

// init registers all biomes that can be used in a world.World.
func init() {
	world.DefaultBiomes.Register(BadlandsPlateau{})
	world.DefaultBiomes.Register(Badlands{})
	world.DefaultBiomes.Register(BambooJungleHills{})
	world.DefaultBiomes.Register(BambooJungle{})
	world.DefaultBiomes.Register(BasaltDeltas{})
	world.DefaultBiomes.Register(Beach{})
	world.DefaultBiomes.Register(BirchForestHills{})
	world.DefaultBiomes.Register(BirchForest{})
	world.DefaultBiomes.Register(CherryGrove{})
	world.DefaultBiomes.Register(ColdOcean{})
	world.DefaultBiomes.Register(CrimsonForest{})
	world.DefaultBiomes.Register(DarkForestHills{})
	world.DefaultBiomes.Register(DarkForest{})
	world.DefaultBiomes.Register(DeepColdOcean{})
	world.DefaultBiomes.Register(DeepDark{})
	world.DefaultBiomes.Register(DeepFrozenOcean{})
	world.DefaultBiomes.Register(DeepLukewarmOcean{})
	world.DefaultBiomes.Register(DeepOcean{})
	world.DefaultBiomes.Register(DeepWarmOcean{})
	world.DefaultBiomes.Register(DesertHills{})
	world.DefaultBiomes.Register(DesertLakes{})
	world.DefaultBiomes.Register(Desert{})
	world.DefaultBiomes.Register(DripstoneCaves{})
	world.DefaultBiomes.Register(End{})
	world.DefaultBiomes.Register(ErodedBadlands{})
	world.DefaultBiomes.Register(FlowerForest{})
	world.DefaultBiomes.Register(Forest{})
	world.DefaultBiomes.Register(FrozenOcean{})
	world.DefaultBiomes.Register(FrozenPeaks{})
	world.DefaultBiomes.Register(FrozenRiver{})
	world.DefaultBiomes.Register(GiantSpruceTaigaHills{})
	world.DefaultBiomes.Register(GiantTreeTaigaHills{})
	world.DefaultBiomes.Register(GravellyMountainsPlus{})
	world.DefaultBiomes.Register(Grove{})
	world.DefaultBiomes.Register(IceSpikes{})
	world.DefaultBiomes.Register(JaggedPeaks{})
	world.DefaultBiomes.Register(JungleEdge{})
	world.DefaultBiomes.Register(JungleHills{})
	world.DefaultBiomes.Register(Jungle{})
	world.DefaultBiomes.Register(LegacyFrozenOcean{})
	world.DefaultBiomes.Register(LukewarmOcean{})
	world.DefaultBiomes.Register(LushCaves{})
	world.DefaultBiomes.Register(MangroveSwamp{})
	world.DefaultBiomes.Register(Meadow{})
	world.DefaultBiomes.Register(ModifiedBadlandsPlateau{})
	world.DefaultBiomes.Register(ModifiedJungleEdge{})
	world.DefaultBiomes.Register(ModifiedJungle{})
	world.DefaultBiomes.Register(ModifiedWoodedBadlandsPlateau{})
	world.DefaultBiomes.Register(MountainEdge{})
	world.DefaultBiomes.Register(MushroomFieldShore{})
	world.DefaultBiomes.Register(MushroomFields{})
	world.DefaultBiomes.Register(NetherWastes{})
	world.DefaultBiomes.Register(Ocean{})
	world.DefaultBiomes.Register(OldGrowthBirchForest{})
	world.DefaultBiomes.Register(OldGrowthPineTaiga{})
	world.DefaultBiomes.Register(OldGrowthSpruceTaiga{})
	world.DefaultBiomes.Register(Plains{})
	world.DefaultBiomes.Register(River{})
	world.DefaultBiomes.Register(SavannaPlateau{})
	world.DefaultBiomes.Register(Savanna{})
	world.DefaultBiomes.Register(ShatteredSavannaPlateau{})
	world.DefaultBiomes.Register(SnowyBeach{})
	world.DefaultBiomes.Register(SnowyMountains{})
	world.DefaultBiomes.Register(SnowyPlains{})
	world.DefaultBiomes.Register(SnowySlopes{})
	world.DefaultBiomes.Register(SnowyTaigaHills{})
	world.DefaultBiomes.Register(SnowyTaigaMountains{})
	world.DefaultBiomes.Register(SnowyTaiga{})
	world.DefaultBiomes.Register(SoulSandValley{})
	world.DefaultBiomes.Register(StonyPeaks{})
	world.DefaultBiomes.Register(StonyShore{})
	world.DefaultBiomes.Register(SunflowerPlains{})
	world.DefaultBiomes.Register(SwampHills{})
	world.DefaultBiomes.Register(Swamp{})
	world.DefaultBiomes.Register(TaigaHills{})
	world.DefaultBiomes.Register(TaigaMountains{})
	world.DefaultBiomes.Register(Taiga{})
	world.DefaultBiomes.Register(TallBirchHills{})
	world.DefaultBiomes.Register(WarmOcean{})
	world.DefaultBiomes.Register(WarpedForest{})
	world.DefaultBiomes.Register(WindsweptForest{})
	world.DefaultBiomes.Register(WindsweptGravellyHills{})
	world.DefaultBiomes.Register(WindsweptHills{})
	world.DefaultBiomes.Register(WindsweptSavanna{})
	world.DefaultBiomes.Register(WoodedBadlandsPlateau{})
	world.DefaultBiomes.Register(WoodedHills{})
}
