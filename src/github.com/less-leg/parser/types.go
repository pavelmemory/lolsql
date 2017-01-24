package parser

import (
	"go/ast"
	"strings"
	"fmt"
	"strconv"
)

const (
	TableName = "TableName"
)

var SupportedBasicTypes = map[string]bool{
	"int": true,
	"int8": true,
	"int16": true,
	"int32": true,
	"int64": true,

	"uint": true,
	"uint8": true,
	"uint16": true,
	"uint32": true,
	"uint64": true,

	"float32": true,
	"float64": true,

	"string": true,

	"time.Time": true,
}

type PackageDefinition struct {
	StructDefinitions map[string]StructDefinition
	PackageDirPath string
}

func NewPackageDefinition(lolDirPath string, parsedStructs []*ParsedStruct) *PackageDefinition {
	structDefinitions := map[string]StructDefinition{}
	for _, psdStruct := range parsedStructs {
		strDef := psdStruct.ToStructDefinition()
		structDefinitions[strDef.Name()] = strDef
		fmt.Printf("%s\n", strDef.String())
	}

	return &PackageDefinition{
		PackageDirPath:lolDirPath,
		StructDefinitions:structDefinitions,
	}
}

func (this *PackageDefinition) ColumnNames(structName string) []string {
	cnames := []string{}
	if sdef, found := this.StructDefinitions[structName]; found {
		for _, fdef := range sdef.FieldDefinitions() {
			switch sfdef := fdef.(type) {
			case *SimpleFieldDefinition:
				cnames = append(cnames, sfdef.ColumnName)
			case *ComplexFieldDefinition:
				if sfdef.Embedded {
					cnames = append(cnames, this.ColumnNames(sfdef.Name())...)
				}
			}
		}
	} else {
		panic("Struct " + structName + " was not found among parsed structs")
	}
	return cnames
}

func (this *PackageDefinition) FieldsToColumns(structName string) [][]string {
	ftoc := [][]string{}
	if sdef, found := this.StructDefinitions[structName]; found {
		for _, fdef := range sdef.FieldDefinitions() {
			switch sfdef := fdef.(type) {
			case *SimpleFieldDefinition:
				ftoc = append(ftoc, []string{sfdef.Name(), sfdef.ColumnName})
			case *ComplexFieldDefinition:
				if sfdef.Embedded {
					embftocs := this.FieldsToColumns(sfdef.name)
					for _, embftoc := range embftocs {
						embftoc[0] = sfdef.name + "." + embftoc[0]
					}
					ftoc = append(ftoc, embftocs...)
				}
			}
		}
	} else {
		panic("Struct " + structName + " was not found among parsed structs")
	}
	return ftoc
}

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

func (this *ParsedStruct) ToStructDefinition() StructDefinition {
	fields, embeddable := this.FieldDefinitions()
	if embeddable {
		return &EmbeddedStructDefinition{name:this.Name, fieldDefinitions:fields}
	} else {
		tableName := this.TableName()
		return &TableStructDefinition{name:this.Name, TableName:tableName, fieldDefinitions:fields}
	}
}

func (this *ParsedStruct) FieldDefinitions() ([]FieldDefinition, bool) {
	fdefs := []FieldDefinition{}
	if this.Type.Fields == nil {
		return fdefs, false
	}

	embeddable := true

	for _, field := range this.Type.Fields.List {
		ftypedef := fieldTypeDefinition(field.Type)
		if len(field.Names) == 0 {
			if ftypedef.IsBasic() {
				panic("Embedded basic field")
			}
			fdefs = append(fdefs, &ComplexFieldDefinition{name:ftypedef.Name(), Embedded:true, fieldType:ftypedef})
		} else {
			if ftypedef.IsBasic() {
				tconf := NewTagConfig(field)
				sdef := &SimpleFieldDefinition{
					name:field.Names[0].Name,
					Primary:tconf.Primary,
					ColumnName:tconf.ColumnName,
					fieldType:ftypedef,
				}
				fdefs = append(fdefs, sdef)
				embeddable = embeddable && !tconf.Primary
			} else {
				fdefs = append(fdefs, &ComplexFieldDefinition{name:fieldName(field), fieldType:ftypedef})
			}
		}
	}
	return fdefs, embeddable
}

func fieldTypeDefinition(expr ast.Expr) *FieldTypeDefinition {
	switch expt := expr.(type) {
	case *ast.ArrayType:
		return &FieldTypeDefinition{Slice:true, Underlying:fieldTypeDefinition(expt.Elt)}
	case *ast.Ident:
		return &FieldTypeDefinition{name:expt.Name}
	case *ast.SelectorExpr:
		return &FieldTypeDefinition{selector:expt.X.(*ast.Ident).Name, Underlying:&FieldTypeDefinition{name:expt.Sel.Name}}
	case *ast.StarExpr:
		return &FieldTypeDefinition{Ptr:true, Underlying:fieldTypeDefinition(expt.X)}
	default:
		panic(fmt.Sprintf("Not supported relation: %#v ", expr))
	}
}

type FieldDefinition interface {
	fmt.Stringer
	Name() string
	FieldType() *FieldTypeDefinition
}

type SimpleFieldDefinition struct {
	name       string
	Primary    bool
	ColumnName string
	fieldType  *FieldTypeDefinition
}

func (this *SimpleFieldDefinition) Name() string {
	return this.name
}

func (this *SimpleFieldDefinition) FieldType() *FieldTypeDefinition {
	return this.fieldType
}

func (this *SimpleFieldDefinition) String() string {
	return fmt.Sprintf("Field[%s] Column[%s] Primary[%s] Type[%s]",
		this.name, this.ColumnName, strconv.FormatBool(this.Primary), this.fieldType.String())
}

type ComplexFieldDefinition struct {
	fmt.Stringer
	name      string
	Embedded  bool
	fieldType *FieldTypeDefinition
}

func (this *ComplexFieldDefinition) Name() string {
	return this.name
}

func (this *ComplexFieldDefinition) String() string {
	return fmt.Sprintf("Complex Field[%s] Type[%s] Embedded[%s]",
		this.name, this.fieldType.String(), strconv.FormatBool(this.Embedded))
}

func (this *ComplexFieldDefinition) FieldType() *FieldTypeDefinition {
	return this.fieldType
}

type StructDefinition interface {
	fmt.Stringer
	Name() string
	FieldDefinitions() []FieldDefinition
	Selectors() []string
}

type TableStructDefinition struct {
	name             string
	TableName        string
	fieldDefinitions []FieldDefinition
}

func (this *TableStructDefinition) String() string {
	fdstr := []string{}
	for _, fd := range this.fieldDefinitions {
		fdstr = append(fdstr, fd.String())
	}
	return fmt.Sprintf("Struct[%s] Table[%s] Fields[\n\t%s]", this.name, this.TableName, strings.Join(fdstr, "\n\t"))
}

func (this *TableStructDefinition) Name() string {
	return this.name
}

func (this *TableStructDefinition) Type() *FieldTypeDefinition {
	panic("Struct has no type")
}

func (this *TableStructDefinition) FieldDefinitions() []FieldDefinition {
	return this.fieldDefinitions
}

func (this *TableStructDefinition) Selectors() []string {
	slcrs := []string{}
	for _, fdef := range this.fieldDefinitions {
		if sel := fdef.FieldType().Selector(); sel != "" {
			slcrs = append(slcrs, sel)
		}
	}
	return slcrs
}

type EmbeddedStructDefinition struct {
	name             string
	fieldDefinitions []FieldDefinition
}

func (this *EmbeddedStructDefinition) String() string {
	fdstr := []string{}
	for _, fd := range this.fieldDefinitions {
		fdstr = append(fdstr, fd.String())
	}
	return fmt.Sprintf("EmbeddedStruct[%s] Fields[\n\t%s]", this.name, strings.Join(fdstr, "\n\t"))
}

func (this *EmbeddedStructDefinition) Name() string {
	return this.name
}

func (this *EmbeddedStructDefinition) Type() *FieldTypeDefinition {
	panic("Struct has no type")
}

func (this *EmbeddedStructDefinition) FieldDefinitions() []FieldDefinition {
	return this.fieldDefinitions
}

func (this *EmbeddedStructDefinition) Selectors() []string {
	slcrs := []string{}
	for _, fdef := range this.fieldDefinitions {
		slcrs = append(slcrs, fdef.FieldType().Selector())
	}
	return slcrs
}

type TagConfig struct {
	Primary    bool
	ColumnName string
}

func NewTagConfig(field *ast.Field) *TagConfig {
	tag := fieldTag(field)
	return &TagConfig{
		Primary: isPrimary(tag),
		ColumnName:columnName(field, tag),
	}
}

func fieldTag(field *ast.Field) string {
	if field.Tag != nil {
		tagStart := strings.Index(field.Tag.Value, "lolsql")
		if tagStart > -1 {
			tagStart = tagStart + len("lolsql:\"")
			tag := field.Tag.Value[tagStart:]
			tagEnd := strings.Index(tag, "\"")
			return string(tag[:tagEnd])
		}
		return ""
	}
	return ""
}

func isPrimary(tag string) bool {
	return strings.Index(tag, "id[true]") >= 0
}

//const MaxColumnNameLength = 32
//var ColumnReg, _ = regexp.Compile(`.*(column\[.{1,` + strconv.Itoa(MaxColumnNameLength) + `}\]).*`)

func columnName(field *ast.Field, tag string) string {
	//regColName := ColumnReg.FindStringSubmatch(tag)
	//if len(regColName) == 2 {
	//	return regColName[1]
	//}
	//return field.Names[0].Name

	colName := field.Names[0].Name
	colTagStart := strings.Index(tag, "column[")
	if colTagStart >= 0 {
		colNameStart := colTagStart + len("column[")
		colNameEnd := strings.Index(tag[colNameStart:], "]")
		if colNameEnd < 1 {
			panic("Columnt name must be at least 1 character long: " + tag)
		}
		colName = string(tag[colNameStart:colNameStart + colNameEnd])
	}
	return colName
}

func fieldName(field *ast.Field) string {
	if len(field.Names) == 0 {
		panic(field)
	}
	return field.Names[0].Name
}

type FieldTypeDefinition struct {
	fmt.Stringer
	name       string
	Ptr        bool
	Slice      bool
	selector   string
	Underlying *FieldTypeDefinition
}

func (this *FieldTypeDefinition) IsNullable() bool {
	if this.Ptr || this.Slice {
		return true
	}
	return false
}

func (this *FieldTypeDefinition) PtrSign() string {
	if this.Ptr || this.Slice {
		return "*"
	}
	return ""
}

func (this *FieldTypeDefinition) IsBasic() bool {
	return SupportedBasicTypes[this.Name()]
}

func (this *FieldTypeDefinition) String() string {
	if this.Ptr {
		return "*" + this.Underlying.String()
	}
	if this.Slice {
		return "[]" + this.Underlying.String()
	}
	if this.selector != "" {
		return this.selector + "." + this.Underlying.String()
	}
	if this.Underlying == nil {
		return this.name
	}
	panic("Unreachable code")
}

func (this *FieldTypeDefinition) Name() string {
	if this.Ptr {
		return this.Underlying.Name()
	}
	if this.Slice {
		return this.Underlying.Name()
	}
	if this.selector != "" {
		return this.selector + "." + this.Underlying.Name()
	}
	if this.Underlying == nil {
		return this.name
	}
	panic("Unreachable code")
}

func (this *FieldTypeDefinition) Selector() string {
	if this.Ptr {
		return this.Underlying.Selector()
	}
	if this.Slice {
		return this.Underlying.Selector()
	}
	if this.selector != "" {
		return this.selector
	}
	return ""
}