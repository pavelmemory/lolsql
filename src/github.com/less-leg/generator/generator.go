package generator

import (
	"os"
	"path/filepath"
	"errors"
	//"go/parser"
	//"go/token"
	//"go/build"
	//
	//"github.com/less-leg/utils"
)

type SQLDefinitionJSON struct {
	Primary   bool
	Type      string
	Name      string
	Query     string
	Converter string
}

type SQLDefinition struct {
	Primary   bool
	Type      string
	Name      string
	Query     string
	Converter SQLConverter
}

type SQLConverter interface {
	ToSQLValue(interface{}) interface{}
	ToGoValue(string) interface{}
}

func Generate() {
	g := new(Generator)
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic(errors.New("GOPATH enviroment variable is not set"))
	}
	filepath.Walk(goPath, filepath.WalkFunc(g.ScanFile))
}

type Generator struct {
	pckgs []string // packages to scan for generation
}

func NewGenerator(pckgs ... string) *Generator {
	return &Generator{pckgs: pckgs}
}

func (g *Generator) ScanFile(path string, info os.FileInfo, err error) error {

	//dir, err :=build.ImportDir(".", build.IgnoreVendor)
	//dir.GoFiles
	//
	//if err != nil {
	//	return err
	//}
	//if info.IsDir() {
	//	return nil
	//}
	//f, err := os.Open(path)
	//if err != nil {
	//	return err
	//}
	//set := token.NewFileSet()
	//ast, err := parser.ParseFile(set, f.Name(), nil, parser.AllErrors)
	//utils.PanicIf(err)
	//for _, decl := range ast.Decls {
	//	decl.
	//}
	return nil
}


func ReadDir() {
	//dir, e := os.Open(".")
	//utils.PanicIf(e)
	//fileNames, e := dir.Readdirnames(-1)
	//utils.PanicIf(e)
	//for _, fileName := range fileNames {
	//	file, e := os.Open(fileName)
	//
	//}
}