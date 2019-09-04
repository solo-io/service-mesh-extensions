package render

import (
	"bytes"
	"context"
	"text/template"

	"k8s.io/helm/pkg/manifest"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/go-utils/installutils"
	"github.com/solo-io/go-utils/installutils/helmchart"
	"github.com/solo-io/go-utils/installutils/kuberesource"
	hubv1 "github.com/solo-io/service-mesh-hub/api/v1"
	"go.uber.org/zap"
)

var (
	MissingInstallSpecError = errors.Errorf("missing installation spec")

	FailedToRenderManifestsError = func(err error) error {
		return errors.Wrapf(err, "error rendering manifests")
	}

	FailedToConvertManifestsError = func(err error) error {
		return errors.Wrapf(err, "error converting manifests to raw resources")
	}

	FailedRenderValueTemplatesError = func(err error) error {
		return errors.Wrapf(err, "error rendering input value templates")
	}

	MissingInputForRequiredLayer = func(err error) error {
		return errors.Wrapf(err, "error retrieving input for required layer")
	}

	IncorrectNumberOfInputLayersError = errors.Errorf("incorrect number of input layers")
)

type SuperglooInfo struct {
	Namespace          string
	ServiceAccountName string
	ClusterRoleName    string
}

type LayerInput struct {
	LayerId, OptionId string
}

type ValuesInputs struct {
	Name             string
	InstallNamespace string
	Flavor           *hubv1.Flavor
	Layers           []LayerInput
	MeshRef          core.ResourceRef

	UserDefinedValues string
	SpecDefinedValues string
	// These map to the params found on versions, flavors, and layers,
	Params map[string]string

	Supergloo SuperglooInfo
}

// Deprecated: use ManifestRenderer.ComputeResourcesForApplication
func ComputeResourcesForApplication(ctx context.Context, inputs ValuesInputs, spec *hubv1.VersionedApplicationSpec) (kuberesource.UnstructuredResources, error) {
	renderer := NewManifestRenderer()
	return renderer.ComputeResourcesForApplication(ctx, inputs, spec)
}

func ValidateInputs(inputs ValuesInputs) error {
	if len(inputs.Layers) != GetRequiredLayerCount(inputs.Flavor) {
		return IncorrectNumberOfInputLayersError
	}

	for _, flavorLayer := range inputs.Flavor.CustomizationLayers {
		layer, err := GetLayer(flavorLayer.Id, inputs.Flavor)
		if err != nil && !flavorLayer.Optional {
			return MissingInputForRequiredLayer(err)
		} else if err != nil && flavorLayer.Optional {
			continue
		}

		var optionId string
		for _, layerInput := range inputs.Layers {
			if layerInput.LayerId == flavorLayer.Id {
				optionId = layerInput.OptionId
			}
		}

		_, err = GetLayerOption(optionId, layer)
		if err != nil {
			return err
		}
	}

	// TODO JOEKELLEY make sure optional layers are handled correctly
	// for _, inputLayer := range inputs.Layers {
	//	layer, err := GetLayer(inputLayer.LayerId, inputs.Flavor)
	//	if err != nil {
	//		return err
	//	}
	//
	//	var flavorLayer *hubv1.Layer
	//	for _, l := range inputs.Flavor.CustomizationLayers {
	//		if layer.Id == l.Id {
	//			flavorLayer = l
	//			break
	//		}
	//	}
	//	if flavorLayer == nil {
	//		return UnexpectedInputLayerIdError
	//	}
	//
	//	if inputLayer.OptionId == "" && !flavorLayer.Optional {
	//		return InvalidLayerConfigError
	//	}
	//
	//}

	return nil
}

/*
 Coalesces spec values yaml, layer values, params, and user-defined values yaml.
 User defined values override params which override layer values which override spec values.
 If there is an error parsing, it is logged and propagated.
*/
func ComputeValueOverrides(ctx context.Context, inputs ValuesInputs) (string, error) {
	valuesMap := make(map[string]interface{})

	specValues, err := ConvertYamlStringToNestedMap(inputs.SpecDefinedValues)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw("Error parsing spec values yaml",
			zap.Error(err),
			zap.String("values", inputs.SpecDefinedValues))
		return "", err
	}
	valuesMap = CoalesceValuesMap(ctx, valuesMap, specValues)

	for _, layerInput := range inputs.Layers {
		option, err := GetLayerOptionTwo(layerInput.LayerId, layerInput.OptionId, inputs.Flavor)
		if err != nil {
			return "", err
		}

		if option.HelmValues != "" {
			layerValues, err := ConvertYamlStringToNestedMap(option.HelmValues)
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorw("Error parsing layer values yaml",
					zap.Error(err),
					zap.String("values", option.HelmValues))
				return "", err
			}
			valuesMap = CoalesceValuesMap(ctx, valuesMap, layerValues)
		}
	}

	paramValues, err := ConvertParamsToNestedMap(inputs.Params)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw("Error parsing install params",
			zap.Error(err),
			zap.Any("params", inputs.Params))
		return "", err
	}
	valuesMap = CoalesceValuesMap(ctx, valuesMap, paramValues)

	userValues, err := ConvertYamlStringToNestedMap(inputs.UserDefinedValues)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw("Error parsing user values yaml",
			zap.Error(err),
			zap.Any("params", inputs.UserDefinedValues))
		return "", err
	}
	valuesMap = CoalesceValuesMap(ctx, valuesMap, userValues)

	values, err := ConvertNestedMapToYaml(valuesMap)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw(err.Error(), zap.Error(err), zap.Any("valuesMap", valuesMap))
		return "", err
	}
	return values, nil
}

func GetManifestsFromApplicationSpec(ctx context.Context, inputs ValuesInputs, spec *hubv1.VersionedApplicationSpec) (helmchart.Manifests, error) {
	var manifests helmchart.Manifests
	switch installationSpec := spec.GetInstallationSpec().(type) {
	case *hubv1.VersionedApplicationSpec_GithubChart:
		githubManifests, err := getManifestsFromGithub(ctx, installationSpec.GithubChart, inputs)
		if err != nil {
			return nil, err
		}
		manifests = githubManifests
	case *hubv1.VersionedApplicationSpec_HelmArchive:
		helmManifests, err := getManifestsFromHelm(ctx, installationSpec.HelmArchive, inputs)
		if err != nil {
			return nil, err
		}
		manifests = helmManifests
	case *hubv1.VersionedApplicationSpec_ManifestsArchive:
		archiveManifests, err := getManifestsFromArchive(ctx, installationSpec.ManifestsArchive, inputs)
		if err != nil {
			return nil, err
		}
		manifests = archiveManifests
	case *hubv1.VersionedApplicationSpec_InstallationSteps:
		archiveManifests, err := getManifestsFromSteps(ctx, installationSpec.InstallationSteps, inputs)
		if err != nil {
			return nil, err
		}
		manifests = archiveManifests
	default:
		return nil, MissingInstallSpecError
	}

	return manifests, nil
}

func GetResourcesFromManifests(ctx context.Context, manifests helmchart.Manifests) (kuberesource.UnstructuredResources, error) {
	rawResources, err := manifests.ResourceList()
	if err != nil {
		wrapped := FailedToConvertManifestsError(err)
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(),
			zap.Error(err))
		return nil, wrapped
	}
	return rawResources, nil
}

func FilterByLabel(ctx context.Context, spec *hubv1.VersionedApplicationSpec, resources kuberesource.UnstructuredResources) kuberesource.UnstructuredResources {
	labels := spec.GetRequiredLabels()
	if len(labels) > 0 {
		contextutils.LoggerFrom(ctx).Infow("Filtering installed resources by label", zap.Any("labels", labels))
		return resources.WithLabels(labels)
	}
	return resources
}

func getManifestsFromHelm(ctx context.Context, helmInstallSpec *hubv1.TgzLocation, inputs ValuesInputs) (helmchart.Manifests, error) {
	values, err := ComputeValueOverrides(ctx, inputs)
	if err != nil {
		return nil, err
	}
	contextutils.LoggerFrom(ctx).Infow("Rendering with values", zap.String("values", values))
	manifests, err := helmchart.RenderManifests(ctx,
		helmInstallSpec.Uri,
		values,
		inputs.Name,
		inputs.InstallNamespace,
		"")
	if err != nil {
		wrapped := FailedToRenderManifestsError(err)
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(),
			zap.Error(err),
			zap.String("chartUri", helmInstallSpec.Uri),
			zap.String("values", values),
			zap.String("releaseName", inputs.Name),
			zap.String("namespace", inputs.InstallNamespace),
			zap.String("kubeVersion", ""))
		return nil, wrapped
	}
	return manifests, nil
}

func getManifestsFromGithub(ctx context.Context, githubInstallSpec *hubv1.GithubRepositoryLocation, inputs ValuesInputs) (helmchart.Manifests, error) {
	ref := helmchart.GithubChartRef{
		Owner:          githubInstallSpec.Org,
		Repo:           githubInstallSpec.Repo,
		Ref:            githubInstallSpec.Ref,
		ChartDirectory: githubInstallSpec.Directory,
	}
	values, err := ComputeValueOverrides(ctx, inputs)
	if err != nil {
		return nil, err
	}
	manifests, err := helmchart.RenderManifestsFromGithub(ctx, ref,
		values,
		inputs.Name,
		inputs.InstallNamespace,
		"")
	if err != nil {
		wrapped := FailedToRenderManifestsError(err)
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(),
			zap.Error(err),
			zap.Any("ref", ref),
			zap.String("values", values),
			zap.String("releaseName", inputs.Name),
			zap.String("namespace", inputs.InstallNamespace),
			zap.String("kubeVersion", ""))
		return nil, wrapped
	}
	return manifests, nil
}

func getManifestsFromArchive(ctx context.Context, manifestsArchive *hubv1.TgzLocation, inputs ValuesInputs) (helmchart.Manifests, error) {
	manifests, err := installutils.GetManifestsFromRemoteTar(manifestsArchive.GetUri())
	if err != nil {
		wrapped := FailedToRenderManifestsError(err)
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(),
			zap.Error(err),
			zap.String("manifestsArchiveUrl", manifestsArchive.GetUri()),
			zap.String("releaseName", inputs.Name),
			zap.String("namespace", inputs.InstallNamespace))
		return nil, wrapped
	}
	return manifests, nil
}

const InstallationStepLabel = "service-mesh-hub.solo.io/installation_step"

func getManifestsFromSteps(ctx context.Context, steps *hubv1.InstallationSteps, inputs ValuesInputs) (helmchart.Manifests, error) {
	if len(steps.Steps) == 0 {
		return nil, errors.Errorf("must provide at least one installation step")
	}
	var combinedManifests []manifest.Manifest
	var uniqueStepNames []string
	for _, step := range steps.Steps {
		if step.Name == "" {
			return nil, errors.Errorf("step must be named")
		}
		for _, name := range uniqueStepNames {
			if step.Name == name {
				return nil, errors.Errorf("step names must be unique; %v duplicated", name)
			}
		}
		uniqueStepNames = append(uniqueStepNames, step.Name)

		manifests, err := getManifestsFromInstallationStep(ctx, inputs, step)
		if err != nil {
			return nil, err
		}
		// add labels for step to every resource in the manifests
		resources, err := manifests.ResourceList()
		if err != nil {
			return nil, err
		}
		for _, resource := range resources {
			labels := resource.GetLabels()
			if labels == nil {
				labels = make(map[string]string)
			}
			labels[InstallationStepLabel] = step.Name
			resource.SetLabels(labels)
		}

		manifests, err = helmchart.ManifestsFromResources(resources)
		if err != nil {
			return nil, err
		}

		combinedManifests = append(combinedManifests, manifests...)
	}
	return combinedManifests, nil
}

func getManifestsFromInstallationStep(ctx context.Context, inputs ValuesInputs, step *hubv1.InstallationSteps_Step) (helmchart.Manifests, error) {
	var manifests helmchart.Manifests
	switch installationSpec := step.Step.(type) {
	case *hubv1.InstallationSteps_Step_GithubChart:
		githubManifests, err := getManifestsFromGithub(ctx, installationSpec.GithubChart, inputs)
		if err != nil {
			return nil, err
		}
		manifests = githubManifests
	case *hubv1.InstallationSteps_Step_HelmArchive:
		helmManifests, err := getManifestsFromHelm(ctx, installationSpec.HelmArchive, inputs)
		if err != nil {
			return nil, err
		}
		manifests = helmManifests
	case *hubv1.InstallationSteps_Step_ManifestsArchive:
		archiveManifests, err := getManifestsFromArchive(ctx, installationSpec.ManifestsArchive, inputs)
		if err != nil {
			return nil, err
		}
		manifests = archiveManifests
	default:
		return nil, MissingInstallSpecError
	}

	return manifests, nil
}

// The SpecDefinedValues, UserDefinedValues, and Params inputs can contain template
// actions (text delimited by "{{" and "}}" ). This function renders the contents of these
// parameters using the data contained in 'input' and updates 'input' with the results.
func ExecInputValuesTemplates(inputs ValuesInputs) (ValuesInputs, error) {

	// Render the helm values string that comes from the extension spec
	buf := new(bytes.Buffer)
	tpl := template.Must(template.New("specValues").Parse(inputs.SpecDefinedValues))
	if err := tpl.Execute(buf, inputs); err != nil {
		return ValuesInputs{}, err
	}
	inputs.SpecDefinedValues = buf.String()
	buf.Reset()

	// Render the helm values string that comes from the user provided overrides
	tpl = template.Must(template.New("userValues").Parse(inputs.UserDefinedValues))
	if err := tpl.Execute(buf, inputs); err != nil {
		return ValuesInputs{}, err
	}
	inputs.UserDefinedValues = buf.String()
	buf.Reset()

	// Render the values of the parameters
	for paramName, paramValue := range inputs.Params {
		t := template.Must(template.New(paramName).Parse(paramValue))
		if err := t.Execute(buf, inputs); err != nil {
			return ValuesInputs{}, err
		}
		inputs.Params[paramName] = buf.String()
		buf.Reset()
	}

	return inputs, nil
}
