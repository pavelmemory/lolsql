package main

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/less-leg/dbmodel/person"
	"github.com/less-leg/parser"
	"strings"
	"github.com/less-leg/generator"
	//"github.com/less-leg/dbmodel/lolsql/person"
)

func main() {
	generate := true

	if generate {
		packageDir := "github.com/less-leg/dbmodel"
		sourceDir := "D:/projects/less-leg/src"
		parsedStructs := parser.Parse(packageDir, sourceDir)
		lolDirPath := createDirectory(filepath.Join(sourceDir, packageDir, "lolsql"))
		pckgDef := parser.NewPackageDefinition(lolDirPath, parsedStructs)
		// code below responsible for code generation
		generateLol(pckgDef)
	} else {
		TestAllThatShit()

		//ids := []int{10, 102}
		//fmt.Println(
		//	person.Select(
		//		person.Id(), person.Name_FirstName(), person.Name_MiddleName()).
		//	Where(person.IdIs().And(person.IdIs(&ids[0], &ids[1]))).Render())
	}
}

func generateLol(pckgDef *parser.PackageDefinition) {
	for structName, structDefinition := range pckgDef.StructDefinitions {
		switch sdef := structDefinition.(type) {
		case *parser.TableStructDefinition:

			structNameLowCase := strings.ToLower(structName)
			entPackageDirPath := createDirectory(filepath.Join(pckgDef.PackageDirPath, structNameLowCase))
			entityFile := createEntityFile(entPackageDirPath, structNameLowCase)

			generator.Package.ExecuteTemplate(entityFile, "", structNameLowCase)
			generator.Imports.ExecuteTemplate(entityFile, "",
				`. "github.com/less-leg/types"
				"strings"
				"fmt"
				"strconv"`)
			generator.Column_interface.ExecuteTemplate(entityFile, "", nil)
			generator.Lol_struct.ExecuteTemplate(entityFile, "", struct{
				TableName string
				Columns   string
			}{
				TableName:sdef.TableName,
				Columns: strings.Join(pckgDef.ColumnNames(structName), ", "),
			})
			generator.Select_func.ExecuteTemplate(entityFile, "", nil)
			generator.LolWhere_struct.ExecuteTemplate(entityFile, "", nil)
			generator.LolConditionAnd_struct.ExecuteTemplate(entityFile, "", nil)
			generator.LolConditionOr_struct.ExecuteTemplate(entityFile, "", nil)
			generator.ColumnStub_struct.ExecuteTemplate(entityFile, "", pckgDef.FieldsToColumns(structName))

			for _, fdef := range sdef.FieldDefinitions() {
				switch fdef := fdef.(type) {
				case *parser.SimpleFieldDefinition:
					if tmpl, found := generator.Conditions[fdef.FieldType.Name()]; found {
						tmpl.ExecuteTemplate(entityFile, "", struct {
							StructName string
							TypeName string
							FieldToColumn parser.TupleStringString
						}{
							StructName:structName,
							TypeName: fdef.FieldType.Name(),
							FieldToColumn:parser.TupleStringString{Value1:fdef.Name(), Value2:fdef.ColumnName},
						})
					} else {
						panic("Generation template doesn't exist for type: " + fdef.String())
					}

				//case *parser.ComplexFieldDefinition:
				//	fdef.Embedded
				}
			}



			entityFile.Close()

		case *parser.EmbeddedStructDefinition:
		default:
			panic("Unreachable code")
		}
	}
}

func createDirectory(fileDir string) string {
	err := os.Mkdir(fileDir, os.ModePerm)
	if os.IsExist(err) {
		err = os.RemoveAll(fileDir)
		if os.IsExist(err) {
			panic(fmt.Sprintf("Directory %s cannot be removed: %s", fileDir, err.Error()))
		}
	}
	return fileDir
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


func createEntityFile(pkgFilePath, entityName string) (entityFile *os.File) {
	entityFileName := entityName + "_lol.go"
	entityFilePath := filepath.Join(pkgFilePath, entityFileName)
	entityFile, err := os.Create(entityFilePath)
	if err != nil {
		panic(err)
	}
	return
}