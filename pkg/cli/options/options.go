package options

import (
	"context"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

type Options struct {
	Ctx              context.Context
	Validate         Validate
	Registry         Registry
	InstallNamespace string
	InstallSpecFile  string
	ManifestFile     string
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

type Registry struct {
	GithubRegistry v1.GithubRepositoryLocation
}

var RegistryDefaults = Registry{
	GithubRegistry: v1.GithubRepositoryLocation{
		Org:       "solo-io",
		Repo:      "service-mesh-hub",
		Ref:       "better-layering-aug",
		Directory: "meshes/v1",
	},
}

func InitializeOptions(ctx context.Context) *Options {
	opts := &Options{
		Ctx:      ctx,
		Validate: ValidateDefaults,
	}
	return opts
}
