package gon

import (
	"errors"
	"fmt"
)

type (
	NodeError struct {
		NodeScalar string
		Cause      error
	}

	DefinitionNotFoundError struct {
		DefinitionName string
	}

	DefinitionNotCallable struct {
		DefinitionName string
	}
)

func NewNodeError(namedNode Named, err error) NodeError {
	return NodeError{
		NodeScalar: namedNode.Scalar(),
		Cause:      err,
	}
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
	return fmt.Sprintf("definition '%s' not found", e.DefinitionName)
}

func (e DefinitionNotCallable) Error() string {
	return fmt.Sprintf("definition '%s' not callable", e.DefinitionName)
}

var (
	_ Node  = NodeError{}
	_ Value = NodeError{}
)
