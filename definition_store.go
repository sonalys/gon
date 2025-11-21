package gon

import (
	"regexp"
	"strings"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/nodes"
)

type (
	// Values defines how scope definitions are configured.
	// The key must be an alphanumeric+underscore+dash string from length 1 to 50, starting with a letter.
	Values map[string]adapters.Value

	definitionStore struct {
		store map[string]adapters.Value
	}
)

var keyValidationRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]{1,50}$`)

func newDefinitionResolver() adapters.DefinitionReadWriter {
	return &definitionStore{
		store: make(map[string]adapters.Value),
	}
}

func (r *definitionStore) Definition(key string) (adapters.Value, bool) {
	parts := strings.Split(key, ".")
	topKey := parts[0]

	value, ok := r.store[topKey]
	if !ok {
		return nodes.Literal(adapters.DefinitionNotFoundError{
			DefinitionKey: topKey,
		}), false
	}

	if len(parts) == 1 {
		return value, true
	}

	resolver, isResolver := value.(adapters.DefinitionReader)
	if isResolver {
		return resolver.Definition(key[len(topKey)+1:])
	}

	return nodes.Literal(adapters.DefinitionNotFoundError{
		DefinitionKey: key,
	}), false
}

func (r *definitionStore) Define(key string, value adapters.Value) error {
	if !keyValidationRegex.MatchString(key) {
		return adapters.InvalidDefinitionKey{
			DefinitionKey: key,
		}
	}

	r.store[key] = value

	return nil
}

var _ adapters.DefinitionReader = &definitionStore{}
