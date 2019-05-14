package validate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/spf13/cobra"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

func Cmd(o *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate a manifest file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(o)
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVar(&o.Validate.ExtensionName, "name", options.ValidateDefaults.ExtensionName,
		"name of the extension that will be validated")
	pflags.StringVar(&o.Validate.Version, "version", options.ValidateDefaults.Version,
		"specification version to be validated")
	pflags.StringVar(&o.Validate.Flavor, "flavor", options.ValidateDefaults.Flavor,
		"name of flavor to be validate")
	pflags.BoolVar(&o.Validate.PrintManifest, "print-manifest", options.ValidateDefaults.PrintManifest,
		"if set, will print the manifest content")
	pflags.StringVar(&o.Validate.InstallNamespace, "install-namespace", options.ValidateDefaults.InstallNamespace,
		fmt.Sprintf("optional, namespace in which to install the app, defaults to placeholder value: %v", options.ValidateDefaults.InstallNamespace))
	pflags.StringVar(&o.Validate.MeshName, "mesh-name", options.ValidateDefaults.MeshName,
		fmt.Sprintf("optional, name of the associated mesh, defaults to placeholder value: %v", options.ValidateDefaults.MeshName))
	pflags.StringVar(&o.Validate.MeshNamespace, "mesh-namespace", options.ValidateDefaults.MeshNamespace,
		fmt.Sprintf("optional, namespace of the associated mesh, defaults to placeholder value: %v", options.ValidateDefaults.MeshNamespace))
	return cmd
}

func validate(o *options.Options) error {

	if o.Validate.ExtensionName == "" {
		return fmt.Errorf("no extension name specified")
	}
	if o.Validate.Flavor == "" {
		return fmt.Errorf("no flavor specified")
	}
	specFilepath := fmt.Sprintf("./extensions/v1/%v/spec.yaml", o.Validate.ExtensionName)
	spec, err := LoadExtensionSpec(specFilepath)
	if err != nil {
		return err
	}
	versionContent, err := getVersionContent(spec.Versions, o.Validate.Version)
	if err != nil {
		return err
	}
	inputValues := render.ValuesInputs{
		Name:             o.Validate.ExtensionName,
		FlavorName:       o.Validate.Flavor,
		InstallNamespace: o.Validate.InstallNamespace,
		MeshRef: core.ResourceRef{
			Namespace: o.Validate.MeshNamespace,
			Name:      o.Validate.MeshName,
		},
		SpecDefinedValues: versionContent.ValuesYaml,
		// TODO - support validation with these parameters
		//SuperglooNamespace: "",
		//UserDefinedValues:  "",
		//FlavorParams:       nil,
		//SpecDefinedValues:  "",
	}

	resources, err := render.ComputeResourcesForApplication(o.Ctx, inputValues, versionContent)
	if err != nil {
		return errors.Wrapf(err, "unable to compute resources on version %v", o.Validate.Version)
	}
	if o.Validate.PrintManifest {
		for _, r := range resources {
			var b []byte
			var err error
			if b, err = r.MarshalJSON(); err != nil {
				return errors.Wrapf(err, "unable to unmarshal unstructured resource")
			}
			if b, err = yaml.JSONToYAML(b); err != nil {
				return err
			}
			fmt.Printf("%v---\n", string(b))
		}
	}
	return nil
}

func LoadExtensionSpec(pathToSpec string) (*v1.ApplicationSpec, error) {
	bytes, err := ioutil.ReadFile(pathToSpec)
	if err != nil {
		return nil, err
	}
	var spec v1.ApplicationSpec
	err = protoutils.UnmarshalYaml(bytes, &spec)
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func getVersionContent(list []*v1.VersionedApplicationSpec, version string) (*v1.VersionedApplicationSpec, error) {
	for _, content := range list {
		if content.Version == version {
			return content, nil
		}
	}
	return nil, fmt.Errorf("could not find version %v in specification", version)
}
