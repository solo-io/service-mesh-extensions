package test

import (
	"fmt"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
)

func GetDefaultParameters(versionedSpec *v1.VersionedApplicationSpec, flavorName string, layerInputs []render.LayerInput) map[string]string {
	var flavor *v1.Flavor
	for _, f := range versionedSpec.Flavors {
		if f.Name == flavorName {
			flavor = f
			break
		}
	}
	if flavor == nil {
		panic(fmt.Sprintf("could not find flavor %s in spec with version %s", flavorName, versionedSpec.Version))
	}

	result := make(map[string]string)
	for _, param := range versionedSpec.Parameters {
		// TODO joekelley extract function to handle params of all types.
		result[param.Name] = param.Default.GetString_()
	}
	for _, param := range flavor.Parameters {
		result[param.Name] = param.Default.GetString_()
	}
	for _, layer := range flavor.CustomizationLayers {
		for _, input := range layerInputs {
			if layer.Id == input.LayerId {
				for _, option := range layer.Options {
					if option.Id == input.OptionId {
						for _, param := range option.Parameters {
							result[param.Name] = param.Default.GetString_()
						}
					}
				}
			}
		}
	}
	return result
}
