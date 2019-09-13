package test

import (
	"github.com/solo-io/service-mesh-hub/pkg/render/inputs"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

func GetRenderValues() inputs.ManifestRenderValues {
	return inputs.ManifestRenderValues{
		MeshRef: core.ResourceRef{
			Name:      "mesh-name",
			Namespace: "mesh-ns",
		},
		InstallNamespace: "install-ns",
		Custom: map[string]interface{}{
			"SomeValue":          "this-is-a-custom-value",
			"SuperglooNamespace": "supergloo-system",
		},
	}
}
