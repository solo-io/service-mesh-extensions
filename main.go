package main

import (
	"context"
	"fmt"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/protoutils"
	"github.com/solo-io/service-mesh-hub/api/v1"
	"github.com/solo-io/service-mesh-hub/pkg/render"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"
	"io/ioutil"
)

func main() {
	ctx := context.Background()
	err := Run(ctx)
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Fatal error during run", zap.Error(err))
	}
}

func Run(ctx context.Context) error {
	inputValues := render.ValuesInputs{
		Name:               "glooshot",
		InstallNamespace:   "default",
		FlavorName:         "istio",
		MeshRef:            core.ResourceRef{},
		SuperglooNamespace: "sm-marketplace",
		UserDefinedValues:  "",
		FlavorParams:       nil,
		SpecDefinedValues:  "",
	}
	spec, err := LoadExtensionSpec("./extensions/v1/glooshot/spec.yaml")
	if err != nil {
		return err
	}

	resources, err := render.ComputeResourcesForApplication(ctx, inputValues, spec.Versions[0])
	if err != nil {
		return err
	}

	fmt.Println(resources)

	return nil
}

func LoadExtensionSpec(pathToSpec string) (*v1.ApplicationSpec, error) {
	bytes, err := ioutil.ReadFile(pathToSpec)
	if err != nil {
		return nil, err
	}
	var spec v1.ApplicationSpec
	err = protoutils.UnmarshalYaml(bytes, &spec)
	if err != nil {
		return nil, err
	}
	return &spec, nil
}
