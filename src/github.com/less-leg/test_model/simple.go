package test_model

import (
	"github.com/less-leg/dbmodel"
	"time"
)

type MyTime time.Time

type User struct {
	Name    string
	Father  *User
	time.Time
	TaxFree dbmodel.Confirmation
}

type (
	Operation struct {
		Begin, End time.Time
		*User
	}
)

var SomeVar string

const SomeConst = 0

func SomeFunc() {
	var op Operation
	op.User = new(User)
}
