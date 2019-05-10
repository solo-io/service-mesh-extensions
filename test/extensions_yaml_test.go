package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
)



var _ = Describe("Extensions Yaml Test", func() {
	const (
		specYamlFilename = "spec.yaml"
		rootDir          = "../extensions/v1"
	)

	extensions, err := ioutil.ReadDir(rootDir)
	if err != nil {
		Fail(err.Error())
	}

	Context("spec yaml validity", func() {
		for _, extension := range extensions {
			It(fmt.Sprintf("extensions/v1/%s/spec.yaml is valid", extension.Name()), func() {
				Expect(extension.IsDir()).To(BeTrue())
				specPath := filepath.Join(rootDir, extension.Name(), specYamlFilename)
				bytes, err := ioutil.ReadFile(specPath)
				Expect(err).NotTo(HaveOccurred())
				var spec v1.ApplicationSpec
				err = protoutils.UnmarshalYaml(bytes, &spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Name).To(BeEquivalentTo(extension.Name()))

			})
		}
	})
})
