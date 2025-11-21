package nodes_test

import (
	"testing"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/internal/nodes"
	"github.com/stretchr/testify/assert"
)

func Test_Banner(t *testing.T) {
	shapedList := []interface {
		adapters.Shaped
		adapters.Named
	}{
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
	}

	for _, shaped := range shapedList {
		t.Run(shaped.Scalar(), func(t *testing.T) {
			assert.NotPanics(t, func() {
				_ = shaped.Shape()
			})
		})
	}
}

func Benchmark_Equal(b *testing.B) {
	scope, _ := gon.NewScope().
		WithContext(b.Context()).
		WithValues(gon.Values{
			"var1": nodes.Literal(1),
			"var2": nodes.Literal(1),
		})

	isEqual := nodes.Equal(
		nodes.Reference("var1"),
		nodes.Reference("var2"),
	)

	for b.Loop() {
		isEqual.Eval(scope)
	}
}
