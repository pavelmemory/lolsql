package test_model

import (
	"github.com/less-leg/dbmodel"
	"time"
)

type MyTime time.Time

type UserKey struct {
	Identifier int
	Version  int
}

type User struct {
	UserKey
	Name   string
	Father *User `json:"omitempty"`
	time.Time
	TaxFree dbmodel.Confirmation
}

type Order struct {
	Id       int `lol:"pk,column[PK]"`
	Version  int `lol:"PK"`
	Customer *string
	// add support to provide mapping of fields to embedded field if they have same corresponding field names
	Owner    *User `lol:"FK[Id,Version -> UserKey.Identifier, UserKey.Version]"`
	TaxFree  dbmodel.Confirmation
	Start    time.Time
}

var SomeVar string

const SomeConst = 0

func SomeFunc() {
	var op Operation
	op.User = new(User)
}
