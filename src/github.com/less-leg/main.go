package main

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/less-leg/dbmodel/person"
	"github.com/less-leg/parser"
	"strings"
	"github.com/less-leg/generator"
	"github.com/less-leg/utils"
	"io"
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
		//	//Where(person.IdIs().And(person.IdIs(&ids[0], &ids[1]))).Render())
		//	Where(person.IdIs(1).And(person.IdIs(ids[0], ids[1]))).Render())
		//
		//now := time.Now()
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).Render())
	}
}

func generateLol(pckgDef *parser.PackageDefinition) {
	for structName, sdef := range pckgDef.StructDefinitions {
		switch sdef := sdef.(type) {
		case *parser.TableStructDefinition:
			structNameLowCase := strings.ToLower(structName)
			entPackageDirPath := createDirectory(filepath.Join(pckgDef.PackageDirPath, structNameLowCase))
			entityFile := createEntityFile(entPackageDirPath, structNameLowCase)

			generator.Package.ExecuteTemplate(entityFile, "", structNameLowCase)
			generator.Imports.ExecuteTemplate(entityFile, "", append([]string{
				`. "github.com/less-leg/types"`,
				`"strings"`,
				`"strconv"`,
				`"github.com/less-leg/utils"`},
				utils.DoubleQuote(sdef.Selectors()...)...))
			generator.Column_interface.ExecuteTemplate(entityFile, "", nil)
			generator.Lol_struct.ExecuteTemplate(entityFile, "", []string{
				sdef.TableName,
				strings.Join(pckgDef.ColumnNames(structName), ", ")})
			generator.Select_func.ExecuteTemplate(entityFile, "", nil)
			generator.LolWhere_struct.ExecuteTemplate(entityFile, "", nil)
			generator.LolConditionAnd_struct.ExecuteTemplate(entityFile, "", nil)
			generator.LolConditionOr_struct.ExecuteTemplate(entityFile, "", nil)
			generator.ColumnStub_struct.ExecuteTemplate(entityFile, "", pckgDef.FieldsToColumns(structName))

			generateFields(entityFile, pckgDef, sdef)
			utils.PanicIf(entityFile.Close())
		case *parser.EmbeddedStructDefinition:
		default:
			panic("Unreachable code")
		}
	}
}

func generateFields(entityFile io.Writer, pckgDef *parser.PackageDefinition, sdef parser.StructDefinition) {
	for _, fdef := range sdef.FieldDefinitions() {
		switch fdef := fdef.(type) {
		case *parser.SimpleFieldDefinition:
			generator.ConditionByField.ExecuteTemplate(entityFile, "", struct {
				TypeName string
				StructName string
				IsNullable string
				FieldToColumn []string
				ValueToStringFunc generator.ValueToStringFunc
			}{
				TypeName:   fdef.FieldType().Name(),
				StructName: sdef.Name(),
				IsNullable: fdef.FieldType().PtrSign(),
				FieldToColumn: []string{fdef.Name(), fdef.ColumnName},
				ValueToStringFunc: generator.ValueToStringFuncs.Get(fdef.FieldType().Name()),
			})
		case *parser.ComplexFieldDefinition:
			if fdef.Embedded {
				if embdStrtDef, found := pckgDef.StructDefinitions[fdef.Name()]; found {
					generateFields(entityFile, pckgDef, embdStrtDef)
				}
			}
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