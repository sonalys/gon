package gon

import (
	"errors"
	"fmt"
)

type NodeError struct {
	Scalar string
	Cause  error
}

func (e NodeError) Error() string {
	if errors.As(e.Cause, &NodeError{}) {
		return fmt.Sprintf("%s.%s", e.Scalar, e.Cause.Error())
	}

	return fmt.Sprintf("%s: %s", e.Scalar, e.Cause.Error())
}

func (e NodeError) Unwrap() error {
	return e.Cause
}

type DefinitionNotFoundError struct {
	DefinitionName string
}

func (e DefinitionNotFoundError) Error() string {
	return fmt.Sprintf("definition '%s' not found", e.DefinitionName)
}

type DefinitionNotCallable struct {
	DefinitionName string
}

func (e DefinitionNotCallable) Error() string {
	return fmt.Sprintf("definition '%s' not callable", e.DefinitionName)
}
