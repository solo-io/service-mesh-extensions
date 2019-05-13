package test

import "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

type RenderValuesForTest struct {
	MeshRef               core.ResourceRef
	SuperglooNamespace    string
	InstallationNamespace string
	Custom                interface{}
}

func GetRenderValues() RenderValuesForTest {
	return RenderValuesForTest{
		MeshRef: core.ResourceRef{
			Name:      "mesh-name",
			Namespace: "mesh-ns",
		},
		SuperglooNamespace:    "supergloo-system",
		InstallationNamespace: "install-ns",
		Custom: map[string]interface{}{
			"SomeValue": "this-is-a-custom-value",
		},
	}
}
