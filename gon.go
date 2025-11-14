package gon

import (
	"github.com/sonalys/gon/adapters"
	"github.com/sonalys/gon/encoding"
	"github.com/sonalys/gon/internal/nodes"
)

type SerializableNode interface {
	adapters.Node
	adapters.Named
	adapters.Shaped
	encoding.AutoRegisterer
}

var (
	Avg            = nodes.Avg
	Call           = nodes.Call
	Equal          = nodes.Equal
	Greater        = nodes.Greater
	GreaterOrEqual = nodes.GreaterOrEqual
	HasPrefix      = nodes.HasPrefix
	HasSuffix      = nodes.HasSuffix
	If             = nodes.If
	Literal        = nodes.Literal
	Not            = nodes.Not
	Or             = nodes.Or
	Reference      = nodes.Reference
	Smaller        = nodes.Smaller
	SmallerOrEqual = nodes.SmallerOrEqual
	Sum            = nodes.Sum
	IsEmpty        = nodes.IsEmpty
)
