package inputs

import "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

type ManifestRenderValues struct {
	Name             string
	InstallNamespace string
	MeshRef          core.ResourceRef

	// Custom values come from the parameters set on a flavor
	Custom interface{}
}
