package plugins

import (
	"fmt"

	"github.com/solo-io/service-mesh-hub/pkg/internal/mocks"
	"github.com/solo-io/service-mesh-hub/pkg/internal/test"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
)

var _ = Describe("loader", func() {
	var (
		loader *staticPluginLoader
		ctrl   *gomock.Controller
		res    *resource.Resource
		unst   unstructured.Unstructured
		rf     *resmap.Factory
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(T)
		rf = test.ResourceMapFactory()
	})

	AfterEach(func() {
		defer ctrl.Finish()
	})

	Context("generators", func() {
		var (
			generatorName = "GeneratorName"
			mockGenerator *mocks.MockNamedGenerator
		)
		BeforeEach(func() {
			mockGenerator = mocks.NewMockNamedGenerator(ctrl)
			mockGenerator.EXPECT().Generate().Times(0)
			unst.SetKind(generatorName)
			res = rf.RF().FromMap(unst.Object)
		})
		It("can load a generator which has been initialized", func() {
			mockGenerator.EXPECT().Name().Times(1).Return(generatorName)
			loader = NewStaticPluginLoader([]NamedGenerator{mockGenerator}, nil)
			result, err := loader.LoadGenerator(nil, res)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(mockGenerator))
		})

		It("fails to load a generator that was never loaded", func() {
			mockGenerator.EXPECT().Name().Times(0)
			loader = NewStaticPluginLoader([]NamedGenerator{}, nil)
			_, err := loader.LoadGenerator(nil, res)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("plugin %s was never loaded", generatorName)))
		})
	})

	Context("transformers", func() {
		var (
			transformerName = "TransformerName"
			mockTransformer *mocks.MockNamedTransformer
		)
		BeforeEach(func() {
			mockTransformer = mocks.NewMockNamedTransformer(ctrl)
			mockTransformer.EXPECT().Transform(gomock.Any()).Times(0)
			unst.SetKind(transformerName)
			res = rf.RF().FromMap(unst.Object)
		})
		It("can load a generator which has been initialized", func() {
			mockTransformer.EXPECT().Name().Times(1).Return(transformerName)
			loader = NewStaticPluginLoader(nil, []NamedTransformer{mockTransformer})
			result, err := loader.LoadTransformer(nil, res)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(mockTransformer))
		})

		It("fails to load a generator that was never loaded", func() {
			mockTransformer.EXPECT().Name().Times(0)
			loader = NewStaticPluginLoader(nil, []NamedTransformer{})
			_, err := loader.LoadTransformer(nil, res)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("plugin %s was never loaded", transformerName)))
		})
	})
})
