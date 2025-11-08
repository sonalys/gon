package gon

import "fmt"

type sumNode struct {
	nodes []Node
}

func Sum(nodes ...Node) Node {
	if len(nodes) == 0 {
		return Literal(NodeError{
			NodeName: "sum",
			Cause:    fmt.Errorf("must receive at least one expression"),
		})
	}

	for i := range nodes {
		if nodes[i] == nil {
			return Literal(NodeError{
				NodeName: "sum",
				Cause:    fmt.Errorf("all expressions should be not-nil"),
			})
		}
	}

	return sumNode{
		nodes: nodes,
	}
}

func (node sumNode) Name() string {
	return "sum"
}

func (node sumNode) Shape() []KeyNode {
	return nil
}

func (node sumNode) Type() NodeType {
	return NodeTypeExpression
}

func (node sumNode) Eval(scope Scope) Value {
	values := make([]any, 0, len(node.nodes))

	for i := range node.nodes {
		curValue, err := scope.Compute(node.nodes[i])
		if err != nil {
			return Literal(NodeError{
				NodeName: node.Name(),
				Cause:    err,
			})
		}

		values = append(values, curValue)
	}

	sum, ok := sumAny(values...)
	if !ok {
		return Literal(NodeError{
			NodeName: node.Name(),
			Cause:    fmt.Errorf("all nodes must be of the same type"),
		})
	}

	return Literal(sum)
}
