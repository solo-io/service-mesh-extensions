package loader

import (
	"context"
	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/installutils"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/service-mesh-hub/pkg/internal/test"
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/spf13/afero"
	"path/filepath"
	"sigs.k8s.io/kustomize/pkg/types"
)

var _ = Describe("kustomize loader", func() {

	var (
		fs      afero.Fs
		ctx     context.Context
		tempDir string
	)

	BeforeEach(func() {
		var err error
		fs = afero.NewOsFs()
		tempDir, err = afero.TempDir(fs, "", "")
		Expect(err).NotTo(HaveOccurred())
		ctx = context.TODO()
	})

	AfterEach(func() {
		err := fs.RemoveAll(tempDir)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("base", func() {
		It("can properly load the directory structure for kustomize", func() {
			manifestSlice, err := installutils.GetManifestsFromRemoteTar("https://github.com/XiaoMi/naftis/releases/download/0.1.4-rc6/manifest.tar.gz")
			Expect(err).NotTo(HaveOccurred())
			manifests := helmchart.Manifests(manifestSlice)
			kloader := NewKustomizeLoader(ctx, fs)
			err = kloader.LoadBase(manifests, tempDir)
			Expect(err).NotTo(HaveOccurred())

			files, err := afero.ReadDir(fs, filepath.Join(tempDir, base))
			Expect(err).NotTo(HaveOccurred())
			foundFiles := 0
			for _, file := range files {
				switch file.Name() {
				case kustYaml:
					byt, err := afero.ReadFile(fs, filepath.Join(tempDir, base, file.Name()))
					Expect(err).NotTo(HaveOccurred())
					var kustOptions types.Kustomization
					err = yaml.Unmarshal(byt, &kustOptions)
					Expect(err).NotTo(HaveOccurred())
					Expect(kustOptions).To(BeEquivalentTo(types.Kustomization{
						Resources: []string{resourceYaml},
					}))
					foundFiles++
				case resourceYaml:
					byt, err := afero.ReadFile(fs, filepath.Join(tempDir, base, file.Name()))
					Expect(err).NotTo(HaveOccurred())
					Expect(string(byt)).To(Equal(manifests.CombinedString()))
					foundFiles++
				default:
					Fail("unintended file found")
				}
			}
			Expect(foundFiles).To(Equal(2))
		})
	})

	Context("remote plugins", func() {
		It("tgz", func() {
			kt, err := test.NewKustomizeTestLayerFromLocalPackages(fs, "../fixtures", "supergloo")
			kl := NewKustomizeLoader(ctx, fs)
			_, err = kl.RetrieveLayers(tempDir, kt)
			Expect(err).NotTo(HaveOccurred())
			files, err := afero.ReadDir(fs, tempDir)
			foundFiles := 0
			Expect(err).NotTo(HaveOccurred())
			for _, file := range files {
				switch filepath.Base(file.Name()) {
				case "supergloo":
					foundFiles++
				case "error":
					foundFiles++
				default:
					Fail("unintended file found")
				}
			}
			Expect(foundFiles).To(Equal(2))
		})
		It("github", func() {
			var (
				guthubDir = "services/operator/kustomize/fixtures"
				githubRef = "64ee5fea427b3bdf50860ac1621d06d0e64f54ce"
			)
			kt := &hubv1.Kustomize{
				Location: &hubv1.Kustomize_Github{
					Github: &hubv1.GithubRepositoryLocation{
						Ref:       githubRef,
						Repo:      "sm-marketplace",
						Org:       "solo-io",
						Directory: guthubDir,
					},
				},
			}
			kl := NewKustomizeLoader(ctx, fs)
			newDir, err := kl.RetrieveLayers(tempDir, kt)
			Expect(err).NotTo(HaveOccurred())
			Expect(newDir).To(ContainSubstring(guthubDir))
			files, err := afero.ReadDir(fs, newDir)
			foundFiles := 0
			Expect(err).NotTo(HaveOccurred())
			for _, file := range files {
				switch filepath.Base(file.Name()) {
				case "supergloo":
					foundFiles++
				case "error":
					foundFiles++
				default:
					Fail("unintended file found")
				}
			}
			Expect(foundFiles).To(Equal(2))
		})
	})

})
