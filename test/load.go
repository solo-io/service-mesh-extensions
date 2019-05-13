package test

import (
	. "github.com/onsi/gomega"

	"io/ioutil"

	"github.com/solo-io/go-utils/protoutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

func LoadExtensionSpec(pathToSpec string) *v1.ApplicationSpec {
	bytes, err := ioutil.ReadFile(pathToSpec)
	Expect(err).NotTo(HaveOccurred())
	var spec v1.ApplicationSpec
	err = protoutils.UnmarshalYaml(bytes, &spec)
	Expect(err).NotTo(HaveOccurred())
	return &spec
}
