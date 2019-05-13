package test

import (
	"path/filepath"

	"github.com/solo-io/go-utils/tarutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/spf13/afero"
)

func NewKustomizeTestLayerFromLocalPackages(fs afero.Fs, dir, overlayPath string) (*v1.Kustomize, error) {
	tempFile, err := afero.TempFile(fs, "", "")
	if err != nil {
		return nil, err
	}
	absolutePath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	err = tarutils.Tar(absolutePath, fs, tempFile)
	if err != nil {
		return nil, err
	}
	kt := &v1.Kustomize{
		OverlayPath: overlayPath,
		Location: &v1.Kustomize_TgzArchive{
			TgzArchive: &v1.TgzLocation{
				Uri: tempFile.Name(),
			},
		},
	}
	return kt, nil
}
