package loader

//go:generate mockgen -package=mocks -destination=../../internal/mocks/loader.go github.com/solo-io/service-mesh-hub/pkg/kustomize/loader Loader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/google/go-github/github"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/githubutils"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/tarutils"
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/spf13/afero"
)

const (
	base         = "base"
	resourceYaml = "resource.yaml"
	kustYaml     = "kustomization.yaml"
)

type Loader interface {
	RetrieveLayers(dir string, kustomize *hubv1.Kustomize) (string, error)
	LoadBase(manifests helmchart.Manifests, dir string) error
}

type kustomizeLoader struct {
	ctx context.Context
	fs  afero.Fs
}

func NewKustomizeLoader(ctx context.Context, fs afero.Fs) *kustomizeLoader {
	return &kustomizeLoader{ctx: ctx, fs: fs}
}

func (kl *kustomizeLoader) RetrieveLayers(dir string, kustomize *hubv1.Kustomize) (string, error) {
	switch kustType := kustomize.GetLocation().(type) {
	case *hubv1.Kustomize_Github:
		return kl.githubRetrieve(dir, kustType)
	case *hubv1.Kustomize_TgzArchive:
		return kl.tgzRetrieve(dir, kustType)
	}
	return "", errors.Errorf("no kustomize location option found")
}

func (kl *kustomizeLoader) githubRetrieve(dir string, kustType *hubv1.Kustomize_Github) (string, error) {
	logger := contextutils.LoggerFrom(kl.ctx)
	client, err := githubutils.GetClient(kl.ctx)
	if err != nil {
		logger.Debugf("no GITHUB_TOKEN found in env, using basic client")
		client = github.NewClient(http.DefaultClient)
	}

	owner, repo, ref := kustType.Github.Org, kustType.Github.Repo, kustType.Github.Ref

	tmpf, err := afero.TempFile(kl.fs, dir, "*.tar.gz")
	if err != nil {
		logger.Errorf("can't create temp file")
		return "", err
	}
	defer kl.fs.Remove(tmpf.Name())

	if err := githubutils.DownloadRepoArchive(kl.ctx, client, tmpf, owner, repo, ref); err != nil {
		return "", err
	}
	if err := tarutils.Untar(dir, tmpf.Name(), kl.fs); err != nil {
		return "", err
	}

	repoFolder, err := getRepoFolder(kl.fs, dir)
	if err != nil {
		return "", err
	}
	pluginDir := filepath.Join(repoFolder, kustType.Github.Directory)

	return pluginDir, nil
}

func (kl *kustomizeLoader) tgzRetrieve(dir string, kustType *hubv1.Kustomize_TgzArchive) (string, error) {
	logger := contextutils.LoggerFrom(kl.ctx)

	archive, err := tarutils.RetrieveArchive(kl.fs, kustType.TgzArchive.GetUri())
	if err != nil {
		return "", err
	}
	defer archive.Close()

	tmpf, err := afero.TempFile(kl.fs, dir, "*.tar.gz")
	if err != nil {
		logger.Errorf("can't create temp file")
		return "", err
	}
	defer kl.fs.Remove(tmpf.Name())

	_, err = io.Copy(tmpf, archive)
	if err != nil {
		return "", err
	}

	err = tarutils.Untar(dir, tmpf.Name(), kl.fs)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func getRepoFolder(fs afero.Fs, tmpd string) (string, error) {
	files, err := afero.ReadDir(fs, tmpd)
	if err != nil {
		return "", err
	}

	var repoDirName string
	for _, file := range files {
		if file.IsDir() {
			repoDirName = file.Name()
			break
		}
	}
	if repoDirName == "" {
		return "", fmt.Errorf("unable to find directory in archive of git repo")
	}
	return filepath.Join(tmpd, repoDirName), nil
}

// loads the kustomize base into the directory which you specify
func (kl *kustomizeLoader) LoadBase(manifests helmchart.Manifests, dir string) error {
	bases, err := gatherBases(kl.fs, dir)
	if err != nil {
		return err
	}
	kustOptions := &types.Kustomization{
		Bases: bases,
	}
	if err := writeKustomizationFile(kl.fs, kustOptions, dir); err != nil {
		return err
	}

	baseDir := filepath.Join(dir, base)
	err = kl.fs.Mkdir(baseDir, 0777)
	if err != nil {
		return err
	}
	err = afero.WriteFile(kl.fs, filepath.Join(baseDir, resourceYaml), []byte(manifests.CombinedString()), 0777)
	if err != nil {
		return err
	}
	kustOptions = &types.Kustomization{
		Resources: []string{resourceYaml},
	}
	if err := writeKustomizationFile(kl.fs, kustOptions, baseDir); err != nil {
		return err
	}
	return nil
}

func writeKustomizationFile(fs afero.Fs, kustOptions *types.Kustomization, dir string) error {
	byt, err := yaml.Marshal(kustOptions)
	if err != nil {
		return err
	}

	err = afero.WriteFile(fs, filepath.Join(dir, kustYaml), byt, 0777)
	if err != nil {
		return err
	}
	return nil
}

// gathers up the bases which have been downloaded
func gatherBases(fs afero.Fs, dir string) ([]string, error) {
	info, err := afero.ReadDir(fs, dir)
	if err != nil {
		return nil, err
	}
	var bases []string
	for _, v := range info {
		if v.IsDir() {
			bases = append(bases, fmt.Sprintf("./%s", filepath.Base(v.Name())))
		}
	}
	return bases, nil
}
