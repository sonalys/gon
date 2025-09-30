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
		Variable(key string) Expression
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

func (s *scope) WithVariables(source map[string]Expression) *scope {
	maps.Insert(s.store, maps.All(source))
	s.store = source
	return s
}

func (s *scope) WithContext(ctx context.Context) *scope {
	s.Context = ctx
	return s
}

func (s *scope) Variable(key string) Expression {
	value, ok := s.store[key]
	if !ok {
		if s.parentScope != nil {
			return s.parentScope.Variable(key)
		}
		return Static(errors.New("fact not found"))
	}

	return value
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
