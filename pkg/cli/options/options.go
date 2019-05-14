package options

import "context"

type Options struct {
	Ctx      context.Context
	Validate Validate
}

type Validate struct {
	ManifestFilepath string
	ExtensionName    string
	VersionIndex     int
	Flavor           string
	Verbose          bool
}

func InitializeOptions(ctx context.Context) *Options {
	opts := &Options{
		Ctx: ctx,
		Validate: Validate{
			Verbose:          false,
			ManifestFilepath: "",
			ExtensionName:    "",
			VersionIndex:     0,
			Flavor:           "",
		},
	}
	return opts
}
