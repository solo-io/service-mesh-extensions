package test

import (
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type TestChart struct {
	Name           string
	Version        string
	HelmUri        string
	DeploymentName string
}

const (
	DefaultFlavorName = "vanilla"
)

var (
	MeshRef = &core.ResourceRef{
		Name:      "name",
		Namespace: "namespace",
	}

	HelloWorldChart1_0 = TestChart{
		Name:           "hello-world",
		Version:        "0.1.0",
		HelmUri:        "https://storage.googleapis.com/solo-helm-charts/helloworld-chart-0.1.0.tgz",
		DeploymentName: "hello-world-helloworld-chart",
	}
	HelloWorldChart1_1 = TestChart{
		Name:           "hello-world",
		Version:        "0.1.1",
		HelmUri:        "https://storage.googleapis.com/solo-helm-charts/helloworld-chart-0.1.1.tgz",
		DeploymentName: "hello-world-helloworld-chart",
	}

	kustomizeLocation = &hubv1.Kustomize_TgzArchive{
		TgzArchive: &hubv1.TgzLocation{
			Uri: "https://storage.googleapis.com/solo-kustomize-plugins/test.tgz",
		},
	}

	KustomizeLayer_Success = &hubv1.Kustomize{
		OverlayPath: "supergloo",
		Location:    kustomizeLocation,
	}

	KustomizeLayer_Failure = &hubv1.Kustomize{
		OverlayPath: "error",
		Location:    kustomizeLocation,
	}

	KustomizeLayer_NotFound = &hubv1.Kustomize{
		OverlayPath: "fails",
		Location:    kustomizeLocation,
	}
)

func GetAppSpec(chart TestChart, kustomizeLayer *hubv1.Kustomize) *hubv1.VersionedApplicationSpec {
	return &hubv1.VersionedApplicationSpec{
		InstallationSpec: &hubv1.VersionedApplicationSpec_HelmArchive{
			HelmArchive: &hubv1.TgzLocation{
				Uri: chart.HelmUri,
			},
		},
		Version: chart.Version,
		Flavors: []*hubv1.Flavor{GetFlavor(kustomizeLayer)},
	}
}

func GetFlavor(kustomizeLayer *hubv1.Kustomize) *hubv1.Flavor {
	flavor := &hubv1.Flavor{
		Name: DefaultFlavorName,
	}
	if kustomizeLayer != nil {
		flavor.CustomizationLayers = []*hubv1.Layer{
			{
				Type: &hubv1.Layer_Kustomize{
					Kustomize: kustomizeLayer,
				},
			},
		}
	}

	return flavor
}
