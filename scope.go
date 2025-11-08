package gon

import (
	"context"
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
		expression  Node
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
	// Empty store can be a direct copy.
	if len(s.definitionStore.store) == 0 {
		s.definitionStore.store = source
		return s, nil
	}

	for key, value := range source {
		if err := s.definitionStore.Define(key, value); err != nil {
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

var (
	_ Scope = &scope{}
)
