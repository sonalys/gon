package gon

import (
	"errors"
	"strings"
)

type (
	definition struct {
		key string
	}

	DefinitionResolver interface {
		Definition(key string) (Expression, bool)
	}

	definitionResolver map[string]Expression

	assignment struct {
		definition string
		expression Expression
	}
)

func Definition(key string) definition {
	return definition{
		key: key,
	}
}

func (d definition) Eval(scope Scope) Value {
	return scope.Definition(d.key).Eval(scope)
}

func (a assignment) Eval(scope Scope) Value {
	scope.Define(a.definition, a.expression)
	return Static(nil)
}

func (s definitionResolver) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := s[topKey]
	if !ok {
		return Static(errors.New("definition not found")), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return Static(errors.New("definition doesn't have children")), false
}
