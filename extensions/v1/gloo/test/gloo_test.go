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

var _ = Describe("gloo extension test", func() {

	const (
		superglooNamesapce = "sm-marketplace"
		namespace          = "gloo-system"
		name               = "gloo"
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
			testInput    = func(flavorName string, layers []render.LayerInput) render.ValuesInputs {
				return render.ValuesInputs{
					Name:             name,
					Flavor:           test.GetFlavor(flavorName, version),
					Layers:           layers,
					InstallNamespace: namespace,
					MeshRef: core.ResourceRef{
						Name:      meshName,
						Namespace: namespace,
					},
					SpecDefinedValues: version.ValuesYaml,
				}
			}
		)

		PContext("0.13.26 with supergloo overlay", func() {
			BeforeEach(func() {
				version = versionMap["0.13.26"]
				layers := []render.LayerInput{{
					LayerId:  "supergloo",
					OptionId: "ingress",
				}}
				inputs = testInput("supergloo", layers)
				rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
				Expect(err).NotTo(HaveOccurred())
				testManifest = NewTestManifestWithResources(rendered)
			})

			It("has the correct number of resources", func() {
				Expect(testManifest.NumResources()).To(Equal(14))
			})

			It("has a mesh ingress", func() {
				testManifest.ExpectCustomResource("MeshIngress", superglooNamesapce, name)
			})
		})
		Context("0.13.26 with vanilla overlay", func() {
			BeforeEach(func() {
				version = versionMap["0.13.26"]
				layers := []render.LayerInput{{
					LayerId:  "custom-resources",
					OptionId: "create",
				}}
				inputs = testInput("vanilla", layers)
				rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
				Expect(err).NotTo(HaveOccurred())
				testManifest = NewTestManifestWithResources(rendered)
			})

			It("has the correct number of resources", func() {
				Expect(testManifest.NumResources()).To(Equal(13))
			})
		})
	})

})
