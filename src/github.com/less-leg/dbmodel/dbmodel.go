package dbmodel

import (
	"fmt"
)

type Person struct {
	Id       int `lolsql:"id[true]"`
	Name     *Name
	Password *string `lolsql:"column[SECRET]"`
}

func (p Person) String() string {
	pas := ""
	if p.Password != nil {
		pas = *p.Password
	}
	return fmt.Sprintf(`"Person":{"Id":%d,%s,"Password":"%s"}`, p.Id, p.Name, pas)
}

func (_ *Person) TableName() string {
	return "PERSONS"
}

type Name struct {
	First  string  `lolsql:"column[FIRST_NAME]"`
	Second *string `lolsql:"column[SECOND_NAME]"`
	//MiddleName *string `lolsql:"column[MIDDLE_NAME]"`
}

func (n Name) String() string {
	sec := ""
	if n.Second != nil {
		sec = *n.Second
	}
	return fmt.Sprintf(`"Name":{"First":"%s","Second":"%s"}`, n.First, sec)
}
