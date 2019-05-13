package kustomize

import (
	"bytes"
	"context"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"
	"path/filepath"

	"github.com/solo-io/go-utils/installutils/helmchart"
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/k8sdeps/kv/plugin"
	"sigs.k8s.io/kustomize/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/pkg/commands/build"
	"sigs.k8s.io/kustomize/pkg/fs"
	kplugins "sigs.k8s.io/kustomize/pkg/plugins"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
)

//go:generate mockgen -package=mocks -mock_names Loader=MockPluginLoader -destination=../../pkg/internal/mocks/loader_factory.go sigs.k8s.io/kustomize/pkg/plugins LoaderFactory

type LayerEngine interface {
	Run(dir string) ([]byte, error)
}

type Kustomizer struct {
	pathLoader loader.Loader
	ctx        context.Context
	manifests  helmchart.Manifests

	//installState *InstallationState
	overlay      *hubv1.Kustomize
}

func NewKustomizer(loader loader.Loader, manifests helmchart.Manifests, layer *hubv1.Kustomize) *Kustomizer {
	return &Kustomizer{overlay: layer, pathLoader: loader, manifests: manifests}
}

func (k *Kustomizer) Run(dir string) ([]byte, error) {

	newDir, err := k.pathLoader.RetrieveLayers(dir, k.overlay)
	if err != nil {
		return nil, err
	}

	err = k.pathLoader.LoadBase(k.manifests, newDir)
	if err != nil {
		return nil, err
	}

	fSys := fs.MakeRealFS()
	buf := &bytes.Buffer{}

	options := build.NewOptions(filepath.Join(newDir, k.overlay.OverlayPath), "")

	pluginConfig := plugin.ActivePluginConfig()

	// Configuration for ConfigMap and Secret generators.
	genMetaArgs := types.GeneratorMetaArgs{
		PluginConfig: pluginConfig,
	}
	uf := kunstruct.NewKunstructuredFactoryWithGeneratorArgs(&genMetaArgs)
	rf := resmap.NewFactory(resource.NewFactory(uf))

	err = options.RunBuild(
		buf, fSys,
		rf,
		transformer.NewFactoryImpl(),
		kplugins.NewLoader(rf, plugins.NewStaticPluginLoader(
			[]plugins.NamedGenerator{plugins.NewManifestRenderPlugin(k.installState)},
			nil,
		)),
	)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
