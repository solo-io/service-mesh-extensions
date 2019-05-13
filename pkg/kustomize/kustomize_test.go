package kustomize_test

import (
	"context"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/internal/test"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/plugins"
	"github.com/solo-io/service-mesh-hub/pkg/render"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/spf13/afero"
)

type incorrectPlugin struct{}

func (p *incorrectPlugin) Name() string {
	return "incorrect"
}

var _ = Describe("kustomize unit tests", func() {

	var (
		fs         afero.Fs
		ctx        context.Context
		tempDir    string
		manifests  helmchart.Manifests
		pathloader loader.Loader
		testNs     = "test"
		appSpec    *v1.VersionedApplicationSpec
		inputs     render.ValuesInputs
	)

	BeforeSuite(func() {
		appSpec = test.GetAppSpec(test.HelloWorldChart1_0, nil)
	})

	BeforeEach(func() {
		var err error
		fs = afero.NewOsFs()
		tempDir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())
		inputs = render.ValuesInputs{
			Name: test.HelloWorldChart1_0.Name,
			InstallNamespace: testNs,
			FlavorName: test.DefaultFlavorName,
		}
		manifests, err = render.GetManifestsFromApplicationSpec(ctx, inputs, appSpec)
		Expect(err).NotTo(HaveOccurred())
		pathloader = loader.NewKustomizeLoader(ctx, fs)
	})

	AfterEach(func() {
		err := fs.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("e2e", func() {
		BeforeEach(func() {
			pathloader = loader.NewKustomizeLoader(ctx, fs)
		})

		It("can work e2e", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "supergloo")
			Expect(err).NotTo(HaveOccurred())
			k, err := kustomize.NewKustomizer(pathloader, manifests, klayer, plugins.NewManifestRenderPlugin(test.GetRenderValues()))
			Expect(err).NotTo(HaveOccurred())
			bytes, err := k.Run(tempDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytes).NotTo(BeEmpty())
		})

		It("errors with an incorrect resource name", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "error")
			Expect(err).NotTo(HaveOccurred())
			k, err := kustomize.NewKustomizer(pathloader, manifests, klayer)
			Expect(err).NotTo(HaveOccurred())
			_, err = k.Run(tempDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("was never loaded"))
		})

		It("errors with no correct base path", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "fails")
			Expect(err).NotTo(HaveOccurred())
			k, err := kustomize.NewKustomizer(pathloader, manifests, klayer)
			Expect(err).NotTo(HaveOccurred())
			_, err = k.Run(tempDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fails: no such file or directory"))
		})

		It("fails when an incorrect object is registered as a plugin", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "fails")
			Expect(err).NotTo(HaveOccurred())

			_, err = kustomize.NewKustomizer(pathloader, manifests, klayer, &incorrectPlugin{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid kustomize plugin"))
		})

	})

})
