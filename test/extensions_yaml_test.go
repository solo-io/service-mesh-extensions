package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/go-utils/protoutils"
	"io/ioutil"
	"path/filepath"
)

var _ = Describe("Extensions Yaml Test", func() {

	const (
		specYamlFilename = "spec.yaml"
	)

	Context("extensions/v1/EXTENSION_NAME/spec.yaml", func() {
		It("is valid yaml", func() {
			rootDir := "../extensions/v1"
			extensions, err := ioutil.ReadDir(rootDir)
			Expect(err).NotTo(HaveOccurred())

			for _, extension := range extensions {
				Expect(extension.IsDir()).To(BeTrue())
				specPath := filepath.Join(rootDir, extension.Name(), specYamlFilename)
				bytes, err := ioutil.ReadFile(specPath)
				Expect(err).NotTo(HaveOccurred())
				var spec v1.ApplicationSpec
				err = protoutils.UnmarshalYaml(bytes, &spec)
				Expect(err).NotTo(HaveOccurred())
				Expect(spec.Name).To(BeEquivalentTo(extension.Name()))
			}
		})
	})
})