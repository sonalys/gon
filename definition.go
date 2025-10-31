package gon

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	definition struct {
		key string
	}

	Definitions map[string]Expression

	DefinitionResolver interface {
		Definition(key string) (Expression, bool)
	}

	Definer interface {
		Define(key string, expression Expression) error
	}

	definitionResolver struct {
		store map[string]Expression
	}

	assignment struct {
		definition string
		expression Expression
	}
)

func (d definition) Banner() (string, []KeyExpression) {
	return d.key, nil
}

func (d definition) Type() NodeType {
	return NodeTypeReference
}

func Reference(key string) definition {
	return definition{
		key: key,
	}
}

func (d definition) Eval(scope Scope) Value {
	expression, ok := scope.Definition(d.key)
	if !ok {
		return Static(fmt.Errorf("definition not found: %s", d.key))
	}

	return expression.Eval(scope)
}

func (a assignment) Eval(scope Scope) Value {
	definer, ok := scope.(Definer)
	if !ok {
		return Static(errors.New("scope is read-only"))
	}

	return Static(definer.Define(a.definition, a.expression))
}

func (s *definitionResolver) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := s.store[topKey]
	if !ok {
		return Static(fmt.Errorf("definition not found: %s", topKey)), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return propagateErr(nil, "definition doesn't have children"), false
}

var definitionNameRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$")

func (s *definitionResolver) Define(key string, expression Expression) error {
	if !definitionNameRegex.MatchString(key) {
		return fmt.Errorf("invalid definition key: %s", key)
	}

	s.store[key] = expression

	return nil
}

var _ DefinitionResolver = &definitionResolver{}
var _ Expression = &definition{}
