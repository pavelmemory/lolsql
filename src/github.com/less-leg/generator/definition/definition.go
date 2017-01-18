package definition
//
//import (
//	"go/ast"
//	"fmt"
//	"strings"
//	"go/build"
//	"path/filepath"
//	"go/token"
//	"go/parser"
//	"reflect"
//	"errors"
//)
//
//type ColumnType int
//
//const (
//	INT ColumnType = 1 << (iota + 1)
//	FLOAT
//	STRING
//	DATE_TIME
//)
//
//var SupportedTypes = map[string]ColumnType {
//	reflect.Int.String(): INT,
//	reflect.Int8.String(): INT,
//	reflect.Int16.String(): INT,
//	reflect.Int32.String(): INT,
//	reflect.Int64.String(): INT,
//
//	reflect.Uint.String(): INT,
//	reflect.Uint8.String(): INT,
//	reflect.Uint16.String(): INT,
//	reflect.Uint32.String(): INT,
//	reflect.Uint64.String(): INT,
//
//	reflect.String.String(): STRING,
//
//	reflect.Float32.String(): FLOAT,
//	reflect.Float64.String(): FLOAT,
//
//	"time": DATE_TIME,
//	"Time": DATE_TIME,
//	"Date": DATE_TIME,
//}
//
//type LolColumnDefinition struct {
//	Name   string
//	IsNull bool
//	Type   ColumnType
//}
//
//func NewLolColumnDefinition(name string, isNull bool, typ ColumnType) *LolColumnDefinition {
//	return &LolColumnDefinition{
//		Name: name,
//		IsNull:isNull,
//		Type: typ,
//	}
//}
//
//func (this *LolColumnDefinition) String() string {
//	return fmt.Sprintf("'LolTableColumn': {'Name': '%s', 'Type': '%s', 'IsNull': %s}", this.Name, this.Type, this.IsNull)
//}
//
//type LolTableDefinition struct {
//	EntityName string
//	TableName  string
//	Columns    []*LolColumnDefinition
//}
//
//func NewLolTableDefinition(entityName string, tableName string) *LolTableDefinition {
//	return &LolTableDefinition{
//		EntityName: entityName,
//		TableName: tableName,
//		Columns: []*LolColumnDefinition{},
//	}
//}
//
//func (this *LolTableDefinition) ColumnNames() [] string {
//	names := []string{}
//	for _, col := range this.Columns {
//		names = append(names, col.Name)
//	}
//	return names
//}
//
//func (this *LolTableDefinition) String() string {
//	colstr := []string{}
//	for _, col := range this.Columns {
//		colstr = append(colstr, col.String())
//	}
//	return fmt.Sprintf("'LolTableDefinition': {'EntityName': '%s', 'TableName': '%s', 'Columns': [%s]}",
//		this.EntityName, this.TableName, strings.Join(colstr, ", "))
//}
//
//type LolPackageDefinition struct {
//	Path              string
//	EntityDefinitions map[string]*LolEntityDefinition
//}
//
//func NewLolPackageDefinition(path string) *LolPackageDefinition {
//	return &LolPackageDefinition{
//		Path: path,
//		EntityDefinitions: map[string]*LolEntityDefinition{},
//	}
//}
//
//// "github.com/less-leg/dbmodel", "D:/projects/less-leg/src"
//func Parse(packageName string, sourceDir string) *LolPackageDefinition {
//
//	pckg, err := build.Import(packageName, sourceDir, build.IgnoreVendor)
//	if err != nil {
//		panic(fmt.Sprintln(err, pckg))
//	}
//
//
//	tmpDefinitions := map[string]*tmpEntityDefinition{}
//	for _, goFile := range pckg.GoFiles {
//		parseFile(pckg.Dir + "/" + goFile, tmpDefinitions)
//	}
//	pkgDefinition := createLolPackageDefinition(filepath.Join(sourceDir, packageName), tmpDefinitions)
//	return pkgDefinition
//}
//
//func processEmbeddedField(field *ast.Field, tmpDefinitions map[string]*tmpEntityDefinition) ([]LolFieldTypeDefinition, error) {
//	if ident, ok := field.Type.(*ast.Ident); ok {
//		if tmpDef, found := tmpDefinitions[ident.Name]; found {
//			embedded := []LolFieldTypeDefinition{}
//			for _, f := range tmpDef.Fields {
//				fieldTypes, err := getFieldType(f, tmpDefinitions)
//				if err != nil {
//					return nil, err
//				}
//				for _, fieldType := range fieldTypes {
//					fieldType.SetName(ident.Name + "." + fieldType.Name())
//				}
//				embedded = append(embedded, fieldTypes...)
//			}
//			return embedded
//		} else {
//			return errors.New("Undefined embedded type: " + ident.Name)
//		}
//	}
//	return errors.New(fmt.Sprintf("Undefined embedded type: %#v", field))
//}
//
//func getFieldType(field *ast.Field, tmpDefinitions map[string]*tmpEntityDefinition) ([]LolFieldTypeDefinition, error) {
//	if slExpr, ok := field.Type.(*ast.ArrayType); ok {
//		return nil, errors.New(fmt.Sprintf("Doesn't support slices for now: %#v", slExpr))
//	}
//	if star, ok := field.Type.(*ast.StarExpr); ok {
//		if ident, ok := star.X.(*ast.Ident); ok {
//			return []LolFieldTypeDefinition{&LolPtrFieldType{underlying:&LolBasicFieldType{name:ident.Name}}}
//		}
//		if selector, ok := star.X.(*ast.SelectorExpr); ok {
//			return []LolFieldTypeDefinition{&LolPtrFieldType{underlying:&LolBasicFieldType{name:selector.Sel.Name}}}
//		}
//		return nil, errors.New(fmt.Sprintf("Doesn't support for now pointer types to: %#v", star.X))
//	}
//	if ident, ok := field.Type.(*ast.Ident); ok {
//		if len(field.Names) == 0 {
//			// it is embedded struct
//			return processEmbeddedField(field, tmpDefinitions)
//			//if tmpDef, found := tmpDefinitions[ident.Name]; found {
//			//	embedded := []LolFieldTypeDefinition{}
//			//	for _, f := range tmpDef.Fields {
//			//		embedded = append(embedded, getFieldType(f, tmpDefinitions)...)
//			//	}
//			//	return embedded
//			//} else {
//			//	return errors.New("Undefined embedded type: " + ident.Name)
//			//}
//		} else {
//			return []LolFieldTypeDefinition{&LolBasicFieldType{name:ident.Name}}
//		}
//	}
//	if selector, ok := field.Type.(*ast.SelectorExpr); ok {
//		return []LolFieldTypeDefinition{&LolPtrFieldType{underlying:&LolBasicFieldType{name:selector.Sel.Name}}}
//	}
//	panic(fmt.Sprintf("Cannot parse go definition: %#v", field))
//}
//
//func getColumnDefinition(field *ast.Field, fieldDefs []*LolFieldTypeDefinition, tmpDefinitions map[string]*tmpEntityDefinition) ([]*LolColumnDefinition, error) {
//	if field.Tag != nil {
//		tagStart := strings.Index(field.Tag.Value, "lolsql")
//		if tagStart > -1 {
//			hasLolTag = true
//			tagStart = tagStart + len("lolsql:\"")
//			tag = field.Tag.Value[tagStart:]
//			tagEnd := strings.Index(tag, "\"")
//			tag = string(tag[:tagEnd])
//		}
//	} else {
//		fieldDefs
//		LolColumnDefinition{Name:}
//	}
//}
//
//func (this *LolPackageDefinition) LolTableColumns(entity string) []*LolColumnDefinition {
//	lolTableCols := []*LolColumnDefinition{}
//	if entityDefinition, found := this.tmpDefinitions[entity]; found {
//		for field, tag := range entityDefinition.Fields {
//			if slExpr, ok := field.Type.(*ast.ArrayType); ok {
//				if star, ok := slExpr.Elt.(*ast.StarExpr); ok {
//					if ident, ok := star.X.(*ast.Ident); ok {
//						// TODO: This part is responsible for Relationships between entities
//						fmt.Printf("Type is not supported : []*%s\n", ident.Name)
//						continue
//					}
//				}
//			} else if star, ok := field.Type.(*ast.StarExpr); ok {
//				if ident, ok := star.X.(*ast.Ident); ok {
//					lolTableCol, err := retrieveLolColumn(field, ident.Name, tag, true)
//					if err == nil {
//						lolTableCols = append(lolTableCols, lolTableCol)
//					} else {
//						// TODO: This part is responsible for Relationships between entities
//						fmt.Printf("Relationship is not supportede for *%s\n", ident.Name)
//					}
//				} else if selector, ok := star.X.(*ast.SelectorExpr); ok {
//					lolTableCol, err := retrieveLolColumn(field, selector.Sel.Name, tag, true)
//					if err == nil {
//						lolTableCols = append(lolTableCols, lolTableCol)
//					} else {
//						//if ident, ok := selector.X.(*ast.Ident); ok {
//						//   TODO: selector for fields declared with 'package.Type' format
//						//}
//						fmt.Printf("Not supported type with selector %s -> %#v\n", selector.Sel.Name, selector.Sel)
//					}
//				} else {
//					panic("NOT HANDLED: " + fmt.Sprintf("%#v", star.X))
//				}
//			} else if ident, ok := field.Type.(*ast.Ident); ok {
//				if len(field.Names) == 0 {
//					// it is embedded struct
//					columns := this.LolTableColumns(ident.Name)
//					lolTableCols = append(lolTableCols, columns...)
//				} else {
//					lolTableCol, err := retrieveLolColumn(field, ident.Name, tag, false)
//					if err == nil {
//						lolTableCols = append(lolTableCols, lolTableCol)
//					} else {
//						fmt.Printf("Not supported type %s\n", ident.Name)
//					}
//				}
//			} else if selector, ok := star.X.(*ast.SelectorExpr); ok {
//				lolTableCol, err := retrieveLolColumn(field, selector.Sel.Name, tag, false)
//				if err == nil {
//					lolTableCols = append(lolTableCols, lolTableCol)
//				} else {
//					//if ident, ok := selector.X.(*ast.Ident); ok {
//					//   TODO: selector for fields declared with 'package.Type' format
//					//}
//					fmt.Printf("Not supported type with selector %s -> %#v\n", selector.Sel.Name, selector.Sel)
//				}
//			} else {
//				panic("Cannot parse go definition")
//			}
//		}
//	}
//	return lolTableCols
//}
//
//func retrieveLolColumn(field *ast.Field, typeName, tag string, isPointer bool) (*LolColumnDefinition, error) {
//	if colType, found := SupportedTypes[typeName]; found {
//		fieldName := field.Names[0].Name
//		if tag == "" {
//			return NewLolColumnDefinition(fieldName, isPointer, colType), nil
//		} else {
//			colTagStart := strings.Index(tag, "column[")
//			if colTagStart < 0 {
//				return NewLolColumnDefinition(fieldName, isPointer, colType), nil
//			} else {
//				colNameStart := colTagStart + len("column[")
//				colNameEnd := strings.Index(tag[colNameStart:], "]")
//				if colNameEnd < 0 {
//					panic(tag)
//				}
//				colName := string(tag[colNameStart:colNameStart + colNameEnd])
//				return NewLolColumnDefinition(colName, isPointer, colType), nil
//			}
//		}
//	} else {
//		return nil, errors.New("Type is not supported :" + typeName)
//	}
//}
//
//func createLolPackageDefinition(path string, tmpDefinitions map[string]*tmpEntityDefinition) *LolPackageDefinition {
//	pkgDefinition := NewLolPackageDefinition(path)
//
//	for entityName, tmpDef := range tmpDefinitions {
//		fieldDefs := []LolFieldTypeDefinition{}
//		for _, field := range tmpDef.Fields {
//			fDefs, err := getFieldType(field, tmpDefinitions)
//			if err != nil {
//				panic(err)
//			}
//			fieldDefs = append(fieldDefs, fDefs...)
//		}
//
//		//fdef := &LolFieldDefinition {
//		//	Name:"",
//		//	Type:LolFieldTypeDefinition{},
//		//}
//		//
//		//col := &LolColumnDefinition{
//		//	Name:"",
//		//	Type:STRING,
//		//	IsNull:false,
//		//}
//		//
//		//lf := &LolField{
//		//	Field: fdef,
//		//	Column: col,
//		//}
//		//
//		//lfs := []*LolField{}
//		//lfs = append(lfs, lf)
//		//
//		//entDef := &LolEntityDefinition {
//		//	Name:entityName,
//		//	TableName:tmpDef.GetTableName(),
//		//	Fields:lfs,
//		//}
//	}
//
//	for entityName, tmpDef := range tmpDefinitions {
//		for _, field := range tmpDef.Fields {
//			columnDefs, err := getColumnDefinition(field, fieldDefs, tmpDefinitions)
//		}
//	}
//
//	return pkgDefinition
//}
//
////
////func (this *LolPackageDefinition) String() string {
////	defstr := []string{}
////	for _, def := range this.EntityDefinitions {
////		defstr = append(defstr, def.String())
////	}
////	tblstr := []string{}
////	for _, tbl := range this.TableDefinitions {
////		tblstr = append(tblstr, tbl.String())
////	}
////	return fmt.Sprintf("'LolPackageDefinition': {\n\t'Path': '%s', \n\t'EntityDefinitions': [\n\t\t%s], \n\t'TableDefinitions': [\n\t\t%s]}",
////		this.Path, strings.Join(defstr, "\n\t\t"), strings.Join(tblstr, "\n\t\t"))
////}
//
//func parseFile(path string, tmpDefinitions map[string]*tmpEntityDefinition) {
//	fset := token.NewFileSet()
//	tree, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, decl := range tree.Decls {
//		currTmpDefs := parseStructDeclarations(decl)
//		for k, v := range currTmpDefs {
//			if _, found := tmpDefinitions[k]; found {
//				panic("Override of stucts is not supported")
//			}
//			tmpDefinitions[k] = v
//		}
//	}
//}
//
////func (this *LolPackageDefinition) parseStructDeclarations(decl ast.Decl) {
////	if genDecl, ok := decl.(*ast.GenDecl); ok {
////		for _, spec := range genDecl.Specs {
////			tspec, ok := spec.(*ast.TypeSpec)
////			if ok {
////				stype, ok := tspec.Type.(*ast.StructType)
////				if ok {
////					if definition, ok := createTmpEntityDefinition(tspec.Name.Name, stype); ok {
////						if _, found := this.tmpDefinitions[definition.Name]; found {
////							panic("Definition overrides is not supported: " + definition.Name)
////						}
////						this.tmpDefinitions[definition.Name] = definition
////					}
////				}
////			}
////		}
////	}
////}
//
//func parseStructDeclarations(decl ast.Decl) map[string]*tmpEntityDefinition {
//	tmpDefinitions := make(map[string]*tmpEntityDefinition)
//	if genDecl, ok := decl.(*ast.GenDecl); ok {
//		for _, spec := range genDecl.Specs {
//			tspec, ok := spec.(*ast.TypeSpec)
//			if ok {
//				stype, ok := tspec.Type.(*ast.StructType)
//				if ok {
//					if tmpDef, ok := createTmpEntityDefinition(tspec.Name.Name, stype); ok {
//						if _, found := tmpDefinitions[tmpDef.Name]; found {
//							panic("Definition overrides is not supported: " + tmpDef.Name)
//						}
//						tmpDefinitions[tmpDef.Name] = tmpDef
//					}
//				}
//			}
//		}
//	// lookup for overridden table names: func (*T) TableName() string { return "T_NAME" }
//	} else if funcDecl, ok := decl.(*ast.FuncDecl); ok {
//		if funcDecl.Name.Name == "TableName" {
//			if funcDecl.Recv != nil && funcDecl.Type.Params != nil && funcDecl.Type.Results != nil {
//				if len(funcDecl.Type.Params.List) == 0 && len(funcDecl.Type.Results.List) == 1 {
//					res := funcDecl.Type.Results.List[0]
//					if ident, ok := res.Type.(*ast.Ident); ok && ident.Name == "string" {
//						for _, recv := range funcDecl.Recv.List {
//							if star, ok := recv.Type.(*ast.StarExpr); ok {
//								if selfIdent, ok := star.X.(*ast.Ident); ok {
//									if def, ok := tmpDefinitions[selfIdent.Name]; ok {
//										if def.Name == selfIdent.Name {
//											def.TableNameFunc = funcDecl
//										} else {
//											panic("definitions has no definition for " + selfIdent.Name + " struct")
//										}
//									}
//								} else {
//									panic("Not a star.X.(*ast.Ident)")
//								}
//							} else {
//								panic("Not a recv.Type.(*ast.StarExpr)")
//							}
//						}
//					} else {
//						panic("Not a string in return")
//					}
//				} else {
//					panic(fmt.Sprintf("Length of parameters is not as expected %d %d\n", len(funcDecl.Recv.List), len(funcDecl.Type.Results.List)))
//				}
//			}
//		} else {
//			panic("Not a TableName function")
//		}
//	}
//	return tmpDefinitions
//}
//
//func (this *LolPackageDefinition) parseTableDeclarations(decl ast.Decl) {
//	// default table names for structs
//	for _, def := range this.tmpDefinitions {
//		if _, found := this.TableDefinitions[def.Name]; !found {
//			this.TableDefinitions[def.Name] = NewLolTableDefinition(def.Name, def.Name)
//		}
//	}
//	// lookup for overridden table names: func (*T) TableName() string { return "T_NAME" }
//	if funcDecl, ok := decl.(*ast.FuncDecl); ok {
//		if funcDecl.Name.Name == "TableName" {
//			if funcDecl.Recv != nil && funcDecl.Type.Params != nil && funcDecl.Type.Results != nil {
//				if len(funcDecl.Type.Params.List) == 0 && len(funcDecl.Type.Results.List) == 1 {
//					res := funcDecl.Type.Results.List[0]
//					if ident, ok := res.Type.(*ast.Ident); ok && ident.Name == "string" {
//						for _, recv := range funcDecl.Recv.List {
//							if star, ok := recv.Type.(*ast.StarExpr); ok {
//								if selfIdent, ok := star.X.(*ast.Ident); ok {
//									if def, ok := this.tmpDefinitions[selfIdent.Name]; ok {
//										if def.Name == selfIdent.Name {
//											for _, stmt := range funcDecl.Body.List {
//												if rstmt, ok := stmt.(*ast.ReturnStmt); ok && len(rstmt.Results) == 1 {
//													tableNameLiteral := rstmt.Results[0].(*ast.BasicLit).Value
//													tableName := tableNameLiteral[1:len(tableNameLiteral) - 1]
//													this.TableDefinitions[def.Name] = NewLolTableDefinition(def.Name, tableName)
//												} else {
//													panic("Not a stmt.(*ast.ReturnStmt) or length of len(rstmt.Results) != 1")
//												}
//											}
//										} else {
//											panic("definitions has no definition for " + selfIdent.Name + " struct")
//										}
//									}
//								} else {
//									panic("Not a star.X.(*ast.Ident)")
//								}
//							} else {
//								panic("Not a recv.Type.(*ast.StarExpr)")
//							}
//						}
//					} else {
//						panic("Not a string in return")
//					}
//				} else {
//					panic(fmt.Sprintf("Length of parameters is not as expected %d %d\n", len(funcDecl.Recv.List), len(funcDecl.Type.Results.List)))
//				}
//			}
//		} else {
//			panic("Not a TableName function")
//		}
//	}
//}
//
//func (this *LolPackageDefinition) parseTableTmpDeclarations(decl ast.Decl, tmpDefinitions map[string]*tmpEntityDefinition) {
//	// default table names for structs
//	for _, def := range tmpDefinitions {
//		if _, found := this.TableDefinitions[def.Name]; !found {
//			this.TableDefinitions[def.Name] = NewLolTableDefinition(def.Name, def.Name)
//		}
//	}
//	// lookup for overridden table names: func (*T) TableName() string { return "T_NAME" }
//	if funcDecl, ok := decl.(*ast.FuncDecl); ok {
//		if funcDecl.Name.Name == "TableName" {
//			if funcDecl.Recv != nil && funcDecl.Type.Params != nil && funcDecl.Type.Results != nil {
//				if len(funcDecl.Type.Params.List) == 0 && len(funcDecl.Type.Results.List) == 1 {
//					res := funcDecl.Type.Results.List[0]
//					if ident, ok := res.Type.(*ast.Ident); ok && ident.Name == "string" {
//						for _, recv := range funcDecl.Recv.List {
//							if star, ok := recv.Type.(*ast.StarExpr); ok {
//								if selfIdent, ok := star.X.(*ast.Ident); ok {
//									if def, ok := this.tmpDefinitions[selfIdent.Name]; ok {
//										if def.Name == selfIdent.Name {
//											for _, stmt := range funcDecl.Body.List {
//												if rstmt, ok := stmt.(*ast.ReturnStmt); ok && len(rstmt.Results) == 1 {
//													tableNameLiteral := rstmt.Results[0].(*ast.BasicLit).Value
//													tableName := tableNameLiteral[1:len(tableNameLiteral) - 1]
//													this.TableDefinitions[def.Name] = NewLolTableDefinition(def.Name, tableName)
//												} else {
//													panic("Not a stmt.(*ast.ReturnStmt) or length of len(rstmt.Results) != 1")
//												}
//											}
//										} else {
//											panic("definitions has no definition for " + selfIdent.Name + " struct")
//										}
//									}
//								} else {
//									panic("Not a star.X.(*ast.Ident)")
//								}
//							} else {
//								panic("Not a recv.Type.(*ast.StarExpr)")
//							}
//						}
//					} else {
//						panic("Not a string in return")
//					}
//				} else {
//					panic(fmt.Sprintf("Length of parameters is not as expected %d %d\n", len(funcDecl.Recv.List), len(funcDecl.Type.Results.List)))
//				}
//			}
//		} else {
//			panic("Not a TableName function")
//		}
//	}
//}
//
//// return true as second return value if declared type has lolsql marker
//func createTmpEntityDefinition(name string, stp *ast.StructType) (*tmpEntityDefinition, bool) {
//	definition := newTmpEntityDefinition(name)
//	definition.Fields = stp.Fields
//	if stp.Fields != nil {
//		for _, field := range stp.Fields.List {
//			if field.Tag != nil && strings.Index(field.Tag.Value, "lolsql") > -1 {
//				return definition, true
//			}
//			//if field.Tag != nil {
//			//	tagStart := strings.Index(field.Tag.Value, "lolsql")
//			//	if tagStart > -1 {
//			//		hasLolTag = true
//			//		tagStart = tagStart + len("lolsql:\"")
//			//		tag = field.Tag.Value[tagStart:]
//			//		tagEnd := strings.Index(tag, "\"")
//			//		tag = string(tag[:tagEnd])
//			//	}
//			//}
//
//		}
//	}
//	return definition, false
//}
//
//type FieldToColumnName struct {
//	FieldName  string
//	ColumnName string
//}
//
////func (this *LolPackageDefinition) FieldsToColumnNames(entityName string) []FieldToColumnName {
////	fToC := make([]FieldToColumnName, 0)
////	if edef, found := this.EntityDefinitions[entityName]; found {
////		if tdef, found := this.TableDefinitions[entityName]; found {
////			for field, _ := range edef.Fields {
////				for _, name := range field.Names {
////					if name != nil {
////						tdef.Columns
////						fToC = append(fToC, name.Name)
////					}
////				}
////			}
////		}
////	}
////	return fToC
////}
//
//type tmpEntityDefinition struct {
//	Name   string
//	Fields []*ast.Field
//	TableNameFunc *ast.FuncDecl
//}
//
//func newTmpEntityDefinition(name string) *LolEntityDefinition {
//	return &tmpEntityDefinition{
//		Name:name,
//		Fields:make(map[*ast.Field]string)}
//}
//
//func (this *tmpEntityDefinition) GetTableName() string {
//	if this.TableNameFunc == nil {
//		return this.Name
//	}
//	for _, stmt := range this.TableNameFunc.Body.List {
//		if rstmt, ok := stmt.(*ast.ReturnStmt); ok && len(rstmt.Results) == 1 {
//			tableNameLiteral := rstmt.Results[0].(*ast.BasicLit).Value
//			return tableNameLiteral[1:len(tableNameLiteral) - 1]
//		} else {
//			panic("Function TableName for " + this.Name + " has no return statement or invalid amount of return values")
//		}
//	}
//	panic("Function TableName for " + this.Name + " has no return statement")
//}
//
//type LolFieldTypeDefinition interface {
//	SetName(string)
//	Name() string
//	IsPtr() bool
//	IsSlice() bool
//	Underlying() LolFieldTypeDefinition
//}
//
//type LolBasicFieldType struct {
//	name string
//	underlying LolFieldTypeDefinition
//}
//func (this *LolBasicFieldType) SetName(name string) {
//	this.name = name
//}
//func (this *LolBasicFieldType) Name() string {
//	return this.name
//}
//func (this *LolBasicFieldType) IsPtr() bool {
//	return false
//}
//func (this *LolBasicFieldType) IsSlice() bool {
//	return false
//}
//func (this *LolBasicFieldType) Underlying() LolFieldTypeDefinition {
//	return nil
//}
//
//type LolSliceFieldType struct {
//	underlying LolFieldTypeDefinition
//}
//func (this *LolSliceFieldType) SetName(string) {}
//func (this *LolSliceFieldType) Name() string {
//	return ""
//}
//func (this *LolSliceFieldType) IsPtr() bool {
//	return false
//}
//func (this *LolSliceFieldType) IsSlice() bool {
//	return true
//}
//func (this *LolSliceFieldType) Underlying() LolFieldTypeDefinition {
//	return this.underlying
//}
//
//type LolPtrFieldType struct {
//	underlying LolFieldTypeDefinition
//}
//func (this *LolPtrFieldType) SetName(name string) {}
//func (this *LolPtrFieldType) Name() string {
//	return ""
//}
//func (this *LolPtrFieldType) IsPtr() bool {
//	return true
//}
//func (this *LolPtrFieldType) IsSlice() bool {
//	return false
//}
//func (this *LolPtrFieldType) Underlying() LolFieldTypeDefinition {
//	return this.underlying
//}
//
//type LolFieldDefinition struct {
//	Name string
//	Type LolFieldTypeDefinition
//}
//
//type LolField struct {
//	Field *LolFieldDefinition
//	Column *LolColumnDefinition
//}
//
//type LolEntityDefinition struct {
//	Name   string
//	TableName string
//	Fields map[string]LolField
//}
//
//
////func (this *LolEntityDefinition) String() string {
////	fields := make([]string, 0)
////	for field, tag := range this.Fields {
////		for _, name := range field.Names {
////			if name != nil {
////				fields = append(fields, fmt.Sprintf("%s `%s`", name.Name, tag))
////			} else {
////				panic(fmt.Sprintf("Name is nil: %s -> %v", this.Name, field))
////			}
////		}
////	}
////	return fmt.Sprintf("'LolEntityDefinition': {'Name': '%s', 'Fields': [%s]}", this.Name, strings.Join(fields, ", "))
////}
