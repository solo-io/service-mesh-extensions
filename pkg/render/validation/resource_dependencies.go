package validation

import v1 "github.com/solo-io/service-mesh-hub/api/v1"

type ValidateResourceDependencies func(deps []*v1.ResourceDependency) error

func NoopValidateResources(deps []*v1.ResourceDependency) error {
	return nil
}
