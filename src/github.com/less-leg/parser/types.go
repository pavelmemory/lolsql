package parser

import (
	"go/ast"
	"go/token"
)

const (
	TableName = "TableName"
)

type ParsedStruct struct {
	Name    string
	Type    *ast.StructType
	Methods map[string]*ast.FuncDecl
}

func (this *ParsedStruct) TableName() string {
	if fdec, found := this.Methods[TableName]; found {
		for _, stmnt := range fdec.Body.List {
			if rstmnt, ok := stmnt.(*ast.ReturnStmt); ok {
				if len(rstmnt.Results) == 1 {
					if blit, ok := rstmnt.Results[0].(*ast.BasicLit); ok {
						tableNameLiteral := blit.Value
						return tableNameLiteral[1:len(tableNameLiteral) - 1]
					}
				}
			}
		}
	}
	return this.Name
}

func (this *ParsedStruct) FieldDefinitions() []FieldDefinition {
	if this.Type.Fields == nil {
		return []FieldDefinition{}
	}

	for _, field := range this.Type.Fields.List {
		switch t := field.Type.(type) {
		case *ast.ArrayType:
			panic(t)
		case *ast.BasicLit:
			switch k := t.Kind {
			case token.IDENT:

			case token.STRING:

			case token.FLOAT:

			case token.INT:

			default:
				panic(k)
			}
		case *ast.StarExpr:

		}
	}

}

type FieldDefinition interface {
	Name() string
}

type BasicFieldDefinition struct {
	name                string
	column              string
	fieldTypeDefinition FieldTypeDefinition
}

func (this *BasicFieldDefinition) Name() string {
	return this.name
}
func (this *BasicFieldDefinition) Type() FieldTypeDefinition {
	return this.fieldTypeDefinition
}
func (this *BasicFieldDefinition) Column() string {
	return this.column
}

type ComplexFieldDefinition struct {
	name string
}

func (this *ComplexFieldDefinition) Name() string {
	return this.name
}

type EmbeddedFieldDefinition struct {
	typeName            string
}

func (this *EmbeddedFieldDefinition) Name() string {
	return this.typeName
}


type FieldTypeDefinition interface {
	Like() bool
	Null() bool
}

type Likable struct {}
func (*Likable) Like() bool {return true}
type NotLikable struct {}
func (*NotLikable) Like() bool {return false}

// String types
type StringFieldTypeDefinition struct {
	Likable
	nullable bool
}
func (this *StringFieldTypeDefinition) Null() bool {
	return this.nullable
}

// Integer types
type IntFieldTypeDefinition struct {
	NotLikable
	nullable bool
}
func (this *IntFieldTypeDefinition) Null() bool {
	return this.nullable
}

// Floating types
type FloatFieldTypeDefinition struct {
	NotLikable
	nullable bool
}
func (this *FloatFieldTypeDefinition) Null() bool {
	return this.nullable
}

//type TimeFieldTypeDefinition struct {
//	NotLikable bool
//}
