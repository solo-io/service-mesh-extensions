package registry_test

import (
	"context"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/registry"

	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {

	var reader registry.SpecReader

	PDescribe("ReaderTest", func() {
		Describe("RemoteSpecReader", func() {
			Describe("GetSpecs", func() {
				It("works", func() {
					reader = registry.NewRemoteSpecReader(context.TODO(), DefaultRemoteUrl)
					specs, err := reader.GetSpecs()
					Expect(err).NotTo(HaveOccurred())
					Expect(len(specs) > 0).To(BeTrue())
				})

				It("errors with a bad url", func() {
					reader = registry.NewRemoteSpecReader(context.TODO(), "fake-url")
					_, err := reader.GetSpecs()
					Expect(err).To(HaveOccurred())
					expectedErr := registry.FailedToDownloadAppSpecsError(errors.Errorf(""))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe("GithubSpecReader", func() {
			Describe("GetSpecs", func() {
				It("works", func() {
					chartsRef := v1.GithubRepositoryLocation{
						Org:       "solo-io",
						Repo:      "service-mesh-hub",
						Ref:       "v0.0.1-initial-api",
						Directory: "extensions/v1",
					}
					reader = registry.NewGithubSpecReader(context.TODO(), chartsRef)
					specs, err := reader.GetSpecs()
					Expect(err).NotTo(HaveOccurred())
					Expect(len(specs)).To(Equal(7))
				})

				It("errors with a bad repo", func() {
					chartsRef := v1.GithubRepositoryLocation{
						Org:       "solo-io",
						Repo:      "service-mesh-hub",
						Ref:       "v0.0.1-initial-api",
						Directory: "fake-directory",
					}
					reader = registry.NewGithubSpecReader(context.TODO(), chartsRef)
					_, err := reader.GetSpecs()
					Expect(err).To(HaveOccurred())
					expectedErr := registry.FailedToGetSpecsFromGithubError(errors.Errorf(""))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})
	})
})
