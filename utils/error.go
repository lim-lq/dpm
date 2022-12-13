package utils

import "errors"

func DpmError(msg string) error {
	return errors.New(msg)
}
