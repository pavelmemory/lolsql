package test_model

import (
	"github.com/less-leg/dbmodel"
	"time"
)

type MyTime time.Time

type User struct {
	Name   string
	Father *User `json:"omitempty"`
	time.Time
	TaxFree dbmodel.Confirmation
}

type Order struct {
	Id       int
	Version  int
	Customer *string
	Owner    *User
	TaxFree dbmodel.Confirmation
	Start    time.Time
}

var SomeVar string

const SomeConst = 0

func SomeFunc() {
	var op Operation
	op.User = new(User)
}
