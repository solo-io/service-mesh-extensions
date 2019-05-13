package render

import (
	"context"

	"github.com/solo-io/service-mesh-hub/pkg/kustomize"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"

	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/spf13/afero"
)

const (
	layerDirPrefix = "layers"
)

var (
	UnknownLayerTypeError = func(layer *hubv1.Layer) error {
		return errors.Errorf("unknown layer type specified %T", layer.GetType())
	}
)

func ApplyLayers(ctx context.Context, installedFlavor *hubv1.Flavor, manifests helmchart.Manifests) (kuberesource.UnstructuredResources, error) {

	if installedFlavor.CustomizationLayers == nil || len(installedFlavor.CustomizationLayers) == 0 {
		return GetResourcesFromManifests(ctx, manifests)
	} else if len(installedFlavor.CustomizationLayers) >= 2 {
		return nil, ExpectedAtMostError("customization", 1, len(installedFlavor.CustomizationLayers))
	}

	fs := afero.NewOsFs()
	execDir, err := afero.TempDir(fs, "", layerDirPrefix)
	if err != nil {
		return nil, err
	}

	layer := installedFlavor.CustomizationLayers[0]
	var layerEngine kustomize.LayerEngine
	switch layerType := layer.GetType().(type) {
	case *hubv1.Layer_Kustomize:
		kustomizeLoader := loader.NewKustomizeLoader(ctx, fs)
		layerEngine = kustomize.NewKustomizer(kustomizeLoader, manifests, layerType.Kustomize)
	default:
		return nil, UnknownLayerTypeError(layer)
	}

	manifestBytes, err := layerEngine.Run(execDir)
	if err != nil {
		return nil, err
	}

	resources, err := YamlToResources(manifestBytes)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
