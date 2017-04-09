package parser

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
)

const (
	TableName = "TableName"
)

var SupportedBasicTypes = map[string]bool{
	"int":   true,
	"int8":  true,
	"int16": true,
	"int32": true,
	"int64": true,

	"uint":   true,
	"uint8":  true,
	"uint16": true,
	"uint32": true,
	"uint64": true,

	"float32": true,
	"float64": true,

	"string": true,

	"time.Time": true,

	"bool": true,
}

type PackageDefinition struct {
	ModelPackage      string
	PackageDirPath    string
	StructDefinitions map[string]StructDefinition
}

func NewPackageDefinition(modelPackage string, lolDirPath string, parsedStructs []*ParsedStruct) *PackageDefinition {
	structDefinitions := map[string]StructDefinition{}
	for _, psdStruct := range parsedStructs {
		sdef := psdStruct.ToStructDefinition(modelPackage)
		structDefinitions[sdef.Name()] = sdef
	}

	return &PackageDefinition{
		ModelPackage:      modelPackage,
		PackageDirPath:    lolDirPath,
		StructDefinitions: structDefinitions,
	}
}

func (this *PackageDefinition) ModelPackageName() string {
	index := strings.LastIndex(this.ModelPackage, "/")
	if index > -1 {
		return this.ModelPackage[index+1:]
	}
	return this.ModelPackage
}

func (this *PackageDefinition) ColumnNames(structName string) []string {
	var cnames []string
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
	var ftoc [][]string
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
				} else {
					ftoc = append(ftoc, []string{sfdef.Name(), sfdef.ColumnName})
				}
			}
		}
	} else {
		panic("Struct " + structName + " was not found among parsed structs")
	}
	return ftoc
}

func (this *PackageDefinition) FieldToColumn(structName string) [][]string {
	var ftoc [][]string
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
						return tableNameLiteral[1 : len(tableNameLiteral)-1]
					}
				}
			}
		}
	}
	return this.Name
}

func (this *ParsedStruct) ToStructDefinition(modelPackage string) StructDefinition {
	if scanMethod, found := this.Methods["Scan"]; found {
		return &CustomUserTypeStructDefinition{
			name:  this.Name,
			scan:  methodHasSignature(scanMethod, []string{"interface"}, []string{"error"}),
			value: methodHasSignature(this.Methods["Value"], nil, []string{"Value", "error"}),
		}
	}

	fields, embeddable := this.FieldDefinitions(modelPackage)
	if embeddable {
		return &EmbeddedStructDefinition{
			name:             this.Name,
			fieldDefinitions: fields,
		}
	} else {
		tableName := this.TableName()
		return &TableStructDefinition{
			name:             this.Name,
			TableName:        tableName,
			fieldDefinitions: fields,
		}
	}
}

// TODO: this function doesn't take into account selectors for parameter types and return types
func methodHasSignature(fdecl *ast.FuncDecl, paramTypeNames []string, returnTypeNames []string) bool {
	if fdecl == nil {
		return false
	}
	return checkHasAll(fdecl.Type.Params, paramTypeNames) && checkHasAll(fdecl.Type.Results, returnTypeNames)
}

func checkHasAll(fieldList *ast.FieldList, expectedTypeNames []string) bool {
	if len(expectedTypeNames) > 0 && fieldList != nil && len(fieldList.List) == len(expectedTypeNames) {
		for indx, field := range fieldList.List {
			if selector, ok := field.Type.(*ast.SelectorExpr); ok {
				if selector.Sel.Name != expectedTypeNames[indx] {
					return false
					//if ident, ok := selector.X.(*ast.Ident); ok {
					//	if ident.Name == "driver" {
					//		resultError := results.List[1]
					//		if ident, ok := resultError.Type.(*ast.Ident); ok && ident.Name == "error" {
					//
					//		}
					//	}
					//}
				}
			} else if ident, ok := field.Type.(*ast.Ident); ok {
				if ident.Name != expectedTypeNames[indx] {
					return false
				}
			}
		}
	}
	return true
}

func (this *ParsedStruct) FieldDefinitions(modelPackage string) ([]FieldDefinition, bool) {
	var fdefs []FieldDefinition
	if this.Type.Fields == nil {
		return fdefs, false
	}

	embeddable := true
	packageName := string(modelPackage[strings.LastIndex(modelPackage, "/")+1:])

	for _, field := range this.Type.Fields.List {
		ftypedef := fieldTypeDefinition(field.Type, packageName)
		if len(field.Names) == 0 {
			if ftypedef.IsBasic() {
				fdefs = append(fdefs, &ComplexFieldDefinition{
					name:      ftypedef.Name(),
					Embedded:  true,
					fieldType: ftypedef,
				})
			} else {
				fdef := &ComplexFieldDefinition{
					name:      ftypedef.RawName(),
					Embedded:  true,
					fieldType: ftypedef,
				}
				fdefs = append(fdefs, fdef)
				fmt.Printf("Gotcha! %s\n", ftypedef.String())
			}
		} else {
			tconf := NewTagConfig(field)
			embeddable = embeddable && !tconf.Primary
			var fdef FieldDefinition
			if ftypedef.IsBasic() {
				fdef = &SimpleFieldDefinition{
					name:       field.Names[0].Name,
					Primary:    tconf.Primary,
					ColumnName: tconf.ColumnName,
					fieldType:  ftypedef,
				}
			} else {
				fdef = &ComplexFieldDefinition{
					name:       fieldName(field),
					fieldType:  ftypedef,
					ColumnName: tconf.ColumnName,
					// TODO package name
				}
			}
			fdefs = append(fdefs, fdef)
		}
	}
	return fdefs, embeddable
}

func fieldTypeDefinition(expr ast.Expr, packageName string) *FieldTypeDefinition {
	switch expt := expr.(type) {
	case *ast.ArrayType:
		return &FieldTypeDefinition{Slice: true, Underlying: fieldTypeDefinition(expt.Elt, packageName)}
	case *ast.Ident:
		if SupportedBasicTypes[expt.Name] {
			return &FieldTypeDefinition{name: expt.Name, likable: expt.Name == "string"}
		} else {
			return &FieldTypeDefinition{selector: packageName, Underlying: &FieldTypeDefinition{name: expt.Name, likable: expt.Name == "string"}}
		}
	case *ast.SelectorExpr:
		return &FieldTypeDefinition{selector: expt.X.(*ast.Ident).Name, Underlying: &FieldTypeDefinition{name: expt.Sel.Name}}
	case *ast.StarExpr:
		return &FieldTypeDefinition{Ptr: true, Underlying: fieldTypeDefinition(expt.X, packageName)}
	default:
		panic(fmt.Sprintf("Not supported relation: %#v ", expr))
	}
}

type FieldDefinition interface {
	fmt.Stringer
	Name() string
	FieldType() *FieldTypeDefinition
	FetchMeta(pckgDef *PackageDefinition) []*FetchMeta
}

type SimpleFieldDefinition struct {
	name       string
	Primary    bool
	ColumnName string
	fieldType  *FieldTypeDefinition
}

func (this *SimpleFieldDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	fm := this.fieldType.FetchMeta()
	fm.FieldName = this.Name()
	return []*FetchMeta{fm}
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
	name       string
	ColumnName string
	Embedded   bool
	fieldType  *FieldTypeDefinition
}

func (this *ComplexFieldDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	if this.Embedded {
		var fms []*FetchMeta
		fmt.Println("FetchMeta():", this.name, this.Embedded, this.ColumnName, this.fieldType)
		if sd, found := pckgDef.StructDefinitions[this.Name()]; found {
			for _, fdef := range sd.FieldDefinitions() {
				metas := fdef.FetchMeta(pckgDef)
				for i := 0; i < len(metas); i++ {
					metas[i].FieldName = this.Name() + "." + metas[i].FieldName
					fms = append(fms, metas[i])
				}
			}
		}
		return fms
	} else {
		fm := this.fieldType.FetchMeta()
		fm.FieldName = this.Name()
		return []*FetchMeta{fm}
	}
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

type FetchMeta struct {
	FieldName  string
	FieldType  string
	IsNullable string
}

type StructDefinition interface {
	fmt.Stringer
	Name() string
	FieldDefinitions() []FieldDefinition
	Selectors(except string) []string
}

type TableStructDefinition struct {
	name             string
	TableName        string
	fieldDefinitions []FieldDefinition
}

func (this *TableStructDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	fmets := make([]*FetchMeta, 0, len(this.FieldDefinitions()))
	for _, fdef := range this.FieldDefinitions() {
		fmets = append(fmets, fdef.FetchMeta(pckgDef)...)
	}
	return fmets
}

func (this *TableStructDefinition) String() string {
	var fdstr []string
	for _, fd := range this.fieldDefinitions {
		fdstr = append(fdstr, fd.String())
	}
	return fmt.Sprintf("TableStructDefinition[%s] Table[%s] Fields[\n\t%s]", this.name, this.TableName, strings.Join(fdstr, "\n\t"))
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

func (this *TableStructDefinition) Selectors(except string) []string {
	var slcrs []string
	for _, fdef := range this.fieldDefinitions {
		if sel := fdef.FieldType().Selector(); sel != "" && sel != except {
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
	var fdstr []string
	for _, fd := range this.fieldDefinitions {
		fdstr = append(fdstr, fd.String())
	}
	return fmt.Sprintf("EmbeddedStructDefinition[%s] Fields[\n\t%s]", this.name, strings.Join(fdstr, "\n\t"))
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

func (this *EmbeddedStructDefinition) Selectors(except string) []string {
	var slcrs []string
	for _, fdef := range this.fieldDefinitions {
		slcrs = append(slcrs, fdef.FieldType().Selector())
	}
	return slcrs
}

type CustomUserTypeStructDefinition struct {
	name             string
	scan             bool
	value            bool
	selectors        []string
	fieldDefinitions []FieldDefinition
}

func (this *CustomUserTypeStructDefinition) FieldDefinitions() []FieldDefinition {
	return this.fieldDefinitions
}

func (this *CustomUserTypeStructDefinition) Name() string {
	return this.name
}

func (this *CustomUserTypeStructDefinition) Selectors(except string) []string {
	return this.selectors
}

func (this *CustomUserTypeStructDefinition) String() string {
	var fdstr []string
	for _, fd := range this.fieldDefinitions {
		fdstr = append(fdstr, fd.String())
	}
	return fmt.Sprintf("CustomUserTypeStructDefinition[%s] Fields[\n\t%s]", this.name, strings.Join(fdstr, "\n\t"))
}

type TagConfig struct {
	Primary    bool
	ColumnName string
}

func NewTagConfig(field *ast.Field) *TagConfig {
	tag := fieldTag(field)
	return &TagConfig{
		Primary:    isPrimary(tag),
		ColumnName: columnName(field, tag),
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

func columnName(field *ast.Field, tag string) string {
	colName := field.Names[0].Name
	colTagStart := strings.Index(tag, "column[")
	if colTagStart >= 0 {
		colNameStart := colTagStart + len("column[")
		colNameEnd := strings.Index(tag[colNameStart:], "]")
		if colNameEnd < 1 {
			panic("Columnt name must be at least 1 character long: " + tag)
		}
		colName = string(tag[colNameStart : colNameStart+colNameEnd])
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
	name       string
	Ptr        bool
	Slice      bool
	selector   string
	likable    bool
	Underlying *FieldTypeDefinition
}

func (this *FieldTypeDefinition) IsLikable() bool {
	if this.Underlying != nil {
		return this.Underlying.IsLikable()
	}
	return this.likable
}

func (this *FieldTypeDefinition) FetchMeta() *FetchMeta {
	fm := &FetchMeta{FieldType: this.Name()}
	if this.IsNullable() {
		fm.IsNullable = "*"
	}
	return fm
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
	switch {
	case this.Ptr:
		return "*" + this.Underlying.String()
	case this.Slice:
		return "[]" + this.Underlying.String()
	case this.selector != "":
		return this.selector + "." + this.Underlying.String()
	case this.Underlying == nil:
		return this.name
	default:
		panic("Unreachable code")
	}
}

func (this *FieldTypeDefinition) Name() string {
	switch {
	case this.Ptr:
		return this.Underlying.Name()
	case this.Slice:
		return this.Underlying.Name()
	case this.selector != "":
		return this.selector + "." + this.Underlying.Name()
	case this.Underlying == nil:
		return this.name
	default:
		panic("Unreachable code")
	}
}

func (this *FieldTypeDefinition) RawName() string {
	switch {
	case this.Ptr:
		return this.Underlying.Name()
	case this.Slice:
		return this.Underlying.Name()
	case this.selector != "":
		return this.Underlying.Name()
	case this.Underlying == nil:
		return this.name
	default:
		panic("Unreachable code")
	}
}

func (this *FieldTypeDefinition) Selector() string {
	switch {
	case this.Ptr:
		return this.Underlying.Selector()
	case this.Slice:
		return this.Underlying.Selector()
	case this.selector != "":
		return this.selector
	default:
		return ""
	}
}
