package gon

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type (
	referenceNode struct {
		definitionName string
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

func (node referenceNode) Name() string {
	return node.definitionName
}

func (node referenceNode) Shape() []KeyExpression {
	return nil
}

func (node referenceNode) Type() NodeType {
	return NodeTypeReference
}

func Reference(key string) referenceNode {
	return referenceNode{
		definitionName: key,
	}
}

func (node referenceNode) Eval(scope Scope) Value {
	expression, ok := scope.Definition(node.definitionName)
	if !ok {
		return Literal(fmt.Errorf("definition not found: %s", node.definitionName))
	}

	return expression.Eval(scope)
}

func (a assignment) Eval(scope Scope) Value {
	definer, ok := scope.(Definer)
	if !ok {
		return Literal(errors.New("scope is read-only"))
	}

	return Literal(definer.Define(a.definition, a.expression))
}

func (r *definitionResolver) Definition(key string) (Expression, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := r.store[topKey]
	if !ok {
		return Literal(fmt.Errorf("definition not found: %s", topKey)), false
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

var nameRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{1,50}$")

func (r *definitionResolver) Define(key string, expression Expression) error {
	if !nameRegex.MatchString(key) {
		return fmt.Errorf("invalid definition key: %s", key)
	}

	r.store[key] = expression

	return nil
}

var _ DefinitionResolver = &definitionResolver{}
var _ Expression = &referenceNode{}
