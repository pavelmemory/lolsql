package generator

import (
	"database/sql"
	"github.com/less-leg/utils"
)

func PropagateInt(to *int) {
	db, err := sql.Open("driver", "datasource")
	utils.PanicIf(err)
	rows, err := db.Query("select f1, f2 from tab")
	if err != nil {

	}
}