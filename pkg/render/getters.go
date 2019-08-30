package render

import (
	"github.com/solo-io/go-utils/errors"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

var (
	NilFlavorError = errors.Errorf("flavor name cannot be nil")

	ExpectedAtMostError = func(item string, desired, found int) error {
		return errors.Errorf("expected at most %d %s, found %d", desired, item, found)
	}

	NoFlavorFoundError = func(name string) error {
		return errors.Errorf("could not find flavor with name: %s", name)
	}

	UnexpectedFlavorError = func(expected, actual string) error {
		return errors.Errorf("user selected flavor %s, but renderer was provided flavor %s", expected, actual)
	}
)

func GetInstalledFlavor(name string, flavors []*v1.Flavor) (*v1.Flavor, error) {
	if name == "" {
		return nil, NilFlavorError
	} else if len(flavors) == 0 {
		return nil, ExpectedAtMostError("flavor", 1, len(flavors))
	}

	for _, flavor := range flavors {
		if flavor.Name == name {
			return flavor, nil
		}
	}
	return nil, NoFlavorFoundError(name)
}

func GetRequiredLayerCount(flavor *v1.Flavor) int {
	count := 0
	for _, layer := range flavor.GetCustomizationLayers() {
		if !layer.Optional {
			count++
		}
	}
	return count
}
