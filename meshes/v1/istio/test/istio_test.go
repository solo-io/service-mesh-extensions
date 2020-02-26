package test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/manifesttestutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
)

var _ = XDescribe("istio extension test", func() {

	const (
		namespace = "istio"
		name      = "istio"
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
		var (
			version      *v1.VersionedApplicationSpec
			inputs       render.ValuesInputs
			testManifest TestManifest
			testInput    = func(flavorName string) render.ValuesInputs {
				return render.ValuesInputs{
					Name:              name,
					Flavor:            test.GetFlavor(flavorName, version),
					InstallNamespace:  namespace,
					SpecDefinedValues: version.ValuesYaml,
				}
			}
		)

		Context("1.1.7", func() {
			Context("vanilla", func() {
				BeforeEach(func() {
					version = versionMap["1.1.7"]
					inputs = testInput("vanilla")
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
				})

				It("has the correct number of resources", func() {
					Expect(testManifest.NumResources()).To(Equal(140))
				})
			})

			Context("banana", func() {
				BeforeEach(func() {
					version = versionMap["1.1.7"]
					inputs = testInput("banana")
					inputs.Layers = []render.LayerInput{
						{LayerId: "security", OptionId: "strict-custom-cert"},
						{LayerId: "gateway", OptionId: "disabled"},
					}
					inputs.Params = map[string]string{"cert.not.implemented": "barbaz"}
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
				})

				It("has the correct number of resources", func() {
					Expect(testManifest.NumResources()).To(Equal(140))
				})
			})
		})

		Context("1.3", func() {
			Context("vanilla", func() {
				BeforeEach(func() {
					version = versionMap["1.3"]
					inputs = testInput("vanilla")
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
				})

				It("has the correct number of resources", func() {
					Expect(testManifest.NumResources()).To(Equal(136))
				})
			})

			Context("banana", func() {
				BeforeEach(func() {
					version = versionMap["1.3"]
					inputs = testInput("banana")
					inputs.Layers = []render.LayerInput{
						{LayerId: "security", OptionId: "strict-custom-cert"},
						{LayerId: "gateway", OptionId: "disabled"},
					}
					inputs.Params = map[string]string{"cert.not.implemented": "barbaz"}
					rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
					Expect(err).NotTo(HaveOccurred())
					testManifest = NewTestManifestWithResources(rendered)
				})

				It("has the correct number of resources", func() {
					Expect(testManifest.NumResources()).To(Equal(136))
				})
			})
		})
	})

})
