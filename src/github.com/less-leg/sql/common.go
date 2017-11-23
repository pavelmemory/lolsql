package sql

import "github.com/less-leg/parser"

type Field interface {
	GetType() parser.TypeIdentity
	GetName() string
}
