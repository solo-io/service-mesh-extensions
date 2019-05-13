package plugins

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/transformers"
)

//go:generate mockgen -package=mocks -destination=../../internal/mocks/generators.go github.com/solo-io/service-mesh-hub/pkg/kustomize/plugins NamedGenerator,NamedTransformer

var (
	PluginNotLoadedError = func(name string) error {
		return errors.Errorf("plugin %s was never loaded", name)
	}
)

type NamedPlugin interface {
	Name() string
}

type NamedGenerator interface {
	NamedPlugin
	transformers.Generator
}

type NamedTransformer interface {
	NamedPlugin
	transformers.Transformer
}

type StaticGeneratorMap map[string]transformers.Generator
type StaticTransformerMap map[string]transformers.Transformer

type staticPluginLoader struct {
	generatorMap   StaticGeneratorMap
	transformerMap StaticTransformerMap
}

func NewStaticPluginLoader(generators []NamedGenerator, tforms []NamedTransformer) *staticPluginLoader {
	generatorMap := make(StaticGeneratorMap)
	for _, v := range generators {
		generatorMap[v.Name()] = v
	}

	transformerMap := make(StaticTransformerMap)
	for _, v := range tforms {
		transformerMap[v.Name()] = v
	}
	return &staticPluginLoader{generatorMap: generatorMap, transformerMap: transformerMap}
}

func (pl *staticPluginLoader) LoadGenerator(ldr ifc.Loader, res *resource.Resource) (transformers.Generator, error) {
	for k, v := range pl.generatorMap {
		if k == res.GetKind() {
			return v, nil
		}
	}
	return nil, PluginNotLoadedError(res.GetKind())
}

func (pl *staticPluginLoader) LoadTransformer(ldr ifc.Loader, res *resource.Resource) (transformers.Transformer, error) {
	for k, v := range pl.transformerMap {
		if k == res.GetKind() {
			return v, nil
		}
	}
	return nil, PluginNotLoadedError(res.GetKind())
}
