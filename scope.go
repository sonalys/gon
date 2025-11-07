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

func (s *scope) Type() NodeType {
	return NodeTypeInvalid
}

func NewScope() *scope {
	return &scope{
		Context:            context.Background(),
		definitionResolver: definitionResolver{store: make(map[string]Expression)},
	}
}

func (s *scope) WithDefinitions(source map[string]Expression) (*scope, error) {
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

func (s *scope) Definition(key string) (Expression, bool) {
	value, ok := s.definitionResolver.Definition(key)
	if !ok {
		if s.parentScope != nil {
			return s.parentScope.Definition(key)
		}
		return nil, false
	}

	return value, true
}

func (s *scope) Eval(scope Scope) Value {
	s.parentScope = scope

	if s.expression != nil {
		return s.expression.Eval(scope)
	}

	return Literal(error(nil))
}

var (
	_ Scope = &scope{}
)
