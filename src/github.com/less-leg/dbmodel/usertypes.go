package dbmodel

import (
	"fmt"
	"database/sql/driver"
)

type Confirmation string

const Yes = Confirmation("Yes")
const No = Confirmation("No")

func (this *Confirmation) Scan(src interface{}) (err error) {

	define := func(v Confirmation) (Confirmation, error) {
		switch v {
		case Yes, No:
			return v, nil
		default:
			return Confirmation(""), fmt.Errorf("Value: %s is not one of supported values: %s|%s", v, Yes, No)
		}
	}

	switch v := src.(type) {
	case []byte:
		*this, err = define(Confirmation(string(v)))
	case string:
		*this, err = define(Confirmation(v))
	default:
		err = fmt.Errorf("Cannot convert %T to Confirmation", src)
	}
	return
}

func (this Confirmation) Value() (driver.Value, error) {
	switch this {
	case Yes, No: return string(this), nil
	default: return nil, fmt.Errorf("Cannot convert Confirmation object to valid dabatase representation: %s", this)
	}
}
