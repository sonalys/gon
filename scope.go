package gon

import (
	"context"
	"errors"
	"maps"
)

type (
	Scope interface {
		context.Context
		Expression

		Definition(key string) Expression
		Define(key string, expression Expression)
	}

	scope struct {
		context.Context
		parentScope Scope
		expression  Expression
		store       map[string]Expression
	}
)

func NewScope() *scope {
	return &scope{
		Context: context.Background(),
		store:   make(map[string]Expression),
	}
}

func (s *scope) WithDefinitions(source map[string]Expression) *scope {
	maps.Insert(s.store, maps.All(source))
	s.store = source
	return s
}

func (s *scope) WithContext(ctx context.Context) *scope {
	s.Context = ctx
	return s
}

func (s *scope) Definition(key string) Expression {
	value, ok := s.store[key]
	if !ok {
		if s.parentScope != nil {
			return s.parentScope.Definition(key)
		}
		return Static(errors.New("fact not found"))
	}

	return value
}

func (s *scope) Define(key string, expression Expression) {
	s.store[key] = expression
}

func (s *scope) Eval(scope Scope) Value {
	s.parentScope = scope

	if s.expression != nil {
		return s.expression.Eval(scope)
	}

	return Static(error(nil))
}

var (
	_ Scope = &scope{}
)
