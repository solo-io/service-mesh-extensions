package test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
)

var _ = Describe("kiali", func() {
	var (
		spec *v1.ApplicationSpec
		versionMap map[string]*v1.VersionedApplicationSpec
	)

	BeforeEach(func() {
		spec = test.LoadExtensionSpec("../spec.yaml")
		versionMap = make(map[string]*v1.VersionedApplicationSpec)
		for _, version := range spec.Versions {
			versionMap[version.Version] = version
		}
	})

	Context("0.12", func() {
		var (
			version *v1.VersionedApplicationSpec
			inputs render.ValuesInputs
			resources kuberesource.UnstructuredResources
		)

		BeforeEach(func() {
			version = versionMap["0.12"]
			inputs = render.ValuesInputs{
				Name: "kiali",
				FlavorName: "istio",
				InstallNamespace: "istio-system",
				MeshRef: v1.ResourceRef{
					Name: "istio",
					Namespace: "istio-system",
				},
				SpecDefinedValues: version.ValuesYaml,
			}
			rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
			Expect(err).NotTo(HaveOccurred())
			resources = rendered
		})

		It("has the correct number of resources", func() {
			Expect(len(resources)).To(Equal(7))
		})
	})
})
