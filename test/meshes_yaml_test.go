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
		rootDir          = "../meshes/v1"
	)

	meshes, err := ioutil.ReadDir(rootDir)
	if err != nil {
		Fail(err.Error())
	}

	Context("spec yaml validity", func() {
		for _, m := range meshes {
			mesh := m
			It(fmt.Sprintf("meshes/v1/%s/spec.yaml is valid", mesh.Name()), func() {
				Expect(mesh.IsDir()).To(BeTrue())
				specPath := filepath.Join(rootDir, mesh.Name(), specYamlFilename)
				spec := LoadApplicationSpec(specPath)
				Expect(spec.Name).To(BeEquivalentTo(mesh.Name()))
				Expect(spec.Type).To(BeEquivalentTo(v1.ApplicationType_MESH))
			})
		}
	})
})
