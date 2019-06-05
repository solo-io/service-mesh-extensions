package cli

import (
	"context"

	"github.com/solo-io/go-utils/clicore"
	"github.com/solo-io/service-mesh-hub/pkg/cli/cmd/validate"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/spf13/cobra"
)

var Version = "0.1.0"
var FileLogPathElements = []string{".glooshot", "log"}

var CommandConfig = clicore.CommandConfig{
	Args:                "",
	Command:             App,
	RootErrorMessage:    "error running the service mesh hub utility",
	OutputModeEnvVar:    "SMH_CLI_OUTPUT_MODE",
	LoggingContext:      []interface{}{"version", Version},
	FileLogPathElements: FileLogPathElements,
	Version:             Version,
}

func App(ctx context.Context, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "smh",
		Short:   "Utility for interacting with Service Mesh Hub",
		Version: version,
	}
	o := options.InitializeOptions(ctx)
	cmd.AddCommand(
		validate.Cmd(o))
	return cmd
}
