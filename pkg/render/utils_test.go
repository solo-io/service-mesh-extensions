package render_test

import (
	"context"

	"github.com/solo-io/service-mesh-hub/pkg/render/validation"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/service-mesh-hub/pkg/render"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

var _ = Describe("utils", func() {
	Context("flavor selection", func() {
		It("returns an error when flavor name is empty", func() {
			_, err := render.GetInstalledFlavor("", nil)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(render.NilFlavorError))
		})
		It("returns an error when no flavors are found", func() {
			_, err := render.GetInstalledFlavor("fsdf", nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(render.ExpectedAtMostError("flavor", 1, 0).Error()))
		})
		It("returns an error when the flavor isn't found", func() {
			_, err := render.GetInstalledFlavor("flavor2", []*v1.Flavor{{Name: "flavor1"}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(render.NoFlavorFoundError("flavor2").Error()))
		})
		It("Succeeds when the relevant flavor is found", func() {
			flavor := &v1.Flavor{Name: "flavor1"}
			result, err := render.GetInstalledFlavor("flavor1", []*v1.Flavor{flavor})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(flavor))
		})
	})

	Context("params", func() {
		It("works", func() {
			input := make(map[string]string)
			input["a.b.c"] = "foo"
			input["a.b.d"] = "bar"
			input["d"] = "baz"
			input["b.c"] = "goo"
			expected := map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "foo",
						"d": "bar",
					},
				},
				"d": "baz",
				"b": map[string]interface{}{
					"c": "goo",
				},
			}
			Expect(render.ConvertParamsToNestedMap(input)).To(Equal(expected))
		})

		It("errors on invalid value", func() {
			input := make(map[string]string)
			input["a"] = "{"
			expectedErr := render.UnableToParseParameterError(errors.Errorf(""), "a", "{")
			out, actualErr := render.ConvertParamsToNestedMap(input)
			Expect(out).To(BeNil())
			Expect(actualErr.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Context("yaml to map", func() {
		It("works", func() {
			yamlString := "foo:\n  bar: baz"
			expected := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			}
			actual, err := render.ConvertYamlStringToNestedMap(yamlString)
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

		It("errors on invalid yaml", func() {
			brokenString := "foo\n:bar"
			actual, err := render.ConvertYamlStringToNestedMap(brokenString)
			Expect(actual).To(BeNil())
			Expect(err.Error()).To(ContainSubstring(render.UnableToParseYamlError(errors.Errorf(""), brokenString).Error()))
		})

		It("works for empty string", func() {
			actual, err := render.ConvertYamlStringToNestedMap("")
			Expect(err).To(BeNil())
			Expect(actual).NotTo(BeNil())
		})
	})

	Context("map to yaml", func() {
		It("works", func() {
			expected := "foo:\n  bar: baz\n"
			nestedMap := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			}
			actual, err := render.ConvertNestedMapToYaml(nestedMap)
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expected))
		})

		It("works for nil", func() {
			actual, err := render.ConvertNestedMapToYaml(nil)
			Expect(err).To(BeNil())
			Expect(actual).To(BeEquivalentTo(""))
		})
	})

	Context("coalesce values", func() {
		It("works for disjoint values", func() {
			initial := map[string]interface{}{
				"foo": "bar",
			}
			expected := map[string]interface{}{
				"foo": "bar",
				"baz1": map[string]interface{}{
					"baz2": "baz3",
				},
			}
			overrides := map[string]interface{}{
				"baz1": map[string]interface{}{
					"baz2": "baz3",
				},
			}
			actual := render.CoalesceValuesMap(context.TODO(), initial, overrides)
			Expect(actual).To(Equal(expected))
		})

		It("works for overriding previous values", func() {
			initial := map[string]interface{}{
				"foo": "bar",
			}
			expected := map[string]interface{}{
				"foo": "baz",
			}
			overrides := map[string]interface{}{
				"foo": "baz",
			}
			actual := render.CoalesceValuesMap(context.TODO(), initial, overrides)
			Expect(actual).To(Equal(expected))
		})

		It("allows overriding value with map", func() {
			initial := map[string]interface{}{
				"foo": "bar",
			}
			expected := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			}
			overrides := expected
			actual := render.CoalesceValuesMap(context.TODO(), initial, overrides)
			Expect(actual).To(Equal(expected))
		})

		It("allows overriding map with value", func() {
			initial := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			}
			expected := map[string]interface{}{
				"foo": "bar",
			}
			overrides := expected
			actual := render.CoalesceValuesMap(context.TODO(), initial, overrides)
			Expect(actual).To(Equal(expected))
		})

		It("works for empty edge cases", func() {
			emptyMap := make(map[string]interface{})
			nonemptyMap := map[string]interface{}{
				"foo": "bar",
			}
			actual := render.CoalesceValuesMap(context.TODO(), emptyMap, emptyMap)
			Expect(actual).To(Equal(emptyMap))
			actual = render.CoalesceValuesMap(context.TODO(), emptyMap, nonemptyMap)
			Expect(actual).To(Equal(nonemptyMap))
			actual = render.CoalesceValuesMap(context.TODO(), nonemptyMap, emptyMap)
			Expect(actual).To(Equal(nonemptyMap))
		})

		It("works for nil edge cases", func() {
			emptyMap := make(map[string]interface{})
			nonemptyMap := map[string]interface{}{
				"foo": "bar",
			}
			actual := render.CoalesceValuesMap(context.TODO(), nil, nil)
			Expect(actual).To(Equal(emptyMap))
			actual = render.CoalesceValuesMap(context.TODO(), nil, nonemptyMap)
			Expect(actual).To(Equal(nonemptyMap))
			actual = render.CoalesceValuesMap(context.TODO(), nonemptyMap, nil)
			Expect(actual).To(Equal(nonemptyMap))
		})
	})

	Context("compute value overrides", func() {
		It("works", func() {
			inputs := render.ValuesInputs{
				UserDefinedValues: "foo: bar",
				Params: map[string]string{
					"baz1.baz2": "baz3",
				},
				SpecDefinedValues: "goo: hoo",
			}
			expected := "baz1:\n  baz2: baz3\nfoo: bar\ngoo: hoo\n"
			Expect(render.ComputeValueOverrides(context.TODO(), inputs)).To(BeEquivalentTo(expected))
		})

		It("prefers flavor params to spec values", func() {
			inputs := render.ValuesInputs{
				Params: map[string]string{
					"foo": "bar",
				},
				SpecDefinedValues: "foo: baz",
			}
			expected := "foo: bar\n"
			Expect(render.ComputeValueOverrides(context.TODO(), inputs)).To(BeEquivalentTo(expected))
		})

		It("prefers user params to flavor params", func() {
			inputs := render.ValuesInputs{
				UserDefinedValues: "foo: bar",
				Params: map[string]string{
					"foo": "baz",
				},
			}
			expected := "foo: bar\n"
			Expect(render.ComputeValueOverrides(context.TODO(), inputs)).To(BeEquivalentTo(expected))
		})

		It("handles empty case", func() {
			inputs := render.ValuesInputs{}
			expected := ""
			Expect(render.ComputeValueOverrides(context.TODO(), inputs)).To(BeEquivalentTo(expected))
		})

		It("errors on invalid user values", func() {
			str := "invalidYaml"
			inputs := render.ValuesInputs{
				UserDefinedValues: str,
				Params: map[string]string{
					"baz1.baz2": "baz3",
				},
				SpecDefinedValues: "goo: hoo",
			}
			_, err := render.ComputeValueOverrides(context.TODO(), inputs)
			Expect(err.Error()).To(ContainSubstring(render.UnableToParseYamlError(errors.Errorf(""), str).Error()))
		})

		It("errors on invalid spec values", func() {
			str := "invalidYaml"
			inputs := render.ValuesInputs{
				SpecDefinedValues: str,
			}
			_, err := render.ComputeValueOverrides(context.TODO(), inputs)
			Expect(err.Error()).To(ContainSubstring(render.UnableToParseYamlError(errors.Errorf(""), str).Error()))
		})

		It("errors on invalid param values", func() {
			key := "invalid"
			invalid := "{{"
			inputs := render.ValuesInputs{
				Params: map[string]string{
					key: invalid,
				},
			}
			_, err := render.ComputeValueOverrides(context.TODO(), inputs)
			Expect(err.Error()).To(ContainSubstring(render.UnableToParseParameterError(errors.Errorf(""), key, invalid).Error()))
		})
	})

	Context("validate inputs", func() {
		It("works in an empty case", func() {
			inputs := render.ValuesInputs{Flavor: &v1.Flavor{}}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).NotTo(HaveOccurred())
		})

		It("works in a non-empty case", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{
					CustomizationLayers: []*v1.Layer{
						{
							Id:      "a",
							Options: []*v1.LayerOption{{Id: "1"}},
						},
					},
					Parameters: []*v1.Parameter{{
						Name:     "foo",
						Required: true,
					}},
				},
				Layers: []render.LayerInput{{LayerId: "a", OptionId: "1"}},
				Params: map[string]string{"foo": "bar", "bar": "baz"},
			}
			version := v1.VersionedApplicationSpec{
				Parameters: []*v1.Parameter{
					{
						Name:     "bar",
						Required: true,
					},
					{
						Name: "optional",
					},
				},
			}
			err := render.ValidateInputs(inputs, version, validation.NoopValidateResources)
			Expect(err).NotTo(HaveOccurred())
		})

		It("works when no layer or param inputs are provided and no layers or params are required", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{
					CustomizationLayers: []*v1.Layer{
						{
							Id:       "a",
							Optional: true,
							Options:  []*v1.LayerOption{{Id: "1"}},
						},
					},
					Parameters: []*v1.Parameter{{Name: "foo"}},
				},
			}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).NotTo(HaveOccurred())
		})

		It("errors if there is a different number of input layers than required layers", func() {
			inputs := render.ValuesInputs{Flavor: &v1.Flavor{CustomizationLayers: []*v1.Layer{{}}}}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(render.IncorrectNumberOfInputLayersError))
		})

		It("errors if a required layer is missing an option", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{
					CustomizationLayers: []*v1.Layer{
						{
							Id: "a",
							Options: []*v1.LayerOption{
								{
									Id: "1",
								},
							},
						},
						{
							Id: "b",
							Options: []*v1.LayerOption{
								{
									Id: "1",
								},
							},
						},
					},
				},
				Layers: []render.LayerInput{{LayerId: "a", OptionId: "1"}, {LayerId: "z", OptionId: "1"}},
			}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(render.MissingInputForRequiredLayer(errors.New("")).Error()))
		})

		It("errors if a required parameter is missing a value", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{
					Parameters: []*v1.Parameter{{Name: "required", Required: true}},
				},
				Params: map[string]string{"required": ""},
			}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(render.MissingInputForRequireParam("required").Error()))
		})

		It("errors if an unrecognized param is provided", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{Parameters: []*v1.Parameter{{Name: "recognized"}}},
				Params: map[string]string{"unrecognized": "param"},
			}
			err := render.ValidateInputs(inputs, v1.VersionedApplicationSpec{}, validation.NoopValidateResources)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(render.UnrecognizedParamError("unrecognized").Error()))
		})

		It("errors if the environment validation function errors", func() {
			inputs := render.ValuesInputs{
				Flavor: &v1.Flavor{
					CustomizationLayers: []*v1.Layer{
						{
							Id:      "a",
							Options: []*v1.LayerOption{{Id: "1"}},
						},
					},
					Parameters: []*v1.Parameter{{
						Name:     "foo",
						Required: true,
					}},
				},
				Layers: []render.LayerInput{{LayerId: "a", OptionId: "1"}},
				Params: map[string]string{"foo": "bar", "bar": "baz"},
			}
			version := v1.VersionedApplicationSpec{
				Parameters: []*v1.Parameter{
					{
						Name:     "bar",
						Required: true,
					},
					{
						Name: "optional",
					},
				},
			}
			err := render.ValidateInputs(inputs, version, func([]*v1.ResourceDependency) error { return errors.Errorf("test") })
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("test"))
		})
	})

	Context("render templates in input values", func() {

		inputs := render.ValuesInputs{
			InstallNamespace: "test-ns",
			MeshRef: core.ResourceRef{
				Namespace: "mesh-ns",
				Name:      "my-mesh",
			},
			UserDefinedValues: "top:\n  nested: {{ .InstallNamespace }}\n",
			Params: map[string]string{
				"my.app.mesh-ref": "{{ .MeshRef.Name }}.{{ .MeshRef.Namespace }}",
				"unchanged":       "still-the-same",
			},
		}

		expected := render.ValuesInputs{
			InstallNamespace: "test-ns",
			MeshRef: core.ResourceRef{
				Namespace: "mesh-ns",
				Name:      "my-mesh",
			},
			UserDefinedValues: "top:\n  nested: test-ns\n",
			Params: map[string]string{
				"my.app.mesh-ref": "my-mesh.mesh-ns",
				"unchanged":       "still-the-same",
			},
		}

		It("correctly renders input values", func() {
			result, err := render.ExecInputValuesTemplates(inputs)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEquivalentTo(expected))
		})

	})
})
