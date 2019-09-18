package options

import (
	"path/filepath"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/service-mesh-hub/pkg/registry"
	"go.uber.org/zap"
)

func MustGetSpecReader(o *Options) registry.SpecReader {
	if o.Registry.LocalDirectory == "" {
		return registry.NewGithubSpecReader(o.Ctx, o.Registry.GithubRegistry)
	}

	absPath, err := filepath.Abs(o.Registry.LocalDirectory)
	if err != nil {
		contextutils.LoggerFrom(o.Ctx).Fatalw("Failed to get absolute path", zap.Error(err))
	}

	return registry.NewLocalSpecReader(o.Ctx, absPath)
}
