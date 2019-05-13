package render

import (
	"regexp"

	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	yaml2json "k8s.io/apimachinery/pkg/util/yaml"
)

var (
	yamlSeparator = regexp.MustCompile("\n---")
)

func YamlToResources(yamlBytes []byte) (kuberesource.UnstructuredResources, error) {
	snippets := yamlSeparator.Split(string(yamlBytes), -1)
	var resources kuberesource.UnstructuredResources
	for _, objectYaml := range snippets {
		if helmchart.IsEmptyManifest(objectYaml) {
			continue
		}
		jsn, err := yaml2json.ToJSON([]byte(objectYaml))
		if err != nil {
			return nil, err
		}

		uncastObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jsn)
		if err != nil {
			return nil, err
		}
		if resourceList, ok := uncastObj.(*unstructured.UnstructuredList); ok {
			for _, item := range resourceList.Items {
				resources = append(resources, &item)
			}
			continue
		}
		resources = append(resources, uncastObj.(*unstructured.Unstructured))
	}
	return resources, nil
}
