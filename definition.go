package gon

import (
	"fmt"
	"regexp"
	"strings"
)

type (
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
)

var nameRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{1,50}$")

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

func (r *definitionResolver) Define(key string, expression Expression) error {
	if !nameRegex.MatchString(key) {
		return fmt.Errorf("invalid definition key: %s", key)
	}

	r.store[key] = expression

	return nil
}

var _ DefinitionResolver = &definitionResolver{}
