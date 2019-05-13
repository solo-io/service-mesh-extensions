package test

import (
	. "github.com/onsi/gomega"

	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"io/ioutil"
)

func LoadExtensionSpec(pathToSpec string) *v1.ApplicationSpec {
	bytes, err := ioutil.ReadFile(pathToSpec)
	Expect(err).NotTo(HaveOccurred())
	var spec v1.ApplicationSpec
	err = protoutils.UnmarshalYaml(bytes, &spec)
	Expect(err).NotTo(HaveOccurred())
	return &spec
}