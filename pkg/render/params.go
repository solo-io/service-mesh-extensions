package render

import (
	"context"
	"fmt"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/ghodss/yaml"

	"github.com/helm/helm/pkg/strvals"
	"github.com/pkg/errors"
)

var (
	UnableToParseParameterError = func(err error, key, value string) error {
		return errors.Wrapf(err, "Unable to parse parameter with key '%s' and value '%s'", key, value)
	}

	UnableToParseYamlError = func(err error, input string) error {
		return errors.Wrapf(err, "Unable to parse yaml string: %s", input)
	}

	UnableToMarshalYamlError = func(err error, input map[string]interface{}) error {
		return errors.Wrapf(err, "Unable to marshal map to yaml: %v", input)
	}
)

func ConvertParamsToNestedMap(params map[string]string) (map[string]interface{}, error) {
	nestedMap := make(map[string]interface{})
	for k, v := range params {
		merged := fmt.Sprintf("%s=%s", k, v)
		err := strvals.ParseInto(merged, nestedMap)
		if err != nil {
			return nil, UnableToParseParameterError(err, k, v)
		}
	}
	return nestedMap, nil
}

func ConvertYamlStringToNestedMap(yamlString string) (map[string]interface{}, error) {
	nestedMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(yamlString), &nestedMap)
	if err != nil {
		return nil, UnableToParseYamlError(err, yamlString)
	}
	return nestedMap, nil
}

func ConvertNestedMapToYaml(nestedMap map[string]interface{}) (string, error) {
	if nestedMap == nil || len(nestedMap) == 0 {
		return "", nil
	}
	yamlString, err := yaml.Marshal(nestedMap)
	if err != nil {
		return "", UnableToMarshalYamlError(err, nestedMap)
	}
	return string(yamlString), nil
}

func CoalesceValuesMap(ctx context.Context, initial map[string]interface{}, overrides map[string]interface{}) map[string]interface{} {
	// this helper prefers the first map over the second
	return coalesceTables(ctx, overrides, initial)
}

// istable is a special-purpose function to see if the present thing matches the definition of a YAML table.
func istable(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// coalesceTables merges a source map into a destination map.
//
// dest is considered authoritative.
func coalesceTables(ctx context.Context, dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	if src == nil {
		src = make(map[string]interface{})
	}
	// Because dest has higher precedence than src, dest values override src
	// values.
	for key, val := range src {
		if istable(val) {
			if innerdst, ok := dst[key]; !ok {
				dst[key] = val
			} else if istable(innerdst) {
				coalesceTables(ctx, innerdst.(map[string]interface{}), val.(map[string]interface{}))
			} else {
				contextutils.LoggerFrom(ctx).Errorw("coalescing table into value, dropping table",
					zap.String("key", key),
					zap.Any("value", innerdst),
					zap.Any("table", val))
			}
			continue
		} else if dv, ok := dst[key]; ok && istable(dv) {
			contextutils.LoggerFrom(ctx).Errorw("coalescing value into table, dropping value",
				zap.String("key", key),
				zap.Any("value", val),
				zap.Any("table", dv))
			continue
		} else if !ok { // <- ok is still in scope from preceding conditional.
			dst[key] = val
			continue
		}
	}
	return dst
}
