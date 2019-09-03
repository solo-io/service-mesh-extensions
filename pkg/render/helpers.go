package render

import (
	"github.com/solo-io/go-utils/errors"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

var (
	InvalidOptionIdError = func(optionId, layerId string) error {
		return errors.Errorf("Layer option %v not found for layer %v", optionId, layerId)
	}

	InvalidLayerIdError = func(layerId, flavorName string) error {
		return errors.Errorf("Layer %v not found for flavor %v", layerId, flavorName)
	}
)

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

func GetLayerOptionTwo(layerId, optionId string, flavor *v1.Flavor) (*v1.LayerOption, error) {
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
