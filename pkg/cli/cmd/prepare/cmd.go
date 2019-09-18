package prepare

import (
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/service-mesh-hub/pkg/cli/installspec"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/spf13/cobra"
)

func Cmd(o *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "prepare an installation spec from which manifests can be rendered",
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepare(o)
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&o.Registry.LocalDirectory, "specs-path", "p", "",
		"local directory to access application specs from, e.g. `extensions/v1`")
	pflags.StringVarP(&o.Registry.GithubRegistry.Org, "registry-org", "", options.RegistryDefaults.GithubRegistry.Org,
		"owner of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Repo, "registry-repo", "", options.RegistryDefaults.GithubRegistry.Repo,
		"repo of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Ref, "registry-ref", "", options.RegistryDefaults.GithubRegistry.Ref,
		"ref of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Directory, "registry-directory", "", options.RegistryDefaults.GithubRegistry.Directory,
		"directory of github registry")
	pflags.StringVarP(&o.InstallNamespace, "namespace", "n", "default",
		"install namespace")
	pflags.StringVarP(&o.InstallSpecFile, "install-spec-file", "i", "",
		"destination for application install spec")
	return cmd
}

func prepare(o *options.Options) error {
	if err := validateOptions(o); err != nil {
		return err
	}

	reader := options.MustGetSpecReader(o)
	installSpec, err := installspec.GetInstallSpec(reader, o.InstallNamespace)
	if err != nil {
		return err
	}

	return installSpec.Save(o.InstallSpecFile)
}

func validateOptions(o *options.Options) error {
	if o.InstallSpecFile == "" {
		return errors.New("-i install spec destination file must be provided")
	}
	return nil
}
