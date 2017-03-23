package dbmodel

type Person struct {
	Id int `lolsql:"id[true]"`
	Name
	Password *string `lolsql:"column[SECRET]"`
}

func (_ *Person) TableName() string {
	return "PERSONS"
}

type Name struct {
	First  string  `lolsql:"column[FIRST_NAME]"`
	Second *string `lolsql:"column[SECOND_NAME]"`
	//MiddleName *string `lolsql:"column[MIDDLE_NAME]"`
}
