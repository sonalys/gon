package gon

import (
	"context"
)

type (
	Scope interface {
		context.Context
		DefinitionResolver
	}

	scope struct {
		definitionResolver
		context.Context

		parentScope Scope
		expression  Expression
	}
)

func NewScope() *scope {
	return &scope{
		Context:            context.Background(),
		definitionResolver: definitionResolver{store: make(Definitions)},
	}
}

func (s *scope) WithDefinitions(source Definitions) (*scope, error) {
	for key, value := range source {
		if err := s.definitionResolver.Define(key, value); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *scope) WithContext(ctx context.Context) *scope {
	s.Context = ctx
	return s
}

func (s *scope) Definition(key string) (Value, bool) {
	value, ok := s.definitionResolver.Definition(key)
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
