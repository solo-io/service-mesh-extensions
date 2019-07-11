package registry_test

import (
	"testing"

	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const DefaultRemoteUrl = "https://storage.googleapis.com/sm-marketplace-registry/extensions4.yaml"

func getAppSpec(name string, appType hubv1.ApplicationType, versions ...*v1.VersionedApplicationSpec) *v1.ApplicationSpec {
	return &v1.ApplicationSpec{
		Name:     name,
		Type:     appType,
		Versions: versions,
	}
}

func getVersion(version string, flavors ...*v1.Flavor) *v1.VersionedApplicationSpec {
	return &v1.VersionedApplicationSpec{
		Version: version,
		Flavors: flavors,
	}
}

func TestRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Suite")
}
