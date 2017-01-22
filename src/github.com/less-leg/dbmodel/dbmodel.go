package dbmodel

type Person struct {
	Id       int     `lolsql:"id[true]"`
	Name
	Password *string `lolsql:"column[SECRET]"`
}

func (_ *Person) TableName() string {
	return "PERSONS"
}

type Name struct {
	FirstName  string  `lolsql:"column[FIRST_NAME]"`
	SecondName *string `lolsql:"column[SECOND_NAME]"`
	MiddleName *string `lolsql:"column[MIDDLE_NAME]"`
}
