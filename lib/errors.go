package lib

import "github.com/pkg/errors"

type UnsupportedError struct {
	message string
}

func NewUnsupportedError(msg string) error {
	return errors.Wrap(&UnsupportedError{msg}, "")
}

func (ue *UnsupportedError) Error() string {
	return ue.message
}

type InvalidAstError struct {
	message string
}

func NewInvalidAstError(msg string) error {
	return errors.Wrap(&InvalidAstError{msg}, "")
}

func (ie *InvalidAstError) Error() string {
	return ie.message
}
