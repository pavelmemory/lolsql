package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

type Package struct {
	Name    string
	Types   map[TypeIdentity]Type // by name dictionary
	Imports map[string]Import     // by name dictionary
}

func ReadPackageInfo(packageName string) (Package, error) {
	pkg := Package{Name: packageName, Imports: make(map[string]Import), Types: make(map[TypeIdentity]Type)}
	srcDirs := build.Default.SrcDirs()
	var (
		pkgInfo *build.Package
		err     error
	)

	// TODO: mostly wan't be in the GOPATH
	for _, srcDir := range srcDirs {
		pkgInfo, err = build.Default.Import(packageName, srcDir, build.IgnoreVendor)
		if err != nil {
			return Package{}, err
		}
	}
	if err != nil {
		return Package{}, err
	}

	srcPkgs, err := parser.ParseDir(token.NewFileSet(), filepath.Join(pkgInfo.SrcRoot, packageName), nil, parser.AllErrors)
	if err != nil {
		return Package{}, err
	}
	if len(srcPkgs) != 1 {
		return Package{}, fmt.Errorf("support only for 1 package, %d provided", len(srcPkgs))
	}

	for _, srcPkg := range srcPkgs {
		for fileName, file := range srcPkg.Files {
			for _, decl := range file.Decls {
				switch tDecl := decl.(type) {
				case *ast.GenDecl:
					switch tDecl.Tok {
					case token.TYPE:
						for _, spec := range tDecl.Specs {
							switch tSpec := spec.(type) {
							case *ast.TypeSpec:
								pkg.AddTypeDeclaration(tSpec)
							default:
								log.Printf("unhandled general declaration specification: %#v\n", spec)
							}
						}
					case token.VAR:
						log.Printf("unhandled var declaration: %#v\n", tDecl)
					case token.CONST:
						log.Printf("unhandled const declaration: %#v\n", tDecl)
					case token.IMPORT:
						for _, spec := range tDecl.Specs {
							tSpec := spec.(*ast.ImportSpec)
							pkg.AddImportDeclaration(tSpec)
						}
					default:
						panic("unreachable code")
					}

				case *ast.BadDecl:
					return Package{}, fmt.Errorf("bad declaration at position %d in file %s", tDecl.From, fileName)
				case *ast.FuncDecl:
					fmt.Println(tDecl)
				}
			}
		}
	}
	return pkg, err
}

func (pkg *Package) AddTypeDeclaration(typeSpec *ast.TypeSpec) {
	// TODO: may be useful
	//typeSpec.Name.IsExported()

	switch ttypeSpec := typeSpec.Type.(type) {
	/*
		type MyTimeType time.Time
	*/
	case *ast.SelectorExpr:
		switch tSpecX := ttypeSpec.X.(type) {
		case *ast.Ident:
			t := UserDefinedAlias{
				TypeIdentity: TypeIdentity{Name: typeSpec.Name.Name, Package: pkg.Name},
				ActualType:   TypeIdentityRef{Name: ttypeSpec.Sel.Name, Selector: tSpecX.Name},
			}
			if _, found := pkg.Types[t.GetIdentity()]; !found {
				pkg.Types[t.GetIdentity()] = t
			}
		}

	/*
		type MyStringType string
	*/
	case *ast.Ident:
		if ttypeSpec.Obj == nil {
			t := UserDefinedAlias{
				TypeIdentity: TypeIdentity{Name: typeSpec.Name.Name, Package: pkg.Name},
				ActualType:   TypeIdentityRef{Name: ttypeSpec.Name},
			}
			if _, found := pkg.Types[t.GetIdentity()]; !found {
				pkg.Types[t.GetIdentity()] = t
			}
		}

		switch ttypeSpec.Obj.Kind {
		case ast.Typ:
			genDecl := ttypeSpec.Obj.Decl.(*ast.GenDecl)
			switch genDecl.Tok {
			case token.TYPE:
				for _, spec := range genDecl.Specs {
					switch spec.(*ast.TypeSpec).Type.(type) {
					case *ast.StructType:
						pkg.AddTypeFromStructType(typeSpec)
					default:
						panic("unhandled spec.(*ast.TypeSpec).Type.(type)")
					}
				}
			default:
				fmt.Printf("unknown genDecl: %#v\n", genDecl)
			}
		default:
			fmt.Printf("unknown ttypeSpec: %#v\n", ttypeSpec)
		}

	/*
		type UserModel {Name string}
	*/
	case *ast.StructType:
		pkg.AddTypeFromStructType(typeSpec)
	default:
		panic("unhandled typeSpec.UserDefinedType:" + fmt.Sprintf(" %#v\n", typeSpec))
	}
}

func (pkg *Package) AddTypeFromStructType(typeSpec *ast.TypeSpec) {
	typeDecl := UserDefinedType{
		TypeIdentity: TypeIdentity{Name: typeSpec.Name.Name, Package: pkg.Name},
		Fields:       make(map[string]Field),
	}
	if _, found := pkg.Types[typeDecl.GetIdentity()]; !found {
		pkg.Types[typeDecl.GetIdentity()] = typeDecl
	}

	for _, fieldSpec := range typeSpec.Type.(*ast.StructType).Fields.List {
		fields := FieldsFromField(fieldSpec)

		for _, field := range fields {
			typeDecl.Fields[field.GetName()] = field
		}
	}
}

func FieldsFromField(fieldSpec *ast.Field) []Field {
	field := FieldFromType(fieldSpec.Type)
	if len(fieldSpec.Names) == 0 {
		/*
			type D struct {string}
		*/
		field.IsEmbedded = true
		return []Field{field}
	} else {
		/*
			type D struct {F, K string}
		*/
		var fields []Field
		for _, name := range fieldSpec.Names {
			field.Name = name.Name
			fields = append(fields, field)
		}
		return fields
	}
}

// returns type's name of struct's field
func FieldFromType(fieldTypeExpr ast.Expr) UserDefinedTypeField {
	var fieldType UserDefinedTypeField
ComplexDeclaration:
	for {
		switch tFieldSpec := fieldTypeExpr.(type) {
		/*
			type D struct {time.Time}
		*/
		case *ast.SelectorExpr:
			tFieldSpecX := tFieldSpec.X.(*ast.Ident)
			fieldType.Name = tFieldSpec.Sel.Name
			fieldType.Selector = tFieldSpecX.Name
			fieldType.TypeName = tFieldSpec.Sel.Name

		/*
			type D struct {string}
		*/
		case *ast.Ident:
			fieldType.Name = tFieldSpec.Name
			fieldType.TypeName = tFieldSpec.Name

		/*
			type D struct {*string}
			type D struct {*time.Time}
		*/
		case *ast.StarExpr:
			fieldType.IsReference = true
			fieldType.OrderedTypeSpecifications = append(fieldType.OrderedTypeSpecifications, Reference)
			fieldTypeExpr = tFieldSpec.X
			continue ComplexDeclaration

		/*
			type D struct {U []string}
			type D struct {[]*time.Time}
		*/
		case *ast.ArrayType:
			fieldType.IsSlice = true
			fieldType.OrderedTypeSpecifications = append(fieldType.OrderedTypeSpecifications, Slice)
			fieldTypeExpr = tFieldSpec.Elt
			continue ComplexDeclaration
		}
		break ComplexDeclaration
	}
	return fieldType
}

func (pkg *Package) AddImportDeclaration(iSpec *ast.ImportSpec) {
	imp := Import{Path: iSpec.Path.Value[1 : len(iSpec.Path.Value)-1]}
	pathParts := strings.Split(imp.Path, "/")
	imp.Alias = pathParts[len(pathParts)-1]
	if iSpec.Name != nil {
		imp.Alias = iSpec.Name.Name
	}
	pkg.Imports[imp.Path] = imp
}
