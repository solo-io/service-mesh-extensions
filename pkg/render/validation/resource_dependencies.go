package validation

import v1 "github.com/solo-io/service-mesh-hub/api/v1"

type ValidateEnvironment func(deps []*v1.EnvironmentRequirements) error

func NoopValidateEnvironment(deps []*v1.EnvironmentRequirements) error {
	return nil
}
