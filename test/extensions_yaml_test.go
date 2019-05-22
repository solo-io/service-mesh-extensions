package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		for _, ext := range extensions {
			extension := ext
			It(fmt.Sprintf("extensions/v1/%s/spec.yaml is valid", extension.Name()), func() {
				Expect(extension.IsDir()).To(BeTrue())
				specPath := filepath.Join(rootDir, extension.Name(), specYamlFilename)
				spec := LoadApplicationSpec(specPath)
				Expect(spec.Name).To(BeEquivalentTo(extension.Name()))
			})
		}
	})
})
