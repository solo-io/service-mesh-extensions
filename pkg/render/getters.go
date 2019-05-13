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
