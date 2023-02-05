package world

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type description struct {
	Identifier             string         `json:"identifier"`
	IsExperimental         bool           `json:"is_experimental"`
	RegisterToCreativeMenu bool           `json:"register_to_creative_menu"`
	Properties             map[string]any `json:"properties,omitempty"`
}

type menu_category struct {
	Category string `json:"category"`
	Group    string `json:"group"`
}

func menu_category_from_map(in map[string]any) menu_category {
	return menu_category{
		Category: in["category"].(string),
		Group:    in["group"].(string),
	}
}

type permutation struct {
	Components map[string]any `json:"components"`
	Condition  string         `json:"condition"`
}

func permutation_from_map(in map[string]any) permutation {
	return permutation{
		Components: in["components"].(map[string]any),
		Condition:  in["condition"].(string),
	}
}

type MinecraftBlock struct {
	Description  description    `json:"description"`
	Components   map[string]any `json:"components,omitempty"`
	MenuCategory menu_category  `json:"menu_category,omitempty"`
	Permutations []permutation  `json:"permutations,omitempty"`
}

func ParseBlock(block protocol.BlockEntry) MinecraftBlock {
	entry := MinecraftBlock{
		Description: description{
			Identifier:             block.Name,
			IsExperimental:         true,
			RegisterToCreativeMenu: true,
		},
	}

	if perms, ok := block.Properties["permutations"].([]any); ok {
		for _, v := range perms {
			entry.Permutations = append(entry.Permutations, permutation_from_map(v.(map[string]any)))
		}
	}

	if comps, ok := block.Properties["components"].(map[string]any); ok {
		delete(comps, "minecraft:creative_category")

		for k, v := range comps {
			if v, ok := v.(map[string]any); ok {
				// fix {"value": 0.1} -> 0.1
				if v, ok := v["value"]; ok {
					comps[k] = v
				}
				// fix {"lightLevel": 15} -> 15
				if v, ok := v["lightLevel"]; ok {
					comps[k] = v
				}
				// fix {"triggerType": "name"} -> "name"
				if v, ok := v["triggerType"]; ok {
					comps[k] = v
				}
				// fix missing * instance
				if k == "minecraft:material_instances" {
					if m, ok := v["materials"].(map[string]any); ok {
						comps[k] = m
					}
				}
			}
		}
		entry.Components = comps
	}

	if menu, ok := block.Properties["menu_category"].(map[string]any); ok {
		entry.MenuCategory = menu_category_from_map(menu)
	}
	if props, ok := block.Properties["properties"].([]any); ok {
		entry.Description.Properties = make(map[string]any)
		for _, v := range props {
			v := v.(map[string]any)
			name := v["name"].(string)
			switch a := v["enum"].(type) {
			case []int32:
				entry.Description.Properties[name] = a
			case []bool:
				entry.Description.Properties[name] = a
			}

		}
	}
	return entry
}
