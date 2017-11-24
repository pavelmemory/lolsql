package sql

import "github.com/less-leg/parser"

type Field interface {
	GetType() parser.TypeIdentity
	GetName() string
}

type SortOrder interface{}

func Desc(Field) SortOrder {
	return nil
}

func Asc(Field) SortOrder {
	return nil
}
