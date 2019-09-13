package test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/manifesttestutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("smi extension test", func() {

	const (
		namespace = "smi"
		name      = "smi"
	)

	var (
		spec       *v1.ApplicationSpec
		versionMap map[string]*v1.VersionedApplicationSpec
	)

	BeforeEach(func() {
		spec = test.LoadApplicationSpec("../spec.yaml")
		versionMap = make(map[string]*v1.VersionedApplicationSpec)
		for _, version := range spec.Versions {
			versionMap[version.Version] = version
		}
	})

	Context("istio", func() {
		const (
			meshName = "istio"
		)

		var (
			version      *v1.VersionedApplicationSpec
			inputs       render.ValuesInputs
			testManifest TestManifest
			testInput    = func(flavorName string) render.ValuesInputs {
				return render.ValuesInputs{
					Name:             name,
					Flavor:           test.GetFlavor(flavorName, version),
					InstallNamespace: namespace,
					MeshRef: core.ResourceRef{
						Name:      meshName,
						Namespace: namespace,
					},
					SpecDefinedValues: version.ValuesYaml,
				}
			}
		)

		Context("0.0.1", func() {
			BeforeEach(func() {
				version = versionMap["0.0.1"]
				inputs = testInput("istio")
				rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
				Expect(err).NotTo(HaveOccurred())
				testManifest = NewTestManifestWithResources(rendered)
			})

			It("has the correct number of resources", func() {
				Expect(testManifest.NumResources()).To(Equal(5))
			})
		})
	})

})
