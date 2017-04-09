package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
				`"github.com/less-leg/types"`,
				`"strings"`,
				`"database/sql"`,
				`"fmt"`,
				`"github.com/less-leg/utils"`},
				utils.DoubleQuotes(append([]string{pckgDef.ModelPackage}, sdef.Selectors(pckgDef.ModelPackageName())...)...)...))

			Scanner_struct.ExecuteTemplate(entityFile, "", struct {
				Package    string
				StructName string
				Fields     []*parser.FetchMeta
			}{
				Package:    pckgDef.ModelPackageName(),
				StructName: sdef.Name(),
				Fields:     sdef.FetchMeta(pckgDef),
			})

			Lol_struct.ExecuteTemplate(entityFile, "", struct {
				Package            string
				StructName         string
				TableNameToColumns []string
			}{
				Package:            pckgDef.ModelPackageName(),
				StructName:         sdef.Name(),
				TableNameToColumns: []string{sdef.TableName, strings.Join(pckgDef.ColumnNames(structName), ", ")},
			})

			Select_func.ExecuteTemplate(entityFile, "", nil)

			LolWhere_struct.ExecuteTemplate(entityFile, "", struct {
				Package    string
				StructName string
			}{
				Package:    pckgDef.ModelPackageName(),
				StructName: sdef.Name(),
			})

			ColumnStub_struct.ExecuteTemplate(entityFile, "", pckgDef.FieldsToColumns(structName))

			generateFields(entityFile, pckgDef, sdef, "")
			utils.PanicIf(entityFile.Close())

		case *parser.EmbeddedStructDefinition:
			fmt.Printf("EmbeddedStructDefinition: %#v\n", sdef)
		case *parser.CustomUserTypeStructDefinition:
			fmt.Printf("CustomUserTypeStructDefinition: %#v\n", sdef)
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
				TypeName      string
				StructName    string
				IsNullable    string
				FieldToColumn []string
				Likable       bool
				Selector      string
			}{
				TypeName:      fdef.FieldType().Name(),
				StructName:    sdef.Name(),
				IsNullable:    fdef.FieldType().PtrSign(),
				FieldToColumn: []string{fdef.Name(), fdef.ColumnName},
				Likable:       fdef.FieldType().IsLikable(),
				Selector:      selector,
			})
		case *parser.ComplexFieldDefinition:
			if fdef.Embedded {
				if embdStrtDef, found := pckgDef.StructDefinitions[fdef.Name()]; found {
					fmt.Printf("Embedded field found: %#v\n", embdStrtDef)
					generateFields(entityFile, pckgDef, embdStrtDef, fdef.Name())
				}
			} else {
				ConditionByField.ExecuteTemplate(entityFile, "", struct {
					TypeName      string
					StructName    string
					IsNullable    string
					FieldToColumn []string
					Likable       bool
					Selector      string
				}{
					TypeName:      fdef.FieldType().Name(),
					StructName:    sdef.Name(),
					IsNullable:    fdef.FieldType().PtrSign(),
					FieldToColumn: []string{fdef.Name(), fdef.ColumnName},
					Likable:       fdef.FieldType().IsLikable(),
					Selector:      selector,
				})
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
