package gon

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Definitions map[string]Value

	DefinitionResolver interface {
		Definition(key string) (Value, bool)
	}

	definitionResolver struct {
		store map[string]Value
	}
)

var nameRegex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]{1,50}$")

func (r *definitionResolver) Definition(key string) (Value, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := r.store[topKey]
	if !ok {
		return Literal(fmt.Errorf("definition '%s' not found", topKey)), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(DefinitionResolver)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return Literal(fmt.Errorf("definition '%s' doesn't have children attributes", topKey)), false
}

func (r *definitionResolver) Define(key string, definition Value) error {
	if !nameRegex.MatchString(key) {
		return NodeError{
			Scalar: "call",
			Cause:  fmt.Errorf("definition key '%s' is invalid", key),
		}
	}

	r.store[key] = definition

	return nil
}

var _ DefinitionResolver = &definitionResolver{}
