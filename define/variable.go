package define

import (
	"fmt"
	"strings"
)

const WEBAPI_VARS string = "$WAPARAM"

func ExpandVariables(src string, vars ...string) string {
	if len(vars) == 0 {
		return src
	}
	result := src
	for i, v := range vars {
		parm := fmt.Sprintf("%s%d", WEBAPI_VARS, i + 1)
		result = strings.Replace(result, parm, v, -1)
	}
	return result
}
