package features

import "github.com/df-mc/dragonfly/server/world"

func init() {
	world.RegisterFeature(&OakTree{})
	world.RegisterFeature(&LargeDarkOakTree{})
	world.RegisterFeature(&SpruceTree{})
	// world.RegisterFeature(&LargeSpruceTree{}) needs leaves done
	world.RegisterFeature(&BirchTree{})
	world.RegisterFeature(&JungleTree{})
	world.RegisterFeature(&LargeJungleTree{})
	world.RegisterFeature(&AcaciaTree{})
	world.RegisterFeature(&AzaleaTree{})
	world.RegisterFeature(&HugeBrownMushroom{})
}
