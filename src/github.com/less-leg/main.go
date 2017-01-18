package main

import (
	"fmt"
	"strings"
	"os"
	"path/filepath"

	. "github.com/less-leg/dbmodel/person"
	"github.com/less-leg/generator/definition"
	"github.com/less-leg/generator"
)

func createEntityFile(pkgFilePath, entityName string) (entityFile *os.File) {
	entityFileName := entityName + "_lol.go"
	entityFilePath := filepath.Join(pkgFilePath, entityFileName)
	entityFile, err := os.Create(entityFilePath)
	if err != nil {
		panic(err)
	}
	return
}

//{{ define "" }}
//{{.}}
//{{ end }}

func main() {
	srcRoot := "D:/projects/less-leg/src"
	pkg := "github.com/less-leg/dbmodel"

	pkgDefinition := definition.Parse(pkg, srcRoot)
	fmt.Println(pkgDefinition)

	// code below responsible for code generation

	pkgFilePath := filepath.Join(pkgDefinition.Path, "lolsql")
	err := os.Mkdir(pkgFilePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, def := range pkgDefinition.EntityDefinitions {
		defNameLowCase := strings.ToLower(def.Name)
		entityFile := createEntityFile(pkgFilePath, defNameLowCase)

		generator.Package.ExecuteTemplate(entityFile, "", defNameLowCase)
		generator.Imports.ExecuteTemplate(entityFile, "", ". \"github.com/less-leg/types\"\n\"strings\"")
		generator.Column_interface.ExecuteTemplate(entityFile, "", nil)
		generator.Lol_struct.ExecuteTemplate(entityFile, "", struct{
			TableName string
			Columns   string
		}{
			TableName:pkgDefinition.TableDefinitions[def.Name].TableName,
			Columns: strings.Join(pkgDefinition.TableDefinitions[def.Name].ColumnNames(), ", "),
		})
		generator.Select_func.ExecuteTemplate(entityFile, "", nil)
		generator.LolWhere_struct.ExecuteTemplate(entityFile, "", nil)
		generator.LolConditionAnd_struct.ExecuteTemplate(entityFile, "", nil)
		generator.LolConditionOr_struct.ExecuteTemplate(entityFile, "", nil)
		generator.ColumnStub.ExecuteTemplate(entityFile, "", def.FieldsToColumnNames())

		entityFile.Close()
	}

	TestAllThatShit()
}

func TestAllThatShit() {
	id := 100
	id2 := 102
	id3 := 103
	var id4 *int
	name := "Pavel"
	name2 := "Pavel2"
	sql := Select().Where(IdIs(&id).And(NameIs(&name))).Or(IdIs(&id2, &id3, id4)).Render()
	fmt.Println(sql)

	fmt.Println(Select(Id(), Name()).Where(NameIs(&name)).Render())
	fmt.Println(Select(Id(), Name()).Render())
	fmt.Println(
		Select(Id(), Name()).
		Where(NameIsNot(&name).And(IdIs(&id))).
			Or(IdIs(&id2, &id2).And(NameIsNot(&name, &name2)).Or(NameIs(&name2))).
			Render())
}