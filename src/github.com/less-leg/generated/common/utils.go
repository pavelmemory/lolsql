package common

import "time"

func MultiTimeTime(times ...time.Time) (vals []interface{}) {
	for _, v := range times {
		vals = append(vals, v)
	}
	return
}
