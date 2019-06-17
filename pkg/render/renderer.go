package render

import (
	"context"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

type ManifestRenderer interface {
	// Given the spec and flavor definition, generate a set of kube resources that represent the exact install manifest.
	// Note: to allow for custom flavor definitions, flavor is passed in here, and it is not required that the spec contains
	// it as one of the built-in flavors, as long as it is valid for the user inputs.
	ComputeResourcesForApplication(ctx context.Context, inputs ValuesInputs, spec *v1.VersionedApplicationSpec, flavor *v1.Flavor) (kuberesource.UnstructuredResources, error)
}

type manifestRenderer struct {
}

func NewManifestRenderer() ManifestRenderer {
	return &manifestRenderer{}
}

func (m *manifestRenderer) ComputeResourcesForApplication(ctx context.Context, inputs ValuesInputs, spec *v1.VersionedApplicationSpec, flavor *v1.Flavor) (kuberesource.UnstructuredResources, error) {
	inputs, err := ExecInputValuesTemplates(inputs)
	if err != nil {
		return nil, FailedRenderValueTemplatesError(err)
	}

	manifests, err := GetManifestsFromApplicationSpec(ctx, inputs, spec)
	if err != nil {
		return nil, err
	}

	if err := ValidateInputsAgainstFlavor(inputs, flavor); err != nil {
		return nil, err
	}

	rawResources, err := ApplyLayers(ctx, inputs, manifests)
	if err != nil {
		return nil, err
	}
	return FilterByLabel(ctx, spec, rawResources), nil
}
