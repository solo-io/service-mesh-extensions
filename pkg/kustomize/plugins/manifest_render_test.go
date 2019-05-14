package plugins

import (
	"fmt"

	"github.com/solo-io/service-mesh-hub/pkg/internal/test"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
)

var _ = Describe("manifest render plugin", func() {
	const (
		emptyYaml     = ``
		invalidYaml_0 = `
	adsa
	asd
`
		invalidYaml_1 = `

  apiVersion: supergloo.solo.io/v1: 32
  kind: MeshIngress
  metadata:
    name: gloo
    namespace: {{ .SuperglooNamespace }}
`
		unknownVariables = `
  apiVersion: supergloo.solo.io/v1
  kind: MeshIngress
  metadata:
    name: gloo
    namespace: {{ .SupergloNamespace }}
`
		validManifest = `
apiVersion: v1
kind: Pod
metadata:
  annotations:
    installNamespace: {{ .InstallNamespace }}
    superglooNamespace: {{ .SuperglooNamespace }}
    customValue: {{ .Custom.SomeValue }}
  name: {{ .MeshRef.Name }}
  namespace: {{ .MeshRef.Namespace }}
spec: {}

`
	)
	var (
		rf     *resmap.Factory
		plugin *manifestRenderPlugin
		unst   unstructured.Unstructured
		res    *resource.Resource
		values = test.GetRenderValues()
	)

	BeforeEach(func() {
		rf = test.ResourceMapFactory()
		plugin = NewManifestRenderPlugin(values)
		unst = unstructured.Unstructured{}
	})

	It("returns name properly", func() {
		Expect(plugin.Name()).To(Equal(ManifestRenderPluginName))
	})

	Context("config", func() {
		It("returns an error if manifest is not valid yaml", func() {
			unst.SetUnstructuredContent(map[string]interface{}{
				"manifest": emptyYaml,
			})
			res = rf.RF().FromMap(unst.Object)
			err := plugin.Config(nil, rf, res)
			Expect(err).To(HaveOccurred())
		})
		It("returns no error if yaml is valid", func() {
			unst.SetUnstructuredContent(map[string]interface{}{
				"manifest": validManifest,
			})
			res = rf.RF().FromMap(unst.Object)
			err := plugin.Config(nil, rf, res)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("generate", func() {
		var testCases = []struct {
			manifest    string
			description string
			err         string
		}{
			{manifest: invalidYaml_0, description: "template variables are incorrect", err: ""},
			{manifest: invalidYaml_1, description: "manifest is not valid yaml", err: "error converting YAML to JSON"},
			{manifest: unknownVariables, description: "template vars are incorrect", err: "can't evaluate field SupergloNamespace"},
		}
		for _, tc := range testCases {
			testCase := tc
			It(fmt.Sprintf("returns an error if %s", testCase.description), func() {
				unst.SetUnstructuredContent(map[string]interface{}{
					"manifest": testCase.manifest,
				})
				res = rf.RF().FromMap(unst.Object)
				err := plugin.Config(nil, rf, res)
				Expect(err).NotTo(HaveOccurred())
				_, err = plugin.Generate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(testCase.err))
			})
		}

		It("creates valid resource when given valid yaml", func() {

			unst.SetUnstructuredContent(map[string]interface{}{
				"manifest": validManifest,
			})
			res = rf.RF().FromMap(unst.Object)
			err := plugin.Config(nil, rf, res)
			Expect(err).NotTo(HaveOccurred())
			res, err := plugin.Generate()
			Expect(err).NotTo(HaveOccurred())
			byt, err := res.EncodeAsYaml()
			Expect(err).NotTo(HaveOccurred())
			var pod corev1.Pod
			err = yaml.Unmarshal(byt, &pod)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod.Name).To(Equal(values.MeshRef.Name))
			Expect(pod.Namespace).To(Equal(values.MeshRef.Namespace))
			Expect(pod.Annotations).To(BeEquivalentTo(map[string]string{
				"superglooNamespace": values.SuperglooNamespace,
				"installNamespace":   values.InstallNamespace,
				"customValue":        (values.Custom.(map[string]interface{})["SomeValue"]).(string),
			}))
		})
	})
})
