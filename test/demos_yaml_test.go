package test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Demos Yaml Test", func() {
	const (
		specYamlFilename = "spec.yaml"
		rootDir          = "../demos/v1"
	)

	demos, err := ioutil.ReadDir(rootDir)
	if err != nil {
		Fail(err.Error())
	}

	Context("spec yaml validity", func() {
		for _, d := range demos {
			demo := d
			It(fmt.Sprintf("demos/v1/%s/spec.yaml is valid", demo.Name()), func() {
				Expect(demo.IsDir()).To(BeTrue())
				specPath := filepath.Join(rootDir, demo.Name(), specYamlFilename)
				spec := LoadApplicationSpec(specPath)
				Expect(spec.Name).To(BeEquivalentTo(demo.Name()))
				Expect(spec.Type).To(BeEquivalentTo(v1.ApplicationType_DEMO))
			})
		}
	})
})
