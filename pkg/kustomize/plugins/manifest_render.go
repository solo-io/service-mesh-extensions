package plugins

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/pkg/errors"
)

const (
	ManifestRenderPluginName = "ManifestRender"
)

var (
	EmptyManifestError = errors.Errorf("manifest in ManifestRender cannot be empty")
)

type manifestRenderPlugin struct {
	values   interface{}
	Manifest string
}

func NewManifestRenderPlugin(values interface{}) *manifestRenderPlugin {
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
