package encoding

import (
	"fmt"
	"time"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/internal/sliceutils"
)

var DefaultExpressionCodex = Codex{
	"if": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, rest, err := argSorter(args, "condition", "then")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'if' node: %w", err)
		}

		return gon.If(orderedArgs["condition"], orderedArgs["then"], rest...), nil
	},
	"or": func(args []gon.KeyNode) (gon.Node, error) {
		_, argsSlice, _ := argSorter(args)

		return gon.Or(argsSlice...), nil
	},
	"equal": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'equal' node: %w", err)
		}
		return gon.Equal(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"lt": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'lt' node: %w", err)
		}
		return gon.Smaller(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"lte": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'lte' node: %w", err)
		}
		return gon.SmallerOrEqual(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"gt": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'gt' node: %w", err)
		}
		return gon.Greater(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"gte": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "first", "second")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'gte' node: %w", err)
		}
		return gon.GreaterOrEqual(orderedArgs["first"], orderedArgs["second"]), nil
	},
	"not": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "expression")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'not' node: %w", err)
		}
		return gon.Not(orderedArgs["expression"]), nil
	},
	"call": func(args []gon.KeyNode) (gon.Node, error) {
		valuer := args[0].Node.(gon.Valued)

		expressionTransform := func(from gon.KeyNode) gon.Node {
			return from.Node
		}

		transformedArgs := sliceutils.Map(args[1:], expressionTransform)

		return gon.Call(valuer.Value().(string), transformedArgs...), nil
	},
	"time": func(args []gon.KeyNode) (gon.Node, error) {
		valuer, ok := args[0].Node.(gon.Valued)
		if !ok {
			return nil, fmt.Errorf("time should be parsed only from string")
		}

		rawTime, ok := valuer.Value().(string)
		if !ok {
			return nil, fmt.Errorf("time should be parsed only from string")
		}

		t, err := time.Parse(time.RFC3339, rawTime)
		if err != nil {
			return nil, fmt.Errorf("time is invalid: %w", err)
		}

		return gon.Literal(t), nil
	},
	"hasPrefix": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "text", "prefix")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'not' node: %w", err)
		}

		return gon.HasPrefix(orderedArgs["text"], orderedArgs["prefix"]), nil
	},
	"hasSuffix": func(args []gon.KeyNode) (gon.Node, error) {
		orderedArgs, _, err := argSorter(args, "text", "prefix")
		if err != nil {
			return nil, fmt.Errorf("error decoding 'not' node: %w", err)
		}

		return gon.HasSuffix(orderedArgs["text"], orderedArgs["prefix"]), nil
	},
	"sum": func(args []gon.KeyNode) (gon.Node, error) {
		_, rest, err := argSorter(args)
		if err != nil {
			return nil, fmt.Errorf("error decoding 'sum' node: %w", err)
		}

		return gon.Sum(rest...), nil
	},
	"avg": func(args []gon.KeyNode) (gon.Node, error) {
		_, rest, err := argSorter(args)
		if err != nil {
			return nil, fmt.Errorf("error decoding 'avg' node: %w", err)
		}

		return gon.Avg(rest...), nil
	},
}
