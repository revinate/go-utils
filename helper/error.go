package helper

import (
	"fmt"
	"strings"
)

func ErrorsToString(errs ...error) string {
	errStrs := []string{}
	for _, err := range errs {
		if err != nil {
			errStrs = append(errStrs, err.Error())
		}
	}
	return strings.Join(errStrs, ", ")
}

func MergeErrors(errs ...error) error {
	errString := ErrorsToString(errs...)
	if errString != "" {
		return fmt.Errorf("Error: %s", errString)
	}
	return nil
}
