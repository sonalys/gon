package gon

import (
	"regexp"
	"strings"
)

type (
	// Definitions defines how scope definitions are configured.
	// The key must be an alphanumeric+underscore+dash string from length 1 to 50, starting with a letter.
	Definitions map[string]Value

	DefinitionReader interface {
		Definition(key string) (Value, bool)
	}

	DefinitionWriter interface {
		Define(key string, value Value) error
	}

	DefinitionReadWriter interface {
		DefinitionReader
		DefinitionWriter
	}

	definitionStore struct {
		store map[string]Value
	}
)

var keyValidationRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]{1,50}$`)

func newDefinitionResolver() DefinitionReadWriter {
	return &definitionStore{
		store: make(map[string]Value),
	}
}

func (r *definitionStore) Definition(key string) (Value, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := r.store[topKey]
	if !ok {
		return Literal(DefinitionNotFoundError{
			DefinitionKey: topKey,
		}), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(DefinitionReader)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return Literal(DefinitionNotFoundError{
		DefinitionKey: key,
	}), false
}

func (r *definitionStore) Define(key string, value Value) error {
	if !keyValidationRegex.MatchString(key) {
		return InvalidDefinitionKey{
			DefinitionKey: key,
		}
	}

	r.store[key] = value

	return nil
}

var _ DefinitionReader = &definitionStore{}
