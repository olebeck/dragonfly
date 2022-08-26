package features

import "github.com/df-mc/dragonfly/server/world"

func init() {
	world.RegisterFeature(&OakTree{})
	world.RegisterFeature(&SpruceTree{})
	world.RegisterFeature(&BirchTree{})
	world.RegisterFeature(&AzaleaTree{})
}
