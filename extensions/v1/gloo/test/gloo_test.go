package test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/go-utils/manifesttestutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("gloo extension test", func() {

	const (
		namespace = "gloo-system"
		name      = "gloo"
	)

	var (
		spec       *v1.ApplicationSpec
		versionMap map[string]*v1.VersionedApplicationSpec
		labels     = map[string]string{
			"app": name,
		}
	)


	BeforeEach(func() {
		spec = test.LoadExtensionSpec("../spec.yaml")
		versionMap = make(map[string]*v1.VersionedApplicationSpec)
		for _, version := range spec.Versions {
			versionMap[version.Version] = version
		}
	})


	Context("linkerd", func() {
		const (
			meshName  = "linkerd"
		)
	})

	Context("istio", func() {
		const (
			meshName  = "istio"
		)

		Context("0.13.26 with default values", func() {
			var (
				version      *v1.VersionedApplicationSpec
				inputs       render.ValuesInputs
				testManifest TestManifest
			)

			BeforeEach(func() {
				version = versionMap["0.13.26"]
				inputs = render.ValuesInputs{
					Name:             name,
					FlavorName:       meshName,
					InstallNamespace: namespace,
					MeshRef: core.ResourceRef{
						Name:      meshName,
						Namespace: namespace,
					},
					SpecDefinedValues: version.ValuesYaml,
				}
				rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
				Expect(err).NotTo(HaveOccurred())
				testManifest = NewTestManifestWithResources(rendered)
			})

			FIt("has the correct number of resources", func() {
				Expect(testManifest.NumResources()).To(Equal(7))
			})

			It("has a demo secret", func() {
				rb := ResourceBuilder{
					Name:      name,
					Namespace: namespace,
					Data: map[string]string{
						"username":   "admin",
						"passphrase": "admin",
					},
					Labels: labels,
				}
				testManifest.ExpectSecret(rb.GetSecret())
			})
		})
	})




})
