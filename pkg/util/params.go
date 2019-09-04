package util

import (
	"strconv"

	v1 "github.com/solo-io/service-mesh-hub/api/v1"
)

func GetDefaultString(param *v1.Parameter) string {
	if param.Default == nil {
		return ""
	}
	switch t := param.Default.Type.(type) {
	case *v1.ParameterValue_String_:
		return t.String_
	case *v1.ParameterValue_Bool:
		if t.Bool {
			return "true"
		}
		return "false"
	case *v1.ParameterValue_Int:
		return strconv.Itoa(int(t.Int))
	case *v1.ParameterValue_Float:
		return strconv.FormatFloat(t.Float, 'E', -1, 64)
	case *v1.ParameterValue_Date:
		return t.Date.String()
	case *v1.ParameterValue_Secret:
		// Default not supported
		return ""
	}
	return ""
}
