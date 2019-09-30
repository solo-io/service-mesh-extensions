package util

import (
	"strconv"

	"github.com/solo-io/go-utils/errors"
	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

var (
	UnsupportedParamType = func(paramValueType interface{}) error {
		return errors.Errorf("The operator does not currently support params of type %T", paramValueType)
	}
)

type GetSecretValue func(value *v1.SecretValue) (string, error)

func PlainTextSecretGetter(value *v1.SecretValue) (string, error) {
	switch t := value.GetType().(type) {
	case *v1.SecretValue_PlainText:
		return t.PlainText, nil
	default:
		// TODO support file secrets
		return "", UnsupportedParamType(t)
	}
}

func ParamValueToString(value *v1.ParameterValue, getSecret GetSecretValue) (string, error) {
	switch t := value.GetType().(type) {
	case *v1.ParameterValue_BooleanValue:
		if t.BooleanValue {
			return "true", nil
		}
		return "false", nil
	case *v1.ParameterValue_DateValue:
		return t.DateValue.String(), nil
	case *v1.ParameterValue_FloatValue:
		return strconv.FormatFloat(t.FloatValue, 'f', -1, 64), nil
	case *v1.ParameterValue_IntValue:
		return strconv.Itoa(int(t.IntValue)), nil
	case *v1.ParameterValue_SecretValue:
		return getSecret(t.SecretValue)
	case *v1.ParameterValue_StringValue:
		return t.StringValue, nil
	default:
		return "", UnsupportedParamType(t)
	}
}
