package test

import (
	"fmt"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/pkg/render/util"
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
		v, err := util.ParamValueToString(param.Default, util.PlainTextSecretGetter)
		if err != nil {
			panic(err.Error())
		}
		result[param.Name] = v
	}
	for _, param := range flavor.Parameters {
		v, err := util.ParamValueToString(param.Default, util.PlainTextSecretGetter)
		if err != nil {
			panic(err.Error())
		}
		result[param.Name] = v
	}
	for _, layer := range flavor.CustomizationLayers {
		for _, input := range layerInputs {
			if layer.Id == input.LayerId {
				for _, option := range layer.Options {
					if option.Id == input.OptionId {
						for _, param := range option.Parameters {
							v, err := util.ParamValueToString(param.Default, util.PlainTextSecretGetter)
							if err != nil {
								panic(err.Error())
							}
							result[param.Name] = v
						}
					}
				}
			}
		}
	}
	return result
}
