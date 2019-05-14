package validate

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/cli/options"
	"github.com/solo-io/service-mesh-hub/pkg/render"
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
	pflags.StringVar(&o.Validate.ExtensionName, "name", "",
		"name of the extension that will be validated")
	pflags.IntVar(&o.Validate.VersionIndex, "version", 0,
		"index of the file resource to be validated")
	pflags.StringVar(&o.Validate.Flavor, "flavor", "",
		"name of flavor to be validate")
	pflags.BoolVar(&o.Validate.Verbose, "verbose", false,
		"if set, will print the manifest content")
	return cmd
}

func validate(o *options.Options) error {

	if o.Validate.ExtensionName == "" {
		return fmt.Errorf("no extension name specified")
	}
	if o.Validate.Flavor == "" {
		return fmt.Errorf("no flavor specified")
	}
	inputValues := render.ValuesInputs{
		Name:       o.Validate.ExtensionName,
		FlavorName: o.Validate.Flavor,
		// TODO - support validation with these parameters
		InstallNamespace: "default",
		SpecDefinedValues: `
kiali:
  enabled: true
`,
		//		SpecDefinedValues: `
		//    customizationLayers:
		//    - kustomize:
		//        github:
		//          org: solo-io
		//          repo: service-mesh-hub
		//          ref: master
		//          directory: extensions/v1/kiali/overlays
		//        overlayPath: kiali-demo-secret
		//`,
		//MeshRef:            core.ResourceRef{},
		//SuperglooNamespace: "",
		//UserDefinedValues:  "",
		//FlavorParams:       nil,
		//SpecDefinedValues:  "",
	}
	specFilepath := fmt.Sprintf("./extensions/v1/%v/spec.yaml", o.Validate.ExtensionName)
	spec, err := LoadExtensionSpec(specFilepath)
	if err != nil {
		return err
	}

	if o.Validate.VersionIndex > len(spec.Versions) {
		return fmt.Errorf("must specify a valid version index, %v exceeds maximum version index: %v",
			o.Validate.VersionIndex,
			len(spec.Versions))
	}
	resources, err := render.ComputeResourcesForApplication(o.Ctx, inputValues, spec.Versions[o.Validate.VersionIndex])
	if err != nil {
		return errors.Wrapf(err, "unable to compute resources on version %v", o.Validate.VersionIndex)
	}
	if o.Validate.Verbose {
		for _, r := range resources {
			b := []byte{}
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
