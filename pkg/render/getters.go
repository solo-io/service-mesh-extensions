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

	InvalidOptionIdError = func(optionId, layerId string) error {
		return errors.Errorf("Layer option %v not found for layer %v", optionId, layerId)
	}

	InvalidLayerIdError = func(layerId, flavorName string) error {
		return errors.Errorf("Layer %v not found for flavor %v", layerId, flavorName)
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

func GetLayerOption(optionId string, layer *v1.Layer) (*v1.LayerOption, error) {
	for _, option := range layer.Options {
		if optionId == option.Id {
			return option, nil
		}
	}
	return nil, InvalidOptionIdError(optionId, layer.Id)
}

func GetLayer(layerId string, flavor *v1.Flavor) (*v1.Layer, error) {
	for _, layer := range flavor.CustomizationLayers {
		if layerId == layer.Id {
			return layer, nil
		}
	}
	return nil, InvalidLayerIdError(layerId, flavor.Name)
}

func GetLayerOptionFromFlavor(layerId, optionId string, flavor *v1.Flavor) (*v1.LayerOption, error) {
	layer, err := GetLayer(layerId, flavor)
	if err != nil {
		return nil, err
	}
	option, err := GetLayerOption(optionId, layer)
	if err != nil {
		return nil, err
	}
	return option, nil
}
