package main_test

import (
	"bytes"
	"fmt"
	"maps"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

type customNode struct {
	input gon.Node
}

func (node *customNode) Type() gon.NodeType {
	return gon.NodeTypeExpression
}

func (node *customNode) Name() string {
	return "customNode"
}

func (node *customNode) Shape() []gon.KeyNode {
	return []gon.KeyNode{
		{Key: "myCustomParam", Node: node.input},
	}
}

func (node *customNode) Eval(scope gon.Scope) gon.Value {
	valued := node.input.Eval(scope)

	fmt.Printf("got: %v\n", valued.Value())

	return gon.Literal(true)
}

var (
	_ gon.Node = &customNode{}
)

func Example_customNode() {
	myExpression := gon.If(&customNode{input: gon.Literal("my-param")}, gon.Literal("works!"))

	buffer := bytes.NewBuffer(make([]byte, 0))

	err := encoding.Encode(buffer, myExpression)
	if err != nil {
		panic(err)
	}

	customCodex := maps.Clone(encoding.DefaultExpressionCodex)

	customCodex["customNode"] = func(args []gon.KeyNode) (gon.Node, error) {
		return &customNode{
			input: args[0].Node,
		}, nil
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

	valued := decodedNode.Eval(scope)
	fmt.Println(valued.Value())

	//Output:
	// if(
	// 	condition: customNode(myCustomParam: "my-param")
	// 	then: "works!"
	// )
	// got: my-param
	// works!
}
