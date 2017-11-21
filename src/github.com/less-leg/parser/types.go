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
	structDefs := map[string]StructDefinition{}
	for _, psdStruct := range parsedStructs {
		structDef := psdStruct.ToStructDefinition(modelPackage)
		structDefs[structDef.Name()] = structDef
	}

	return &PackageDefinition{
		ModelPackage:      modelPackage,
		PackageDirPath:    lolDirPath,
		StructDefinitions: structDefs,
	}
}

func (pkgDef *PackageDefinition) ModelPackageName() string {
	index := strings.LastIndex(pkgDef.ModelPackage, "/")
	if index > -1 {
		return pkgDef.ModelPackage[index+1:]
	}
	return pkgDef.ModelPackage
}

func (pkgDef *PackageDefinition) ColumnNames(structName string) []string {
	var columnNames []string
	if structDef, found := pkgDef.StructDefinitions[structName]; found {
		for _, fieldDef := range structDef.FieldDefinitions() {
			switch typedFieldDef := fieldDef.(type) {
			case *SimpleFieldDefinition:
				columnNames = append(columnNames, typedFieldDef.ColumnName)
			case *ComplexFieldDefinition:
				if typedFieldDef.Embedded {
					columnNames = append(columnNames, pkgDef.ColumnNames(typedFieldDef.Name())...)
				}
			}
		}
	} else {
		panic("Struct " + structName + " was not found among parsed structs")
	}
	return columnNames
}

func (pkgDef *PackageDefinition) FieldsToColumns(structName string) [][]string {
	var fieldsToColumns [][]string
	if structDef, found := pkgDef.StructDefinitions[structName]; found {
		for _, fieldDef := range structDef.FieldDefinitions() {
			switch typedFieldDef := fieldDef.(type) {
			case *SimpleFieldDefinition:
				fieldsToColumns = append(fieldsToColumns, []string{typedFieldDef.Name(), typedFieldDef.ColumnName})
			case *ComplexFieldDefinition:
				if typedFieldDef.Embedded {
					embeddedFieldsToColumns := pkgDef.FieldsToColumns(typedFieldDef.name)
					for _, embeddedFieldToColumn := range embeddedFieldsToColumns {
						embeddedFieldToColumn[0] = typedFieldDef.name + "." + embeddedFieldToColumn[0]
					}
					fieldsToColumns = append(fieldsToColumns, embeddedFieldsToColumns...)
				} else {
					fieldsToColumns = append(fieldsToColumns, []string{typedFieldDef.Name(), typedFieldDef.ColumnName})
				}
			}
		}
	} else {
		panic("Struct " + structName + " was not found among parsed structs")
	}
	return fieldsToColumns
}

type ParsedStruct struct {
	Name    string
	Type    *ast.StructType
	Methods map[string]*ast.FuncDecl
}

func (prsdStruct *ParsedStruct) TableName() string {
	if methodDef, found := prsdStruct.Methods[TableName]; found {
		for _, statement := range methodDef.Body.List {
			if returnStatement, ok := statement.(*ast.ReturnStmt); ok {
				if len(returnStatement.Results) == 1 {
					if basicLiteral, ok := returnStatement.Results[0].(*ast.BasicLit); ok {
						tableNameLiteral := basicLiteral.Value
						return tableNameLiteral[1 : len(tableNameLiteral)-1]
					}
				}
			}
		}
	}
	return prsdStruct.Name
}

func (prsdStruct *ParsedStruct) ToStructDefinition(modelPackage string) StructDefinition {
	if scanMethod, found := prsdStruct.Methods["Scan"]; found {
		return &CustomUserTypeStructDefinition{
			name:  prsdStruct.Name,
			scan:  methodHasSignature(scanMethod, []string{"interface"}, []string{"error"}),
			value: methodHasSignature(prsdStruct.Methods["Value"], nil, []string{"Value", "error"}),
		}
	}

	fields, embeddable := prsdStruct.FieldDefinitions(modelPackage)
	if embeddable {
		return &EmbeddedStructDefinition{
			name:             prsdStruct.Name,
			fieldDefinitions: fields,
		}
	} else {
		tableName := prsdStruct.TableName()
		return &TableStructDefinition{
			name:             prsdStruct.Name,
			TableName:        tableName,
			fieldDefinitions: fields,
		}
	}
}

// TODO: this function doesn't take into account selectors for parameter types and return types
func methodHasSignature(funcDeclaration *ast.FuncDecl, paramTypeNames []string, returnTypeNames []string) bool {
	if funcDeclaration == nil {
		return false
	}
	return checkHasAll(funcDeclaration.Type.Params, paramTypeNames) &&
		checkHasAll(funcDeclaration.Type.Results, returnTypeNames)
}

func checkHasAll(fieldList *ast.FieldList, expectedTypeNames []string) bool {
	if len(expectedTypeNames) > 0 && fieldList != nil && len(fieldList.List) == len(expectedTypeNames) {
		for index, field := range fieldList.List {
			if selector, ok := field.Type.(*ast.SelectorExpr); ok {
				if selector.Sel.Name != expectedTypeNames[index] {
					return false
					//if ident, ok := selector.X.(*ast.Ident); ok {
					//	if ident.Name == "driver" {
					//		resultError := results.List[1]
					//		if ident, ok := resultError.UserDefinedType.(*ast.Ident); ok && ident.Name == "error" {
					//
					//		}
					//	}
					//}
				}
			} else if identifier, ok := field.Type.(*ast.Ident); ok {
				if identifier.Name != expectedTypeNames[index] {
					return false
				}
			}
		}
	}
	return true
}

func (prsdStruct *ParsedStruct) FieldDefinitions(modelPackage string) ([]FieldDefinition, bool) {
	var fieldDefs []FieldDefinition
	if prsdStruct.Type.Fields == nil {
		return fieldDefs, false
	}

	embeddable := true
	packageName := string(modelPackage[strings.LastIndex(modelPackage, "/")+1:])

	for _, field := range prsdStruct.Type.Fields.List {
		fieldTypeDef := fieldTypeDefinition(field.Type, packageName)
		if len(field.Names) == 0 {
			if fieldTypeDef.IsBasic() {
				fieldDefs = append(fieldDefs, &ComplexFieldDefinition{
					name:      fieldTypeDef.Name(),
					Embedded:  true,
					fieldType: fieldTypeDef,
				})
			} else {
				fieldDef := &ComplexFieldDefinition{
					name:      fieldTypeDef.RawName(),
					Embedded:  true,
					fieldType: fieldTypeDef,
				}
				fieldDefs = append(fieldDefs, fieldDef)
				fmt.Printf("Gotcha! %s\n", fieldTypeDef.String())
			}
		} else {
			tagConf := NewTagConfig(field)
			embeddable = embeddable && !tagConf.Primary
			var fieldDef FieldDefinition
			if fieldTypeDef.IsBasic() {
				fieldDef = &SimpleFieldDefinition{
					name:       field.Names[0].Name,
					Primary:    tagConf.Primary,
					ColumnName: tagConf.ColumnName,
					fieldType:  fieldTypeDef,
				}
			} else {
				fieldDef = &ComplexFieldDefinition{
					name:       fieldName(field),
					fieldType:  fieldTypeDef,
					ColumnName: tagConf.ColumnName,
					// TODO package name
				}
			}
			fieldDefs = append(fieldDefs, fieldDef)
		}
	}
	return fieldDefs, embeddable
}

func fieldTypeDefinition(expr ast.Expr, packageName string) *FieldTypeDefinition {
	switch typedExpression := expr.(type) {
	case *ast.ArrayType:
		return &FieldTypeDefinition{Slice: true, Underlying: fieldTypeDefinition(typedExpression.Elt, packageName)}
	case *ast.Ident:
		if SupportedBasicTypes[typedExpression.Name] {
			return &FieldTypeDefinition{name: typedExpression.Name, likable: typedExpression.Name == "string"}
		} else {
			return &FieldTypeDefinition{
				selector:   packageName,
				Underlying: &FieldTypeDefinition{name: typedExpression.Name, likable: typedExpression.Name == "string"}}
		}
	case *ast.SelectorExpr:
		return &FieldTypeDefinition{selector: typedExpression.X.(*ast.Ident).Name,
			Underlying: &FieldTypeDefinition{name: typedExpression.Sel.Name}}
	case *ast.StarExpr:
		return &FieldTypeDefinition{Ptr: true, Underlying: fieldTypeDefinition(typedExpression.X, packageName)}
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

func (sfd *SimpleFieldDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	fm := sfd.fieldType.FetchMeta()
	fm.FieldName = sfd.Name()
	return []*FetchMeta{fm}
}

func (sfd *SimpleFieldDefinition) Name() string {
	return sfd.name
}

func (sfd *SimpleFieldDefinition) FieldType() *FieldTypeDefinition {
	return sfd.fieldType
}

func (sfd *SimpleFieldDefinition) String() string {
	return fmt.Sprintf("Field[%s] Column[%s] Primary[%s] UserDefinedType[%s]",
		sfd.name, sfd.ColumnName, strconv.FormatBool(sfd.Primary), sfd.fieldType.String())
}

type ComplexFieldDefinition struct {
	name       string
	ColumnName string
	Embedded   bool
	fieldType  *FieldTypeDefinition
}

func (cfd *ComplexFieldDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	if cfd.Embedded {
		var fms []*FetchMeta
		fmt.Println("FetchMeta():", cfd.name, cfd.Embedded, cfd.ColumnName, cfd.fieldType)
		if sd, found := pckgDef.StructDefinitions[cfd.Name()]; found {
			for _, fieldDef := range sd.FieldDefinitions() {
				metas := fieldDef.FetchMeta(pckgDef)
				for i := 0; i < len(metas); i++ {
					metas[i].FieldName = cfd.Name() + "." + metas[i].FieldName
					fms = append(fms, metas[i])
				}
			}
		}
		return fms
	} else {
		fm := cfd.fieldType.FetchMeta()
		fm.FieldName = cfd.Name()
		return []*FetchMeta{fm}
	}
}

func (cfd *ComplexFieldDefinition) Name() string {
	return cfd.name
}

func (cfd *ComplexFieldDefinition) String() string {
	return fmt.Sprintf("Complex Field[%s] UserDefinedType[%s] Embedded[%s]",
		cfd.name, cfd.fieldType.String(), strconv.FormatBool(cfd.Embedded))
}

func (cfd *ComplexFieldDefinition) FieldType() *FieldTypeDefinition {
	return cfd.fieldType
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

func (tsd *TableStructDefinition) FetchMeta(pckgDef *PackageDefinition) []*FetchMeta {
	fetchMetas := make([]*FetchMeta, 0, len(tsd.FieldDefinitions()))
	for _, fieldDef := range tsd.FieldDefinitions() {
		fetchMetas = append(fetchMetas, fieldDef.FetchMeta(pckgDef)...)
	}
	return fetchMetas
}

func (tsd *TableStructDefinition) String() string {
	var fieldDefStrings []string
	for _, fd := range tsd.fieldDefinitions {
		fieldDefStrings = append(fieldDefStrings, fd.String())
	}
	return fmt.Sprintf("TableStructDefinition[%s] Table[%s] Fields[\n\t%s]",
		tsd.name, tsd.TableName, strings.Join(fieldDefStrings, "\n\t"))
}

func (tsd *TableStructDefinition) Name() string {
	return tsd.name
}

func (tsd *TableStructDefinition) Type() *FieldTypeDefinition {
	panic("Struct has no type")
}

func (tsd *TableStructDefinition) FieldDefinitions() []FieldDefinition {
	return tsd.fieldDefinitions
}

func (tsd *TableStructDefinition) Selectors(except string) []string {
	var selectors []string
	for _, fieldDef := range tsd.fieldDefinitions {
		if sel := fieldDef.FieldType().Selector(); sel != "" && sel != except {
			selectors = append(selectors, sel)
		}
	}
	return selectors
}

type EmbeddedStructDefinition struct {
	name             string
	fieldDefinitions []FieldDefinition
}

func (esd *EmbeddedStructDefinition) String() string {
	var fieldDefStrings []string
	for _, fd := range esd.fieldDefinitions {
		fieldDefStrings = append(fieldDefStrings, fd.String())
	}
	return fmt.Sprintf("EmbeddedStructDefinition[%s] Fields[\n\t%s]", esd.name, strings.Join(fieldDefStrings, "\n\t"))
}

func (esd *EmbeddedStructDefinition) Name() string {
	return esd.name
}

func (esd *EmbeddedStructDefinition) Type() *FieldTypeDefinition {
	panic("Struct has no type")
}

func (esd *EmbeddedStructDefinition) FieldDefinitions() []FieldDefinition {
	return esd.fieldDefinitions
}

func (esd *EmbeddedStructDefinition) Selectors(except string) []string {
	var selectors []string
	for _, fieldDef := range esd.fieldDefinitions {
		selectors = append(selectors, fieldDef.FieldType().Selector())
	}
	return selectors
}

type CustomUserTypeStructDefinition struct {
	name             string
	scan             bool
	value            bool
	selectors        []string
	fieldDefinitions []FieldDefinition
}

func (cutsd *CustomUserTypeStructDefinition) FieldDefinitions() []FieldDefinition {
	return cutsd.fieldDefinitions
}

func (cutsd *CustomUserTypeStructDefinition) Name() string {
	return cutsd.name
}

func (cutsd *CustomUserTypeStructDefinition) Selectors(except string) []string {
	return cutsd.selectors
}

func (cutsd *CustomUserTypeStructDefinition) String() string {
	var fieldDefStings []string
	for _, fd := range cutsd.fieldDefinitions {
		fieldDefStings = append(fieldDefStings, fd.String())
	}
	return fmt.Sprintf("CustomUserTypeStructDefinition[%s] Fields[\n\t%s]", cutsd.name, strings.Join(fieldDefStings, "\n\t"))
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

func (ftd *FieldTypeDefinition) IsLikable() bool {
	if ftd.Underlying != nil {
		return ftd.Underlying.IsLikable()
	}
	return ftd.likable
}

func (ftd *FieldTypeDefinition) FetchMeta() *FetchMeta {
	fm := &FetchMeta{FieldType: ftd.Name()}
	if ftd.IsNullable() {
		fm.IsNullable = "*"
	}
	return fm
}

func (ftd *FieldTypeDefinition) IsNullable() bool {
	if ftd.Ptr || ftd.Slice {
		return true
	}
	return false
}

func (ftd *FieldTypeDefinition) PtrSign() string {
	if ftd.Ptr || ftd.Slice {
		return "*"
	}
	return ""
}

func (ftd *FieldTypeDefinition) IsBasic() bool {
	return SupportedBasicTypes[ftd.Name()]
}

func (ftd *FieldTypeDefinition) String() string {
	switch {
	case ftd.Ptr:
		return "*" + ftd.Underlying.String()
	case ftd.Slice:
		return "[]" + ftd.Underlying.String()
	case ftd.selector != "":
		return ftd.selector + "." + ftd.Underlying.String()
	case ftd.Underlying == nil:
		return ftd.name
	default:
		panic("Unreachable code")
	}
}

func (ftd *FieldTypeDefinition) Name() string {
	switch {
	case ftd.Ptr:
		return ftd.Underlying.Name()
	case ftd.Slice:
		return ftd.Underlying.Name()
	case ftd.selector != "":
		return ftd.selector + "." + ftd.Underlying.Name()
	case ftd.Underlying == nil:
		return ftd.name
	default:
		panic("Unreachable code")
	}
}

func (ftd *FieldTypeDefinition) RawName() string {
	switch {
	case ftd.Ptr:
		return ftd.Underlying.Name()
	case ftd.Slice:
		return ftd.Underlying.Name()
	case ftd.selector != "":
		return ftd.Underlying.Name()
	case ftd.Underlying == nil:
		return ftd.name
	default:
		panic("Unreachable code")
	}
}

func (ftd *FieldTypeDefinition) Selector() string {
	switch {
	case ftd.Ptr:
		return ftd.Underlying.Selector()
	case ftd.Slice:
		return ftd.Underlying.Selector()
	case ftd.selector != "":
		return ftd.selector
	default:
		return ""
	}
}
