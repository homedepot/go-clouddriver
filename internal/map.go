package internal

import (
	"fmt"
	"strings"
)

// DeleteNilVaules returns the given map with all keys
// with nil values removed.
func DeleteNilValues(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})

	for k, v := range m {
		if v != nil {
			if isMap(v) {
				newMap[k] = DeleteNilValues(v.(map[string]interface{}))
			} else {
				newMap[k] = v
			}
		}
	}

	return newMap
}

// isMap returns true if the argument is of type map
//
// see https://stackoverflow.com/questions/20759803/how-to-check-variable-type-is-map-in-go-language
func isMap(x interface{}) bool {
	t := fmt.Sprintf("%T", x)
	return strings.HasPrefix(t, "map[")
}
