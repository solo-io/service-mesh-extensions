package validate

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/protoutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func Cmd(o *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate a local spec file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(o)
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.StringVar(&o.Validate.ApplicationType, "type", options.ValidateDefaults.ApplicationType,
		fmt.Sprintf("type of the application that will be validated. Available: %v, %v, %v",
			v1.ApplicationType_EXTENSION.String(), v1.ApplicationType_DEMO.String(), v1.ApplicationType_MESH.String()))
	pflags.StringVar(&o.Validate.ApplicationName, "name", options.ValidateDefaults.ApplicationName,
		"name of the application that will be validated")
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

// map the manifest type to the directory
var directoryNamaes = map[string]string{
	v1.ApplicationType_EXTENSION.String(): "extensions",
	v1.ApplicationType_DEMO.String():      "demos",
	v1.ApplicationType_MESH.String():      "meshes",
}

func validate(o *options.Options) error {

	if o.Validate.ApplicationName == "" {
		return fmt.Errorf("no extension name specified")
	}
	appDir, ok := directoryNamaes[o.Validate.ApplicationType]
	if !ok {
		return fmt.Errorf("must provide a valid application type with --type; available options: "+
			"%v, %v, %v", v1.ApplicationType_EXTENSION.String(), v1.ApplicationType_DEMO.String(), v1.ApplicationType_MESH.String())
	}
	if o.Validate.Flavor == "" {
		return fmt.Errorf("no flavor specified")
	}

	specFilepath := fmt.Sprintf("./%v/v1/%v/spec.yaml", appDir, o.Validate.ApplicationName)
	spec, err := LoadExtensionSpec(specFilepath)
	if err != nil {
		return err
	}
	versionContent, err := getVersionContent(spec.Versions, o.Validate.Version)
	if err != nil {
		return err
	}
	flavorContent, err := getFlavorContent(versionContent, o.Validate.Flavor)
	if err != nil {
		return err
	}
	inputValues := render.ValuesInputs{
		Name:             o.Validate.ApplicationName,
		Flavor:           flavorContent,
		InstallNamespace: o.Validate.InstallNamespace,
		MeshRef: core.ResourceRef{
			Namespace: o.Validate.MeshNamespace,
			Name:      o.Validate.MeshName,
		},
		SpecDefinedValues: versionContent.ValuesYaml,
		// TODO - support validation with these parameters
		//UserDefinedValues:  "",
		//Params:       nil,
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

func getFlavorContent(version *v1.VersionedApplicationSpec, name string) (*v1.Flavor, error) {
	for _, content := range version.Flavors {
		if name == content.Name {
			return content, nil
		}
	}
	return nil, fmt.Errorf("could not find flavor %v in specification", name)
}
