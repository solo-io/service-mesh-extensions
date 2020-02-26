package render

import (
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/installutils/kuberesource"
)

func GetResources(manifests helmchart.Manifests) (kuberesource.UnstructuredResources, error) {
	manifestBytes := []byte(manifests.CombinedString())
	resources, err := YamlToResources(manifestBytes)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
