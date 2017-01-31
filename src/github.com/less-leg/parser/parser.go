package parser

import (
	"go/build"
	"go/token"
	"go/parser"
	"go/ast"
	"github.com/less-leg/utils"
)

func Parse(packageName string, sourceDir string) []*ParsedStruct {
	pckg, err := build.Import(packageName, sourceDir, build.IgnoreVendor)
	utils.PanicIf(err)

	var decls []ast.Decl
	for _, goFile := range pckg.GoFiles {
		decls = append(decls, parseFile(pckg.Dir + "/" + goFile)...)
	}

	structs := map[string]*ast.StructType{}
	methods := map[string]map[string]*ast.FuncDecl{}

	for _, decl := range decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if tspec, ok := spec.(*ast.TypeSpec); ok {
					if stype, ok := tspec.Type.(*ast.StructType); ok {
						structs[tspec.Name.Name] = stype
					}
				}
			}
		} else if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv != nil && len(funcDecl.Recv.List) == 1 {
			for _, recv := range funcDecl.Recv.List {
				if star, ok := recv.Type.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						if ms, found := methods[ident.Name]; found {
							ms[funcDecl.Name.Name] = funcDecl;
						} else {
							methods[ident.Name] = map[string]*ast.FuncDecl{funcDecl.Name.Name: funcDecl}
						}
					}
				} else if ident, ok := recv.Type.(*ast.Ident); ok {
					if ms, found := methods[ident.Name]; found {
						ms[funcDecl.Name.Name] = funcDecl;
					} else {
						methods[ident.Name] = map[string]*ast.FuncDecl{funcDecl.Name.Name: funcDecl}
					}
				}
			}
		}
	}

	var ps []*ParsedStruct
	for structName, structType := range structs {
		ps = append(ps, &ParsedStruct{Name:structName, Type:structType, Methods:methods[structName]})
	}
	return ps
}

func parseFile(path string) []ast.Decl {
	tree, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
	utils.PanicIf(err)
	return tree.Decls
}
