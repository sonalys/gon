package gon

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	// Definitions defines how scope definitions are configured.
	// The key must be an alphanumeric+underscore+dash string from length 1 to 50, starting with a letter.
	Definitions map[string]Value

	DefinitionResolver interface {
		Definition(key string) (Value, bool)
	}

	definitionStore struct {
		store map[string]Value
	}
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]{1,50}$`)

func (r *definitionStore) Definition(key string) (Value, bool) {
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

func (r *definitionStore) Define(key string, value Value) error {
	if !nameRegex.MatchString(key) {
		return fmt.Errorf("definition key '%s' is invalid", key)
	}

	r.store[key] = value

	return nil
}

var _ DefinitionResolver = &definitionStore{}
