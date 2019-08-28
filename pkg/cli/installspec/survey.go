package installspec

import (
	"fmt"
	"strings"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/registry"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"gopkg.in/AlecAivazis/survey.v1"
)

func GetInstallSpec(reader registry.SpecReader, installNamespace string) (*InstallSpec, error) {
	spec, err := selectApplication(reader)
	if err != nil {
		return nil, err
	}
	version, err := selectVersion(spec)
	if err != nil {
		return nil, err
	}
	values, err := GetValuesInputs(spec, version, installNamespace)
	if err != nil {
		return nil, err
	}

	return &InstallSpec{
		Values:  *values,
		Version: version,
	}, nil
}

func GetValuesInputs(spec *v1.ApplicationSpec, version *v1.VersionedApplicationSpec, installNamespace string) (*render.ValuesInputs, error) {
	values := render.ValuesInputs{
		InstallNamespace: installNamespace,
	}

	values.Name = spec.Name
	values.SpecDefinedValues = version.ValuesYaml
	flavor, err := selectFlavor(version)
	if err != nil {
		return nil, err
	}
	values.FlavorName = flavor.Name
	values.Params = make(map[string]string)
	for _, layer := range flavor.GetCustomizationLayers() {
		for _, param := range layer.GetParameters() {
			val, err := selectParam(param)
			if err != nil {
				return nil, err
			}
			values.Params[param.Name] = val
		}
	}
	return &values, nil
}

func selectApplication(reader registry.SpecReader) (*v1.ApplicationSpec, error) {
	specs, err := reader.GetSpecs()
	if err != nil {
		return nil, err
	}
	var names []string
	nameToSpec := make(map[string]*v1.ApplicationSpec)
	for _, spec := range specs {
		nameToSpec[spec.Name] = spec
		names = append(names, spec.Name)
	}
	specName := ""
	prompt := &survey.Select{
		Options:  names,
		Message:  "What application to install?",
		PageSize: 10,
	}
	err = survey.AskOne(prompt, &specName, survey.Required)
	if err != nil {
		return nil, err
	}
	return nameToSpec[specName], nil
}

func selectVersion(spec *v1.ApplicationSpec) (*v1.VersionedApplicationSpec, error) {
	var versions []string
	versionToSpec := make(map[string]*v1.VersionedApplicationSpec)
	for _, versionedSpec := range spec.Versions {
		versionToSpec[versionedSpec.Version] = versionedSpec
		versions = append(versions, versionedSpec.Version)
	}
	specVersion := ""
	prompt := &survey.Select{
		Options:  versions,
		Message:  "What version to install?",
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &specVersion, survey.Required)
	if err != nil {
		return nil, err
	}
	return versionToSpec[specVersion], nil
}

func selectFlavor(spec *v1.VersionedApplicationSpec) (*v1.Flavor, error) {
	var flavors []string
	nameToFlavor := make(map[string]*v1.Flavor)
	for _, flavor := range spec.GetFlavors() {
		if strings.Contains(flavor.Name, "supergloo") {
			// These flavors require supergloo with a cluster-admin role.
			continue
		}
		nameToFlavor[flavor.Name] = flavor
		flavors = append(flavors, flavor.Name)
	}
	flavor := ""
	prompt := &survey.Select{
		Options:  flavors,
		Message:  "What flavor to install?",
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &flavor, survey.Required)
	if err != nil {
		return nil, err
	}
	return nameToFlavor[flavor], nil
}

func selectParam(spec *v1.Parameter) (string, error) {
	prompt := &survey.Input{
		Default: spec.Default,
		Message: fmt.Sprintf("[%s] %s", spec.Description, spec.Name),
	}
	input := ""
	err := survey.AskOne(prompt, &input, nil)
	return input, err
}
