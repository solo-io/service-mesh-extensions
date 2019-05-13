package test

import (
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/k8sdeps/kv/plugin"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
)

func ResourceMapFactory() *resmap.Factory {
	pluginConfig := plugin.ActivePluginConfig()
	genMetaArgs := types.GeneratorMetaArgs{
		PluginConfig: pluginConfig,
	}
	uf := kunstruct.NewKunstructuredFactoryWithGeneratorArgs(&genMetaArgs)
	rf := resmap.NewFactory(resource.NewFactory(uf))
	return rf
}
