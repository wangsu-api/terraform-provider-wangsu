package common

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateAllowedStringValue checks if a string is in a slice of strings.
func ValidateAllowedStringValue(ss []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !IsContains(ss, value) {
			errors = append(errors, fmt.Errorf("%q must contain a valid string value must in array %#v, got %q", k, ss, value))
		}
		return
	}
}
