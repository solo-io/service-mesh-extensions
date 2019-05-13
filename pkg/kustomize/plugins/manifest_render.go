package plugins

import (
	"bytes"
	"encoding/json"
	"text/template"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"

	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
)

const (
	ManifestRenderPluginName = "ManifestRender"
)

var (
	EmptyManifestError = errors.Errorf("manifest in ManifestRender cannot be empty")
)

type ManifestRenderValues struct {
	MeshRef               v1.ResourceRef
	SuperglooNamespace    string
	InstallationNamespace string
}

type manifestRenderPlugin struct {
	values   ManifestRenderValues
	Manifest string
}

func NewManifestRenderPlugin() *manifestRenderPlugin {
	values := ManifestRenderValues{
		MeshRef: v1.ResourceRef{
			Name:      "", // state.Mesh.Name,
			Namespace: "", // state.Mesh.Namespace,
		},
		SuperglooNamespace:    "supergloo-system",
		InstallationNamespace: "", //state.InstallNamespace,
	}
	return &manifestRenderPlugin{values: values}
}

func (p *manifestRenderPlugin) Name() string {
	return ManifestRenderPluginName
}

func (p *manifestRenderPlugin) Config(ldr ifc.Loader, rf *resmap.Factory, k ifc.Kunstructured) error {
	byt, err := k.MarshalJSON()
	if err != nil {
		return err
	}

	var plugin manifestRenderPlugin
	err = json.Unmarshal(byt, &plugin)
	if err != nil {
		return err
	}
	if plugin.Manifest == "" {
		return EmptyManifestError
	}
	p.Manifest = plugin.Manifest

	return nil
}

func (p *manifestRenderPlugin) Generate() (resmap.ResMap, error) {
	var buf bytes.Buffer

	temp := template.Must(template.New("manifest").Parse(p.Manifest))
	err := temp.Execute(&buf, p.values)
	if err != nil {
		return nil, err
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()))
	return rf.NewResMapFromBytes(buf.Bytes())
}
