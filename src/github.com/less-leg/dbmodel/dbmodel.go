package dbmodel

type Person struct {
	Id       *int      `lolsql:"id[true] column[ID]"`
	Name               `lolsql:"embedded[true]"`
	Address  *Address  `lolsql:"fk[PERSON_ID] fv[Person.Id]"`
	Children []*Person `lolsql:"jtable[PERSON_CHILD] jby[Person.Id] jon[FATHER_ID] jf[PERSON_ID] jv[Person.Id]"`
}

func (_ *Person) TableName() string {
	return "PERSONS"
}

type Name struct {
	FirstName  string  `lolsql:"column[FIRST_NAME]"`
	SecondName *string `lolsql:"column[SECOND_NAME]"`
	MiddleName *string `lolsql:"column[MIDDLE_NAME]"`
}

type Address struct {
	Identifier int     `lolsql:"id[true] column[ID]"`
	Street     *string
	Building   *string
	Block      *string
	Flat       *string
	Person     *Person `lolsql:"fk[PERSON_ID] fv[Person.Id]"`
}

func (a *Address) TableName() string {
	return "CITY_ADDRESS"
}


//select ID, FIRST_NAME, SECOND_NAME, MIDDLE_NAME from PERSONS where FIRST_NAME = 'Bob';
//select ID, Street, Building, Block, Flat from CITY_ADDRESS where PERSON_ID in (select ID from PERSONS where FIRST_NAME = 'Bob')
//
//select t1.ID, t1.FIRST_NAME, t1.SECOND_NAME, t1.MIDDLE_NAME,
//t2.ID, t2.Street, t2.Building, t2.Block, t2.Flat
//from PERSONS t1 join CITY_ADDRESS t2 on t1.ID = t2.PERSON_ID and t1.FIRST_NAME = 'Bob';
//
//rows, err := db.Query()
//
