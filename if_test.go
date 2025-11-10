package gon_test

import (
	"github.com/sonalys/gon"
	"github.com/sonalys/gon/ast"
)

var (
	_ ast.ParseableNode = &gon.IfNode{}
)
