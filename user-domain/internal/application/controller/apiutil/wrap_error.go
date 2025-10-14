package apiutil

import "fmt"

func WrapError(base, errContext error) error {
	return fmt.Errorf("%v: %w", base, errContext)
}
