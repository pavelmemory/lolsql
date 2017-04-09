package main

import (
	//_ "github.com/go-sql-driver/mysql"

	"path/filepath"

	"github.com/less-leg/utils"
	"github.com/less-leg/parser"
	"github.com/less-leg/sql/generator"
	"fmt"
	. "github.com/less-leg/dbmodel/lolsql/djangoadminlog"
	"github.com/less-leg/dbmodel"
	//"database/sql"
	"github.com/less-leg/dbmodel/lolsql/person"
)

func main() {
	generate := true

	if generate {
		packageDir := "github.com/less-leg/dbmodel"
		sourceDir := "D:/workspace/GoProjects/lolsql/src"
		parsedStructs := parser.Parse(packageDir, sourceDir)

		fmt.Println("PARSED")
		for _, parsedStruct := range parsedStructs {
			fmt.Printf("%#v\n", parsedStruct)
		}

		fmt.Println("DEFINITIONS")
		lolDirPath := utils.RecreateDirectory(filepath.Join(sourceDir, packageDir, "lolsql"))
		pckgDef := parser.NewPackageDefinition(packageDir, lolDirPath, parsedStructs)
		for _, strDef := range pckgDef.StructDefinitions {
			fmt.Printf("%#v\n", strDef)
		}

		generator.Generate(pckgDef)
	} else {

		//ids := []int{10, 102}
		//fmt.Println(
		//	person.Select(person.Id(), person.Password()).
		//	Where(person.IdIs(10).And(person.IdIs(ids[0], ids[1]))).Or(person.PasswordIsNotNull().And(person.PasswordLike("%10", "001_"))).Render())

		//now := time.Now()
		//fmt.Println(handsome.Select(handsome.Login(), handsome.DateOfBirth()).Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now).Or(handsome.SalaryIs(10.2, 100.2).Or(handsome.SalaryIs(-100)))).And(handsome.LoginIs("LoginIs")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIs("LoginIs", "LoginIs2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot", "LoginIsNot2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("%P1")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("%P1")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("LoginLike", "LoginLike2")).Render())
		//fmt.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("LoginNotLike", "LoginNotLike2")).Render())

		//db, err := sql.Open("mysql", "root:root@/akeos?parseTime=true")
		//utils.PanicIf(err)
		//defer db.Close()
		////
		//dals, err := djangoadminlog.Select().Where(djangoadminlog.ObjectReprLike("Harry%", "Lol%").And(djangoadminlog.ActionFlagIs(1)).And(djangoadminlog.ActionFlagIs(1))).Fetch(db)
		fmt.Println(
			Select(ChangeMessage(), ActionTime()).
				Where(ObjectReprIs("Harry%", "Lol%").And(ActionFlagIs(1)).And(ActionFlagIs(1))).And(AskPasswordIs(dbmodel.No)).
				Render(),
		)
		//fmt.Println(len(dals))
		//var db *sql.DB
		// Select().Where(AskPasswordIs(dbmodel.Yes)).Fetch(db)
		//fmt.Println(len(dals))

		//ps, err := person.Select().Fetch(db)
		//fmt.Println(err, ps)
		//

		//sqlRend := person.Select(person.Id(), person.Password()).
		sqlRend := person.Select(person.Name_First()).
			Where(person.IdIs(10).And(person.IdIs(10, 1000))).
				Or(person.PasswordIsNotNull().And(person.PasswordLike("%10", "001_"))).
				Or(person.FirstNotLike("Pavrl")).
				Or(person.SecondIsNull()).
			Render()
		fmt.Println(sqlRend)

		//ps, err := person.Select(person.Id(), person.Password()).
		//	Where(person.IdIs(10).And(
		//		person.IdIs(10, 1000))).Or(
		//			person.PasswordIsNotNull().And(
		//			person.PasswordLike("%10", "001_")),
		//).Fetch(db)

		//r,e := db.Query("select Id, SECRET from PERSONS where SECRET like ?", 10)
		//fmt.Println(e, r)
		//
		//
		//fmt.Println(person.Select(person.Id(), person.Password()).Where(person.PasswordLike("%")).Render())
		//ps, err := person.Select(person.Id(), person.Password()).Where(person.PasswordLike("%")).Fetch(db)
		//
		//fmt.Println(err, ps)
	}
}