package test

import (
	. "github.com/onsi/gomega"

	"io/ioutil"

	"github.com/solo-io/go-utils/protoutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

func LoadApplicationSpec(pathToSpec string) *v1.ApplicationSpec {
	bytes, err := ioutil.ReadFile(pathToSpec)
	Expect(err).NotTo(HaveOccurred())
	var spec v1.ApplicationSpec
	err = protoutils.UnmarshalYaml(bytes, &spec)
	Expect(err).NotTo(HaveOccurred())
	return &spec
}

func GetFlavor(name string, spec *v1.VersionedApplicationSpec) *v1.Flavor {
	for _, flavor := range spec.Flavors {
		if flavor.Name == name {
			return flavor
		}
	}
	return nil
}
