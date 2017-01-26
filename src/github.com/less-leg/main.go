package main

import (
	"path/filepath"

	"github.com/less-leg/utils"
	"github.com/less-leg/parser"
	"github.com/less-leg/generator"
	"fmt"
	"github.com/less-leg/dbmodel/lolsql/djangoadminlog"
)

func main() {
	generate := false

	if generate {
		packageDir := "github.com/less-leg/dbmodel"
		sourceDir := "D:/projects/less-leg/src"
		parsedStructs := parser.Parse(packageDir, sourceDir)
		lolDirPath := utils.RecreateDirectory(filepath.Join(sourceDir, packageDir, "lolsql"))
		pckgDef := parser.NewPackageDefinition(lolDirPath, parsedStructs)

		generator.Generate(pckgDef)
	} else {

		//ids := []int{10, 102}
		//fmt.Println(
		//	person.Select(
		//		person.Id(), person.Name_FirstName(), person.Name_MiddleName()).
		//	//Where(person.IdIs().And(person.IdIs(&ids[0], &ids[1]))).Render())
		//	Where(person.IdIs(1).And(person.IdIs(ids[0], ids[1]))).Render())

		//now := time.Now()
		//fmt.Println(handsome.Select(handsome.Login(), handsome.DateOfBirth()).Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now).Or(handsome.SalaryIs(10.2, 100.2).Or(handsome.SalaryIs(-100)))).And(handsome.LoginIs("LoginIs")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIs("LoginIs", "LoginIs2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot", "LoginIsNot2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("%P1")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("%P1")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("LoginLike", "LoginLike2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("LoginNotLike", "LoginNotLike2")).Render())
		fmt.Println(djangoadminlog.Select().Where(djangoadminlog.ActionFlagIs(10).Render()))
	}
}