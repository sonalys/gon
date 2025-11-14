package gon

import (
	"context"
	"fmt"
)

type (
	// Scope defines a block capable of evaluating expressions.
	// It should be able to act as a context, as well as resolve definitions.
	Scope interface {
		context.Context
		DefinitionReader
		Compute(Node) (any, error)
	}

	scope struct {
		store DefinitionReadWriter
		context.Context

		parentScope Scope
	}
)

// NewScope initializes a new scope.
// A Scope can be used to evaluate expressions under specific conditions.
// It can also define context for evaluation and define data for the expressions.
// It starts with a background context by default.
func NewScope() *scope {
	return &scope{
		Context: context.Background(),
		store:   newDefinitionResolver(),
	}
}

func (s *scope) WithContext(ctx context.Context) *scope {
	s.Context = ctx
	return s
}

func (s *scope) WithDefinitions(source Definitions) (*scope, error) {
	for key, value := range source {
		if !keyValidationRegex.MatchString(key) {
			return nil, fmt.Errorf("invalid definition name: %s", key)
		}
		if err := s.store.Define(key, value); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *scope) Definition(key string) (Value, bool) {
	value, ok := s.store.Definition(key)
	if !ok {
		if s.parentScope != nil {
			return s.parentScope.Definition(key)
		}
		return value, false
	}

	return value, true
}

// Compute will evaluate the final value for the root node.
// If the value is of type error, it will be returned as error instead.
func (s *scope) Compute(node Node) (any, error) {
	result := node.Eval(s)
	switch t := result.Value().(type) {
	case error:
		return nil, t
	default:
		return t, nil
	}
}

var (
	_ Scope = &scope{}
)
