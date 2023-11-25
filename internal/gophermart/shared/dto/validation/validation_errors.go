package validation

import "strings"

type ValidationErrors []error

func (ve ValidationErrors) String() string {
	errStrings := make([]string, len(ve))
	for i, err := range ve {
		errStrings[i] = err.Error()
	}
	return strings.Join(errStrings, "\n")
}
