package dbmodel

import (
	"fmt"
	"database/sql/driver"
)

type Confirmation struct {
	string
}

var Yes = Confirmation{"Yes"}
var No = Confirmation{"No"}

func (this *Confirmation) Scan(src interface{}) (err error) {

	define := func(v string) (Confirmation, error) {
		switch v {
		case Yes.string:
			return Yes, nil
		case No.string:
			return No, nil
		default:
			return Confirmation{""}, fmt.Errorf("Value: %s is not one of supported values: %s|%s", v, Yes, No)
		}
	}

	switch v := src.(type) {
	case []byte:
		*this, err = define(string(v))
	case string:
		*this, err = define(v)
	case nil:
		*this, err = define("")
	default:
		err = fmt.Errorf("Cannot convert %T to Confirmation", src)
	}
	return
}

func (this Confirmation) Value() (driver.Value, error) {
	switch this.string {
	case Yes.string, No.string: return this.string, nil
	default: return nil, fmt.Errorf("Cannot convert Confirmation object to valid dabatase representation: %s", this.string)
	}
}

