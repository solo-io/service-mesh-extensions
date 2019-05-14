package options

import "context"

type Options struct {
	Ctx      context.Context
	Validate Validate
}

type Validate struct {
	ManifestFilepath string
	ExtensionName    string
	Version          string
	Flavor           string
	PrintManifest    bool
	MeshName         string
	MeshNamespace    string
	InstallNamespace string
}

var ValidateDefaults = Validate{
	ManifestFilepath: "",
	ExtensionName:    "",
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
