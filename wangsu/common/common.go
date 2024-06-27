package common

import (
	"strings"
)

// IsContains returns whether value is within array
func IsContains(array interface{}, value interface{}) bool {
	switch arr := array.(type) {
	case []string:
		for _, v := range arr {
			if v == value {
				return true
			}
		}
	case map[string]interface{}:
		if _, ok := arr[value.(string)]; ok {
			return true
		}
	case string:
		if strings.Contains(arr, value.(string)) {
			return true
		}
	}
	return false
}
