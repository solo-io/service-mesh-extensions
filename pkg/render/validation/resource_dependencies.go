package validation

import (
	"errors"
	"fmt"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

type ValidateResourceDependencies func(deps []*v1.ResourceDependency) error

func NoopValidateResources(deps []*v1.ResourceDependency) error {
	return nil
}

func ValidateResources(deps []*v1.ResourceDependency) error {
	for _, dep := range deps {
		switch dep.GetType().(type) {
		case *v1.ResourceDependency_SecretDependency: {
			if dep.GetSecretDependency().GetName() == "" {
				return errors.New("secret dependency has no name")
			} else if len(dep.GetSecretDependency().GetKeys()) == 0 {
				return errors.New("secret dependency has no keys")
			}
		}
		default: {
			return fmt.Errorf("unknown resource dependency type %v", dep.GetType())
		}
		}
	}
	return nil
}