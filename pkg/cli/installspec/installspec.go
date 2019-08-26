package installspec

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/solo-io/go-utils/protoutils"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/service-mesh-hub/pkg/util"
)

type InstallSpec struct {
	Values  render.ValuesInputs
	Version *v1.VersionedApplicationSpec
}

// Workaround for being unable to marshal/unmarshal oneofs on proto messages nested in standard structs.
// Contains human-readable yaml strings for each field in the values and versioned spec structs.
type installSpecYaml struct {
	Values  string
	Version string
}

func (i *InstallSpec) Load(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	spec := &installSpecYaml{}
	if err := yaml.Unmarshal(bytes, spec); err != nil {
		return err
	}

	values := &render.ValuesInputs{}
	if err := yaml.Unmarshal([]byte(spec.Values), values); err != nil {
		return err
	}

	versionJson, err := yaml.YAMLToJSON([]byte(spec.Version))
	if err != nil {
		return err
	}
	version := &v1.VersionedApplicationSpec{}
	if err := protoutils.UnmarshalBytes(versionJson, version); err != nil {
		return err
	}

	i.Values = *values
	i.Version = version

	return nil
}

func (i *InstallSpec) Save(filename string) error {
	versionBytes, err := protoutils.MarshalBytes(i.Version)
	if err != nil {
		return err
	}

	versionBytes, err = yaml.JSONToYAML(versionBytes)
	if err != nil {
		return err
	}

	valuesBytes, err := yaml.Marshal(i.Values)
	if err != nil {
		return err
	}

	persistedSpec := &installSpecYaml{}
	persistedSpec.Version = string(versionBytes)
	persistedSpec.Values = string(valuesBytes)

	bytes, err := yaml.Marshal(persistedSpec)
	if err != nil {
		return err
	}

	return util.SaveFile(filename, string(bytes))
}
