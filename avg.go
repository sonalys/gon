package gon

import "fmt"

type avgNode struct {
	nodes []Node
}

func Avg(nodes ...Node) Node {
	if len(nodes) == 0 {
		return Literal(NodeError{
			NodeName: "avg",
			Cause:    fmt.Errorf("must receive at least one expression"),
		})
	}

	for i := range nodes {
		if nodes[i] == nil {
			return Literal(NodeError{
				NodeName: "avg",
				Cause:    fmt.Errorf("all expressions should be not-nil"),
			})
		}
	}

	return avgNode{
		nodes: nodes,
	}
}

func (node avgNode) Name() string {
	return "avg"
}

func (node avgNode) Shape() []KeyNode {
	return nil
}

func (node avgNode) Type() NodeType {
	return NodeTypeExpression
}

func (node avgNode) Eval(scope Scope) Value {
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

	sum, ok := avgAny(values...)
	if !ok {
		return Literal(NodeError{
			NodeName: node.Name(),
			Cause:    fmt.Errorf("all nodes must be of the same type"),
		})
	}

	return Literal(sum)
}
