package main_test

import (
	"bytes"
	"fmt"
	"maps"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/encoding"
	"github.com/sonalys/gon/gonutils"
)

type customNode struct {
	input adapters.Node
}

func (node *customNode) Type() adapters.NodeType {
	return adapters.NodeTypeExpression
}

func (node *customNode) Scalar() string {
	return "customNode"
}

func (node *customNode) Shape() []adapters.KeyNode {
	return []adapters.KeyNode{
		{Key: "myCustomParam", Node: node.input},
	}
}

func (node *customNode) Eval(scope adapters.Scope) adapters.Value {
	valued := node.input.Eval(scope)

	fmt.Printf("got: %v\n", valued.Value())

	return gon.Literal(true)
}

func (node *customNode) Register(codex adapters.Codex) error {
	return codex.Register(node.Scalar(), func(args []adapters.KeyNode) (adapters.Node, error) {
		orderedArgs, _, err := gonutils.SortArgs(args, "myCustomParam")
		if err != nil {
			return nil, err
		}

		return &customNode{
			input: orderedArgs["myCustomParam"],
		}, nil
	})
}

var (
	// Ensure these constraints are met if you want your node to encode/decode properly.
	_ adapters.SerializableNode = &customNode{}
)

func Example_customNode() {
	myExpression := gon.If(&customNode{input: gon.Literal("my-param")}, gon.Literal("works!"))

	buffer := bytes.NewBuffer(make([]byte, 0))

	err := encoding.HumanEncode(buffer, myExpression)
	if err != nil {
		panic(err)
	}

	customCodex := maps.Clone(encoding.DefaultExpressionCodex)

	err = customCodex.AutoRegister(&customNode{})
	if err != nil {
		panic(err)
	}

	decodedNode, err := encoding.Decode(buffer.Bytes(), customCodex)
	if err != nil {
		panic(err)
	}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"var": gon.Literal("my-var"),
		})
	if err != nil {
		panic(err)
	}

	fmt.Println(buffer.String())

	value, err := scope.Compute(decodedNode)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	//Output:
	// if(
	// 	condition: customNode(myCustomParam: "my-param"),
	// 	then: "works!"
	// )
	// got: my-param
	// works!
}
