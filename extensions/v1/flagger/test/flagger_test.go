package test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	. "github.com/solo-io/go-utils/manifesttestutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/test"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("flagger", func() {

	const (
		meshNamespace            = "istio-system"
		name                     = "flagger"
		meshName                 = "my-istio"
		superglooIstioFlavor     = "istio-supergloo"
		superglooClusterRoleName = "supergloo-cluster-role"
		flaggerSgCrbName         = "flagger-supergloo"
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

	Context("0.12.0 with istio-supergloo flavor (default parameters)", func() {
		var (
			version      *v1.VersionedApplicationSpec
			inputs       render.ValuesInputs
			testManifest TestManifest
			rendered     kuberesource.UnstructuredResources
			err          error
		)

		BeforeEach(func() {
			version = versionMap["0.12.0"]
			layers := []render.LayerInput{{LayerId: superglooIstioFlavor, OptionId: "cluster-role-binding"}}
			inputs = render.ValuesInputs{
				Name:             name,
				Flavor:           test.GetFlavor(superglooIstioFlavor, version),
				InstallNamespace: meshNamespace,
				MeshRef: core.ResourceRef{
					Name:      meshName,
					Namespace: meshNamespace,
				},
				SpecDefinedValues: version.ValuesYaml,
				Params:            test.GetDefaultParameters(version, superglooIstioFlavor, layers),
				Layers:            layers,
			}

			rendered, err = render.ComputeResourcesForApplication(context.TODO(), inputs, version)
			Expect(err).NotTo(HaveOccurred())
			testManifest = NewTestManifestWithResources(rendered)
		})

		It("has the correct number of resources", func() {
			// 5 resources from the flagger chart, plus 1 additional ClusterRoleBinding
			// required by flagger to control supergloo resources
			Expect(testManifest.NumResources()).To(Equal(6))
		})

		It("the flagger service account is correctly bound to the supergloo cluster role", func() {
			superglooCrb := rendered.Filter(func(resource *unstructured.Unstructured) bool {
				return !(resource.GetKind() == "ClusterRoleBinding" && resource.GetName() == flaggerSgCrbName)
			})

			Expect(superglooCrb).To(HaveLen(1))

			obj, err := kuberesource.ConvertUnstructured(superglooCrb[0])
			Expect(err).NotTo(HaveOccurred())

			crb, ok := obj.(*rbacv1.ClusterRoleBinding)
			Expect(ok).To(BeTrue())

			Expect(crb.Name).To(Equal(flaggerSgCrbName))
			Expect(crb.RoleRef.Name).To(Equal(superglooClusterRoleName))
			Expect(crb.RoleRef.Kind).To(Equal("ClusterRole"))
			Expect(crb.Subjects).To(HaveLen(1))
			Expect(crb.Subjects[0].Kind).To(Equal("ServiceAccount"))
			Expect(crb.Subjects[0].Name).To(Equal(name))
			Expect(crb.Subjects[0].Namespace).To(Equal(meshNamespace))
		})

		It("the flagger pods get started with the expected arguments", func() {
			flaggerDeployment := rendered.Filter(func(resource *unstructured.Unstructured) bool {
				return !(resource.GetKind() == "Deployment" && resource.GetName() == name)
			})

			Expect(flaggerDeployment).To(HaveLen(1))

			obj, err := kuberesource.ConvertUnstructured(flaggerDeployment[0])
			Expect(err).NotTo(HaveOccurred())

			deployment, ok := obj.(*appsv1.Deployment)
			Expect(ok).To(BeTrue())

			expectedMeshProviderArg := fmt.Sprintf("-mesh-provider=supergloo:%s.%s", meshName, meshNamespace)
			expectedPrometheusArg := fmt.Sprintf("-metrics-server=http://prometheus.%s:9090", meshNamespace)
			Expect(deployment.Spec.Template.Spec.Containers).To(HaveLen(1))
			Expect(deployment.Spec.Template.Spec.Containers[0].Command).To(ContainElement(expectedMeshProviderArg))
			Expect(deployment.Spec.Template.Spec.Containers[0].Command).To(ContainElement(expectedPrometheusArg))

		})
	})
})
