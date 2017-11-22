package composer

import (
	"github.com/less-leg/generated/order"
	"github.com/less-leg/sql"
	"github.com/less-leg/test_model"
	"testing"
	"github.com/less-leg/dbmodel"
)

func TestName(t *testing.T) {
	var (
		orders []test_model.Order
		err    error
	)
	orders, err = order.Order(order.Id(), order.Owner()).
		Where(
			sql.Or(
				order.Id().Equal(1).Or(order.Owner().Name().Like("Pavlo%")),
				order.TaxFree().Equal(dbmodel.Confirmation("yes")),
				order.Owner().Version().NotIn(1, 4, 9))).
		OrderBy(sql.Desc(order.Id())).
		GroupBy(order.Id(), order.Owner().Name()).
		Having(order.Id().Equal(1)).
		Get()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orders)
}
