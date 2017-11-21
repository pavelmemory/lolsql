package test_model

import "time"

type (
	Operation struct {
		Begin, End time.Time
		*User
		Actions []*byte
	}
)
