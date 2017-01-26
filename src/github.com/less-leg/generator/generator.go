package generator

import (
	"strings"
	"path/filepath"
	"io"
	"os"
	"fmt"

	"github.com/less-leg/parser"
	"github.com/less-leg/utils"
)

func Generate(pckgDef *parser.PackageDefinition) {
	for structName, sdef := range pckgDef.StructDefinitions {
		switch sdef := sdef.(type) {
		case *parser.TableStructDefinition:
			structNameLowCase := strings.ToLower(structName)
			entPackageDirPath := utils.RecreateDirectory(filepath.Join(pckgDef.PackageDirPath, structNameLowCase))
			entityFile := createEntityFile(entPackageDirPath, structNameLowCase)

			Package.ExecuteTemplate(entityFile, "", structNameLowCase)
			Imports.ExecuteTemplate(entityFile, "", append([]string{
				`. "github.com/less-leg/types"`,
				`"strings"`,
				`"github.com/less-leg/utils"`},
				utils.DoubleQuote(sdef.Selectors()...)...))
			Column_interface.ExecuteTemplate(entityFile, "", nil)
			Lol_struct.ExecuteTemplate(entityFile, "", []string{
				sdef.TableName,
				strings.Join(pckgDef.ColumnNames(structName), ", ")})
			Select_func.ExecuteTemplate(entityFile, "", nil)
			LolWhere_struct.ExecuteTemplate(entityFile, "", nil)
			LolConditionAnd_struct.ExecuteTemplate(entityFile, "", nil)
			LolConditionOr_struct.ExecuteTemplate(entityFile, "", nil)
			ColumnStub_struct.ExecuteTemplate(entityFile, "", pckgDef.FieldsToColumns(structName))

			generateFields(entityFile, pckgDef, sdef, "")
			utils.PanicIf(entityFile.Close())

		case *parser.EmbeddedStructDefinition:
		default:
			panic("Unreachable code")
		}
	}
}

func generateFields(entityFile io.Writer, pckgDef *parser.PackageDefinition, sdef parser.StructDefinition, selector string) {
	for _, fdef := range sdef.FieldDefinitions() {
		switch fdef := fdef.(type) {
		case *parser.SimpleFieldDefinition:
			ConditionByField.ExecuteTemplate(entityFile, "", struct {
				TypeName string
				StructName string
				IsNullable string
				FieldToColumn []string
				Likable bool
				Selector string
			}{
				TypeName:   fdef.FieldType().Name(),
				StructName: sdef.Name(),
				IsNullable: fdef.FieldType().PtrSign(),
				FieldToColumn: []string{fdef.Name(), fdef.ColumnName},
				Likable: fdef.FieldType().Name() == "string",
				Selector: selector,
			})
		case *parser.ComplexFieldDefinition:
			if fdef.Embedded {
				if embdStrtDef, found := pckgDef.StructDefinitions[fdef.Name()]; found {
					generateFields(entityFile, pckgDef, embdStrtDef, fdef.Name())
				}
			}
		}
	}
}

func createEntityFile(pkgFilePath, entityName string) *os.File {
	entityFileName := entityName + "_lol.go"
	entityFilePath := filepath.Join(pkgFilePath, entityFileName)
	entityFile, err := os.Create(entityFilePath)
	if os.IsExist(err) {
		panic(fmt.Sprintf("File %s cannot be created: %s", entityFilePath, err.Error()))
	}
	return entityFile
}