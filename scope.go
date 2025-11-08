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
		DefinitionResolver
	}

	scope struct {
		definitionStore
		context.Context

		parentScope Scope
	}
)

// NewScope initializes a new scope.
// It starts with a background context by default.
func NewScope() *scope {
	return &scope{
		Context:         context.Background(),
		definitionStore: definitionStore{store: make(Definitions)},
	}
}

func (s *scope) WithContext(ctx context.Context) *scope {
	s.Context = ctx
	return s
}

func (s *scope) WithDefinitions(source Definitions) (*scope, error) {
	for key, value := range source {
		if !nameRegex.MatchString(key) {
			return nil, fmt.Errorf("invalid definition name: %s", key)
		}
		if err := s.Define(key, value); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *scope) Definition(key string) (Value, bool) {
	value, ok := s.definitionStore.Definition(key)
	if !ok {
		if s.parentScope != nil {
			return s.parentScope.Definition(key)
		}
		return value, false
	}

	return value, true
}

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
