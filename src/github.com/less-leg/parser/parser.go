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
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if tspec, ok := spec.(*ast.TypeSpec); ok {
					if stype, ok := tspec.Type.(*ast.StructType); ok {
						structs[tspec.Name.Name] = stype
					}
				}
			}
		case *ast.FuncDecl:
			if decl.Recv == nil || len(decl.Recv.List) != 1 {
				break
			}

			recvType := decl.Recv.List[0].Type
			if star, ok := recvType.(*ast.StarExpr); ok {
				addMethod(star.X, decl, methods)
			} else {
				addMethod(recvType, decl, methods)
			}
		}
	}

	var ps []*ParsedStruct
	for structName, structType := range structs {
		ps = append(ps, &ParsedStruct{Name:structName, Type:structType, Methods:methods[structName]})
	}
	return ps
}

func addMethod(expr ast.Expr, decl *ast.FuncDecl, methods map[string]map[string]*ast.FuncDecl) {
	if ident, ok := expr.(*ast.Ident); ok {
		if ms, found := methods[ident.Name]; found {
			ms[decl.Name.Name] = decl;
		} else {
			methods[ident.Name] = make(map[string]*ast.FuncDecl)
			methods[ident.Name][decl.Name.Name] = decl
		}
	}
}

func parseFile(path string) []ast.Decl {
	tree, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
	utils.PanicIf(err)
	return tree.Decls
}
