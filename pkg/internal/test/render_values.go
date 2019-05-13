package test

import (
	"github.com/solo-io/service-mesh-hub/api/v1"
)

type RenderValuesForTest struct {
	MeshRef               v1.ResourceRef
	SuperglooNamespace    string
	InstallationNamespace string
	Custom                interface{}
}

func GetRenderValues() RenderValuesForTest {
	return RenderValuesForTest{
		MeshRef: v1.ResourceRef{
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
