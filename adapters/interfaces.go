package adapters

import (
	"context"
)

type (
	// Valued defines a value wrapper, any struct capable of returning it's internal value.
	Valued interface {
		Value() any
	}

	// Value represents a value node, it can be used inside expressions or as definitions.
	Value interface {
		Node
		Valued
	}

	// Typed defines a node capable of returning it's NodeType.
	// This interface is used for encoding and syntax purposes.
	Typed interface {
		Type() NodeType
	}

	// Named abstracts a node scalar getter.
	Named interface {
		Scalar() string
	}

	// Shaped defines a node capable of returning it's shape.
	// The shape is a slice of named or unamed parameters required for constructing the node.
	// This abstraction is used for constructing nodes from named parameters.
	Shaped interface {
		Shape() []KeyNode
	}

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

	// Scope defines a block capable of evaluating expressions.
	// It should be able to act as a context, as well as resolve definitions.
	Scope interface {
		context.Context
		DefinitionReader
		Compute(Node) (any, error)
	}

	// Node is the building block of any expression.
	// It can be used to represent values, evaluations or operations.
	// Nodes can be evaluated under a scope.
	Node interface {
		Named
		Eval(scope Scope) Value
	}

	// KeyNode defines a key-node pair, used for named parameters.
	KeyNode struct {
		Key  string
		Node Node
	}

	Codex interface {
		Register(name string, constructor func([]KeyNode) (Node, error)) error
	}

	// Callable defines a node that can be called.
	// It represents a function as a node.
	Callable interface {
		Node
		Call(ctx context.Context, funcName string, argValues ...Value) Value
	}
)
