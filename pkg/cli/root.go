package cli

import (
	"context"

	"github.com/solo-io/go-utils/clicore"
	"github.com/solo-io/service-mesh-hub/pkg/cli/cmd/prepare"
	"github.com/solo-io/service-mesh-hub/pkg/cli/cmd/render"
	"github.com/solo-io/service-mesh-hub/pkg/cli/cmd/validate"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/solo-io/service-mesh-hub/pkg/internal/version"
	"github.com/spf13/cobra"
)

var FileLogPathElements = []string{".hubctl", "log"}

var CommandConfig = clicore.CommandConfig{
	Args:                "",
	Command:             App,
	RootErrorMessage:    "error running the service mesh hub utility",
	OutputModeEnvVar:    "HUBCTL_OUTPUT_MODE",
	LoggingContext:      []interface{}{"version", version.Version},
	FileLogPathElements: FileLogPathElements,
	Version:             version.Version,
}

func App(ctx context.Context, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hubctl",
		Short:   "Utility for interacting with Service Mesh Hub",
		Version: version,
	}
	o := options.InitializeOptions(ctx)
	cmd.AddCommand(
		prepare.Cmd(o),
		render.Cmd(o),
		validate.Cmd(o))
	return cmd
}
