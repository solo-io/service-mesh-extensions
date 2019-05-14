package test

import "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

// TODO: get rid of this and just use the ValuesInputs type
type RenderValuesForTest struct {
	MeshRef            core.ResourceRef
	SuperglooNamespace string
	InstallNamespace   string
	Custom             interface{}
}

func GetRenderValues() RenderValuesForTest {
	return RenderValuesForTest{
		MeshRef: core.ResourceRef{
			Name:      "mesh-name",
			Namespace: "mesh-ns",
		},
		SuperglooNamespace: "supergloo-system",
		InstallNamespace:   "install-ns",
		Custom: map[string]interface{}{
			"SomeValue": "this-is-a-custom-value",
		},
	}
}
