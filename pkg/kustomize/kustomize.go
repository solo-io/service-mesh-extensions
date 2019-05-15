 package kustomize

import (
	"bytes"
	"context"
	"path/filepath"

	"github.com/solo-io/service-mesh-hub/pkg/kustomize/loader"
	"github.com/solo-io/service-mesh-hub/pkg/kustomize/plugins"

	"github.com/pkg/errors"

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

//go:generate mockgen -package=mocks -mock_names Loader=MockPluginLoader -destination=../internal/mocks/loader_factory.go sigs.k8s.io/kustomize/pkg/plugins LoaderFactory

var (
	PluginTypeError = func(plugin plugins.NamedPlugin) error {
		return errors.Errorf("invalid kustomize plugin: plugin %s of type %T must implement "+
			"either the generator or the transformer interfaces", plugin.Name(), plugin)
	}
)

type LayerEngine interface {
	Run(dir string) ([]byte, error)
}

type Kustomizer struct {
	pathLoader   loader.Loader
	ctx          context.Context
	manifests    helmchart.Manifests
	overlay      *hubv1.Kustomize
	pluginLoader kplugins.LoaderFactory
}

func NewKustomizer(loader loader.Loader, manifests helmchart.Manifests, layer *hubv1.Kustomize, kPlugins ...plugins.NamedPlugin) (*Kustomizer, error) {

	var generators []plugins.NamedGenerator
	var transformers []plugins.NamedTransformer

	for _, p := range kPlugins {
		assigned := false

		genPlugin, ok := p.(plugins.NamedGenerator)
		if ok {
			generators = append(generators, genPlugin)
			assigned = true
		}

		trPlugin, ok := p.(plugins.NamedTransformer)
		if ok {
			transformers = append(transformers, trPlugin)
			assigned = true
		}

		if !assigned {
			return nil, PluginTypeError(p)
		}
	}

	return &Kustomizer{
		overlay:      layer,
		pathLoader:   loader,
		manifests:    manifests,
		pluginLoader: plugins.NewStaticPluginLoader(generators, transformers),
	}, nil
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

	err = options.RunBuild(buf, fSys, rf, transformer.NewFactoryImpl(), kplugins.NewLoader(rf, k.pluginLoader))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
