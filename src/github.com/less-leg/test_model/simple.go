package test_model

import "time"

type MyTime time.Time

type User struct {
	Name   string
	Years  int
	Father *User
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
