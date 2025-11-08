package gon

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

	// Named defines a named node.
	Named interface {
		Name() string
	}

	// Shaped defines a node capable of returning it's shape.
	// The shape is a slice of named or unamed parameters required for constructing the node.
	// This abstraction is used for constructing nodes from named parameters.
	Shaped interface {
		Shape() []KeyNode
	}

	// Computable represents any abstraction capable of evaluating under a scope.
	Computable interface {
		Eval(scope Scope) Value
	}

	// Node is the building block of any expression.
	// It can be used to represent values, evaluations or operations.
	// Nodes can be evaluated under a scope.
	Node interface {
		Named
		Shaped
		Computable
	}

	// KeyNode defines a key-node pair, used for named parameters.
	KeyNode struct {
		Key  string
		Node Node
	}
)
