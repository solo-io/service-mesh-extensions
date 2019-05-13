package kustomize_test

import (
	"context"

	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"
	"github.com/solo-io/sm-marketplace/services/operator/kustomize"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/spf13/afero"
)

var _ = Describe("kustomize unit tests", func() {

	var (
		fs         afero.Fs
		ctx        context.Context
		tempDir    string
		manifests  helmchart.Manifests
		pathloader loader.Loader
		testNs     = "test"
		//appState   *v1.ApplicationState
	)

	BeforeSuite(func() {
		appState = test.GetState(test.HelloWorldChart1_0, testNs, test.InstallState_Vanilla, nil)
	})

	BeforeEach(func() {
		var err error
		fs = afero.NewOsFs()
		tempDir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())
		manifests, err = utils.GetManifestsFromApplicationSpec(ctx, appState)
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
			appState = test.GetState(test.HelloWorldChart1_0, testNs, test.InstallState_Vanilla, klayer)
			k := kustomize.NewKustomizer(pathloader, manifests, klayer, test.InstallState_Vanilla)
			bytes, err := k.Run(tempDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(bytes).NotTo(BeEmpty())
		})

		It("errors with an incorrect resource name", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "error")
			Expect(err).NotTo(HaveOccurred())
			k := kustomize.NewKustomizer(pathloader, manifests, klayer, test.InstallState_Vanilla)
			_, err = k.Run(tempDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("was never loaded"))
		})

		It("errors with no correct base path", func() {
			klayer, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "fixtures", "fails")
			Expect(err).NotTo(HaveOccurred())
			k := kustomize.NewKustomizer(pathloader, manifests, klayer, test.InstallState_Vanilla)
			_, err = k.Run(tempDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fails: no such file or directory"))
		})

	})

})
