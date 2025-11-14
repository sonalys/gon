package encoding

import (
	"fmt"

	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/nodes"
)

type (
	NodeConstructor func(keyedNodes []adapters.KeyNode) (adapters.Node, error)
	Codex           map[string]NodeConstructor
	DecodeConfig    struct {
		NodeCodex Codex
	}
)

var DefaultExpressionCodex = Codex{}

func (c *Codex) Register(name string, constructor func([]adapters.KeyNode) (adapters.Node, error)) error {
	if _, conflicts := (*c)[name]; conflicts {
		return fmt.Errorf("node with name '%s' is already registered", name)
	}

	(*c)[name] = constructor

	return nil
}

type AutoRegisterer interface {
	Register(codex adapters.Codex) error
}

func (c *Codex) AutoRegister(nodes ...AutoRegisterer) error {
	for _, registerer := range nodes {
		if err := registerer.Register(c); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	err := DefaultExpressionCodex.AutoRegister(
		&nodes.AvgNode{},
		&nodes.CallNode{},
		&nodes.EqualNode{},
		&nodes.GreaterNode{},
		&nodes.HasPrefixNode{},
		&nodes.HasSuffixNode{},
		&nodes.IfNode{},
		&nodes.LiteralNode{},
		&nodes.NotNode{},
		&nodes.OrNode{},
		&nodes.SmallerNode{},
		&nodes.SumNode{},
		&nodes.IsEmptyNode{},
	)
	if err != nil {
		panic(fmt.Errorf("unexpected error registering default nodes: %s", err))
	}
}
