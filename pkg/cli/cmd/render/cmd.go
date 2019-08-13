package render

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/service-mesh-hub/pkg/cli/installspec"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/solo-io/service-mesh-hub/pkg/registry"
	renderutil "github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/pkg/util"
	"github.com/spf13/cobra"
)

func Cmd(o *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "render a manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return render(o)
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&o.Registry.GithubRegistry.Org, "registry-org", "", options.RegistryDefaults.GithubRegistry.Org,
		"owner of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Repo, "registry-repo", "", options.RegistryDefaults.GithubRegistry.Repo,
		"repo of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Ref, "registry-ref", "", options.RegistryDefaults.GithubRegistry.Ref,
		"ref of github registry")
	pflags.StringVarP(&o.Registry.GithubRegistry.Directory, "registry-directory", "", options.RegistryDefaults.GithubRegistry.Directory,
		"directory of github registry")
	pflags.StringVarP(&o.InstallSpecFile, "install-spec-file", "i", "",
		"optional install spec to generate manifests from")
	pflags.StringVarP(&o.ManifestFile, "manifest-file", "m", "",
		"optional destination for rendered manifest, otherwise print to stdout")
	return cmd
}

func render(o *options.Options) error {
	var installSpec *installspec.InstallSpec
	var err error
	if o.InstallSpecFile == "" {
		reader := registry.NewGithubSpecReader(o.Ctx, o.Registry.GithubRegistry)
		if installSpec, err = installspec.GetInstallSpec(reader, o.InstallNamespace); err != nil {
			return err
		}
	} else {
		installSpec = &installspec.InstallSpec{}
		if err = installSpec.Load(o.InstallSpecFile); err != nil {
			return err
		}
	}

	manifest, err := renderManifest(o.Ctx, installSpec)
	if err != nil {
		return err
	}

	if o.ManifestFile != "" {
		return util.SaveFile(o.ManifestFile, manifest)
	}
	fmt.Print(manifest + "\n")
	return nil
}

func renderManifest(ctx context.Context, spec *installspec.InstallSpec) (string, error) {
	resources, err := renderutil.ComputeResourcesForApplication(ctx, spec.Values, spec.Version)
	if err != nil {
		return "", err
	}
	manifests, err := helmchart.ManifestsFromResources(resources)
	if err != nil {
		return "", err
	}
	return manifests.CombinedString() + "\n", nil
}
