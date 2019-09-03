package test

import (
	"fmt"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

func GetDefaultParameters(versionedSpec *v1.VersionedApplicationSpec, flavorName string) map[string]string {
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
	for _, param := range flavor.Parameters {
		result[param.Name] = param.Default.GetString_()
	}
	return result
}
