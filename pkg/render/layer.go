package render

import (
	"context"

	renderinputs "github.com/solo-io/service-mesh-hub/pkg/render/inputs"
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
	for _, layerInput := range inputs.Layers {
		option, err := GetLayerOptionFromFlavor(layerInput.LayerId, layerInput.OptionId, inputs.Flavor)
		if err != nil {
			return nil, err
		}

		if option.Kustomize != nil {
			layerEngine, err := kustomize.NewKustomizer(kustomizeLoader, manifests, option.Kustomize,
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
	manifestBytes := []byte(manifests.CombinedString())
	resources, err := YamlToResources(manifestBytes)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func getRenderValues(inputs ValuesInputs) (interface{}, error) {
	customValues, err := ConvertParamsToNestedMap(inputs.Params)
	if err != nil {
		return nil, err
	}

	return renderinputs.ManifestRenderValues{
		Name:             inputs.Name,
		InstallNamespace: inputs.InstallNamespace,
		MeshRef:          inputs.MeshRef,
		Custom:           customValues,
	}, nil
}
