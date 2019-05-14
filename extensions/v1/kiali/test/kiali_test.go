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

var _ = Describe("kiali", func() {

	const (
		namespace = "istio-system"
		name      = "kiali"
		meshName  = "istio"
	)

	var (
		spec         *v1.ApplicationSpec
		versionMap   map[string]*v1.VersionedApplicationSpec
		version      *v1.VersionedApplicationSpec
		inputs       render.ValuesInputs
		testManifest TestManifest
		labels       map[string]string
	)

	BeforeEach(func() {
		spec = test.LoadExtensionSpec("../spec.yaml")
		versionMap = make(map[string]*v1.VersionedApplicationSpec)
		for _, version := range spec.Versions {
			versionMap[version.Version] = version
		}
		inputs = render.ValuesInputs{
			Name:             name,
			FlavorName:       meshName,
			InstallNamespace: namespace,
			MeshRef: core.ResourceRef{
				Name:      meshName,
				Namespace: namespace,
			},
		}
		labels     = map[string]string {
			"app": "kiali",
		}
	})

	bindVersion := func(versionString string) {
		version = versionMap[versionString]
		inputs.SpecDefinedValues = version.ValuesYaml
		rendered, err := render.ComputeResourcesForApplication(context.TODO(), inputs, version)
		Expect(err).NotTo(HaveOccurred())
		testManifest = NewTestManifestWithResources(rendered)
	}

	testDemoSecret := func() {
		rb := ResourceBuilder{
			Name: name,
			Namespace: namespace,
			Data: map[string]string {
				"username": "admin",
				"passphrase": "admin",
			},
			Labels: labels,
		}
		testManifest.ExpectSecret(rb.GetSecret())
	}

	Context("0.16 with default values", func() {
		BeforeEach(func() {
			bindVersion("0.16")
			labels = map[string]string {
				"chart": "kiali",
				"heritage": "Tiller",
				"release": "kiali",
				"app": "kiali",
			}
		})

		It("has the correct number of resources", func() {
			Expect(testManifest.NumResources()).To(Equal(8))
		})

		It("has a demo secret", func() {
			testDemoSecret()
		})
	})

	Context("0.12 with default values", func() {
		BeforeEach(func() {
			bindVersion("0.12")
		})

		It("has the correct number of resources", func() {
			Expect(testManifest.NumResources()).To(Equal(7))
		})

		It("has a demo secret", func() {
			testDemoSecret()
		})
	})
})
