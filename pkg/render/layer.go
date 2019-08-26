package render

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/helm/pkg/releaseutil"

	"github.com/solo-io/service-mesh-hub/pkg/kustomize/plugins"

	"github.com/solo-io/service-mesh-hub/pkg/kustomize"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"

	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	"github.com/spf13/afero"
)

const (
	layerDirPrefix = "layers"
)

var (
	FailedToCalculateValues = func(err error) error {
		return errors.Wrapf(err, "failed to calculate values for layer rendering")
	}
)

func ApplyLayers(ctx context.Context, inputs ValuesInputs, manifests helmchart.Manifests) (kuberesource.UnstructuredResources, error) {

	fs := afero.NewOsFs()
	execDir, err := afero.TempDir(fs, "", layerDirPrefix)
	if err != nil {
		return nil, err
	}

	values, err := getRenderValues(inputs)
	if err != nil {
		return nil, FailedToCalculateValues(err)
	}

	kustomizeLoader := loader.NewKustomizeLoader(ctx, fs)
	var manifestBytes []byte
	for _, layerInput := range inputs.Layers {
		if layerInput.Option != nil && layerInput.Option.Kustomize != nil {
			layerEngine, err := kustomize.NewKustomizer(kustomizeLoader, manifests, layerInput.Option.Kustomize,
				plugins.NewManifestRenderPlugin(values))
			if err != nil {
				return nil, err
			}
			manifestBytes, err := layerEngine.Run(execDir)
			if err != nil {
				return nil, err
			}
			manifests = helmchart.Manifests{{Head: &releaseutil.SimpleHead{}, Content: string(manifestBytes)}}

		}
	}
	if manifestBytes == nil {
		manifestBytes = []byte(manifests.CombinedString())
	}
	resources, err := YamlToResources(manifestBytes)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func getRenderValues(inputs ValuesInputs) (interface{}, error) {

	// TODO: get rid of this and just use the ValuesInputs type
	type manifestRenderValues struct {
		Name             string
		InstallNamespace string
		FlavorName       string
		MeshRef          core.ResourceRef

		Supergloo SuperglooInfo

		// Custom values come from the parameters set on a flavor
		Custom interface{}
	}

	customValues, err := ConvertParamsToNestedMap(inputs.FlavorParams)
	if err != nil {
		return nil, err
	}

	return manifestRenderValues{
		Name:             inputs.Name,
		InstallNamespace: inputs.InstallNamespace,
		FlavorName:       inputs.FlavorName,
		MeshRef:          inputs.MeshRef,
		Supergloo:        inputs.Supergloo,
		Custom:           customValues,
	}, nil
}
