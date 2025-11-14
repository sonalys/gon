package gon

import (
	"errors"
	"fmt"
)

type (
	StringError string

	NodeError struct {
		NodeScalar string
		Cause      error
	}

	DefinitionNotFoundError struct {
		DefinitionKey string
	}

	DefinitionNotCallableError struct {
		DefinitionKey string
	}

	InvalidDefinitionKey struct {
		DefinitionKey string
	}

	IncompatiblePairError struct {
		First  any
		Second any
	}

	DefinitionError interface {
		Key() string
	}
)

const (
	ErrAllNodesMustMatch StringError = "all nodes must be of the same type"
	ErrAllNodesMustBeSet StringError = "all nodes must be set"
	ErrMustHaveArguments StringError = "must have at least one argument"
)

func NewNodeError(namedNode Named, err error) NodeError {
	return NodeError{
		NodeScalar: namedNode.Scalar(),
		Cause:      err,
	}
}

func (e StringError) Error() string {
	return string(e)
}

func (e NodeError) Value() any {
	return e
}

func (e NodeError) Eval(scope Scope) Value {
	return e
}

func (e NodeError) Scalar() string {
	return "nodeError"
}

func (e NodeError) Error() string {
	if errors.As(e.Cause, &NodeError{}) {
		return fmt.Sprintf("%s.%s", e.NodeScalar, e.Cause.Error())
	}

	return fmt.Sprintf("%s: %s", e.NodeScalar, e.Cause.Error())
}

func (e NodeError) Unwrap() error {
	return e.Cause
}

func (e DefinitionNotFoundError) Error() string {
	return fmt.Sprintf("definition '%s' not found", e.DefinitionKey)
}

func (e DefinitionNotCallableError) Error() string {
	return fmt.Sprintf("definition '%s' not callable", e.DefinitionKey)
}

func (e IncompatiblePairError) Error() string {
	return fmt.Sprintf("types %T and %T are not compatible", e.First, e.Second)
}

func (e InvalidDefinitionKey) Error() string {
	return fmt.Sprintf("definition key '%s' is invalid", e.DefinitionKey)
}

func (e DefinitionNotCallableError) Key() string {
	return e.DefinitionKey
}

func (e DefinitionNotFoundError) Key() string {
	return e.DefinitionKey
}

func (e InvalidDefinitionKey) Key() string {
	return e.DefinitionKey
}

var (
	_ Node            = NodeError{}
	_ Value           = NodeError{}
	_ DefinitionError = DefinitionNotCallableError{}
	_ DefinitionError = DefinitionNotFoundError{}
	_ DefinitionError = InvalidDefinitionKey{}
)
