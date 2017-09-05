package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/less-leg/dbmodel/lolsql/person"
	"github.com/less-leg/parser"
	"github.com/less-leg/sql/generator"
	"github.com/less-leg/utils"
)

func main() {
	generate := true

	if generate {
		packageDir := "github.com/less-leg/dbmodel" // TODO: it is must be an input argument
		goPathValues := os.Getenv("GOPATH")
		if goPathValues == "" {
			log.Fatalln("OS env var GOPATH is not defined")
		}
		paths := strings.Split(goPathValues, fmt.Sprint(os.PathListSeparator))

		var parsedStructs []*parser.ParsedStruct
		for _, path := range paths {
			goPathValues = filepath.Join(path, "src")
			parsed, err := parser.Parse(packageDir, goPathValues)
			if err == nil {
				parsedStructs = append(parsedStructs, parsed...)
			} else {
				log.Fatalln(err)
			}
		}

		//log.Println("PARSED")
		//for _, parsedStruct := range parsedStructs {
		//	log.Printf("%#v\n", parsedStruct)
		//}

		//log.Println("DEFINITIONS")
		pckgDef := parser.NewPackageDefinition(packageDir, filepath.Join(paths[0], packageDir, "lolsql"), parsedStructs)
		for _, strDef := range pckgDef.StructDefinitions {
			log.Printf("%#v\n", strDef)
		}

		generator.Generate(pckgDef)
	} else {

		//ids := []int{10, 102}
		//log.Println(
		//	person.Select(person.Id(), person.Password()).
		//	Where(person.IdIs(10).And(person.IdIs(ids[0], ids[1]))).Or(person.PasswordIsNotNull().And(person.PasswordLike("%10", "001_"))).Render())

		//now := time.Now()
		//log.Println(handsome.Select(handsome.Login(), handsome.DateOfBirth()).Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now).Or(handsome.SalaryIs(10.2, 100.2).Or(handsome.SalaryIs(-100)))).And(handsome.LoginIs("LoginIs")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIs("LoginIs", "LoginIs2")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginIsNot("LoginIsNot", "LoginIsNot2")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("%P1")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("%P1")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginLike("LoginLike", "LoginLike2")).Render())
		//log.Println(handsome.Select().Where(handsome.DateOfBirthIsNot(&now)).Or(handsome.SalaryIs(10.2, 100.2)).And(handsome.LoginNotLike("LoginNotLike", "LoginNotLike2")).Render())

		db, err := sql.Open("mysql", "root:root@/Development?parseTime=true")
		utils.PanicIfNotNil(err)
		defer db.Close()
		////
		//dals, err := djangoadminlog.Select().Where(djangoadminlog.ObjectReprLike("Harry%", "Lol%").And(djangoadminlog.ActionFlagIs(1)).And(djangoadminlog.ActionFlagIs(1))).Fetch(db)
		//log.Println(
		//	Select(ChangeMessage(), ActionTime()).
		//		Where(ObjectReprIs("Harry%", "Lol%").And(ActionFlagIs(1)).And(ActionFlagIs(1))).And(AskPasswordIs(dbmodel.No)).
		//		Render(),
		//)
		//log.Println(len(dals))
		//var db *sql.DB
		// Select().Where(AskPasswordIs(dbmodel.Yes)).Fetch(db)
		//log.Println(len(dals))

		//ps, err := person.Select().Fetch(db)
		//log.Println(err, ps)
		//

		//sqlRend := person.Select(person.Id(), person.Password()).
		//sqlRend := person.Select(person.Name_First()).
		//	Where(person.IdIs(10).And(person.IdIs(10, 1000))).
		//		Or(person.PasswordIsNotNull().And(person.PasswordLike("%10", "001_"))).
		//		Or(person.FirstNotLike("Pavrl")).
		//		Or(person.SecondIsNull()).
		//	Render()
		//log.Println(sqlRend)

		// select Id, FIRST_NAME, SECOND_NAME, SECRET from PERSONS where (Id = ?) or (Id = ?)
		// select Id, FIRST_NAME, SECOND_NAME, SECRET from PERSONS where Id = ? and (Id = ?)
		qb := person.Select().
			Where(person.IdIs(1).Or(person.IdIs(2)))
			//Or(person.PasswordIsNotNull().And(person.PasswordLike("%10", "001_"))).
			//Or(person.FirstNotLike("Pavrl")).
			//Or(person.SecondIsNull())
		log.Println(qb.Render())
		persons, err := qb.Fetch(db)
		if err != nil {
			log.Fatalln(err)
		}
		for _, per := range persons {
			log.Println(per)
		}

		//ps, err := person.Select(person.Id(), person.Password()).
		//	Where(person.IdIs(10).And(
		//		person.IdIs(10, 1000))).Or(
		//			person.PasswordIsNotNull().And(
		//			person.PasswordLike("%10", "001_")),
		//).Fetch(db)

		//r,e := db.Query("select Id, SECRET from PERSONS where SECRET like ?", 10)
		//log.Println(e, r)
		//
		//
		//log.Println(person.Select(person.Id(), person.Password()).Where(person.PasswordLike("%")).Render())
		//ps, err := person.Select(person.Id(), person.Password()).Where(person.PasswordLike("%")).Fetch(db)
		//
		//log.Println(err, ps)
	}
}
