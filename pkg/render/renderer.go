package render

import (
	"context"

	"github.com/solo-io/go-utils/installutils/kuberesource"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render/validation"
)

//go:generate mockgen -source=./renderer.go -package mocks -destination=./mocks/mock_render.go ManifestRenderer

type ManifestRenderer interface {
	// Given the spec and values inputs, generate a set of kube resources that represent the exact install manifest.
	ComputeResourcesForApplication(ctx context.Context, inputs ValuesInputs, spec *v1.VersionedApplicationSpec) (kuberesource.UnstructuredResources, error)
}

type manifestRenderer struct {
	validateEnvironment validation.ValidateResourceDependencies
}

func NewManifestRenderer(validateFn validation.ValidateResourceDependencies) ManifestRenderer {
	return &manifestRenderer{validateEnvironment: validateFn}
}

func (m *manifestRenderer) ComputeResourcesForApplication(ctx context.Context, inputs ValuesInputs, spec *v1.VersionedApplicationSpec) (kuberesource.UnstructuredResources, error) {
	if err := ValidateInputs(inputs, *spec, m.validateEnvironment); err != nil {
		return nil, err
	}

	inputs, err := ExecInputValuesTemplates(inputs)
	if err != nil {
		return nil, FailedRenderValueTemplatesError(err)
	}

	manifests, err := GetManifestsFromApplicationSpec(ctx, inputs, spec)
	if err != nil {
		return nil, err
	}

	rawResources, err := GetResources(manifests)
	if err != nil {
		return nil, err
	}
	return FilterByLabel(ctx, spec, rawResources), nil
}
