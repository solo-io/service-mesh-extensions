package registry_test

import (
	"context"
	registry "github.com/solo-io/service-mesh-hub/pkg/read"

	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"

	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader", func() {

	var reader registry.SpecReader

	Describe("ReaderTest", func() {
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
					chartsRef := hubv1.GithubRepositoryLocation{
						Org:       "solo-io",
						Repo:      "service-mesh-hub",
						Ref:       "4e9dd4176db09b32ab6de83c12d0ca1908082155",
						Directory: "extensions/v1",
					}
					reader = registry.NewGithubSpecReader(context.TODO(), chartsRef)
					specs, err := reader.GetSpecs()
					Expect(err).NotTo(HaveOccurred())
					Expect(len(specs)).To(Equal(3))
				})

				It("errors with a bad repo", func() {
					chartsRef := hubv1.GithubRepositoryLocation{
						Org:       "solo-io",
						Repo:      "service-mesh-hub",
						Ref:       "4e9dd4176db09b32ab6de83c12d0ca1908082155",
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
