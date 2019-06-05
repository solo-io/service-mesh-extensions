package options

import (
	"context"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

type Options struct {
	Ctx      context.Context
	Validate Validate
}

type Validate struct {
	ManifestFilepath string
	ApplicationName  string
	ApplicationType  string
	Version          string
	Flavor           string
	PrintManifest    bool
	MeshName         string
	MeshNamespace    string
	InstallNamespace string
}

var ValidateDefaults = Validate{
	ManifestFilepath: "",
	ApplicationName:  "",
	ApplicationType:  v1.ApplicationType_EXTENSION.String(),
	Version:          "",
	Flavor:           "",
	PrintManifest:    false,
	MeshName:         "mesh-name",
	MeshNamespace:    "default",
	InstallNamespace: "default",
}

func InitializeOptions(ctx context.Context) *Options {
	opts := &Options{
		Ctx:      ctx,
		Validate: ValidateDefaults,
	}
	return opts
}
