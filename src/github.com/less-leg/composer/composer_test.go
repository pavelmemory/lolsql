package composer

import (
	"testing"
	"time"

	"database/sql"

	"github.com/less-leg/dbmodel"
	"github.com/less-leg/generated/order"
	llsql "github.com/less-leg/sql"
	"github.com/less-leg/test_model"
)

func TestName(t *testing.T) {
	var (
		orders []test_model.Order
		err    error
	)
	orders, err = order.UnitOfWork{}.Order(order.Id(), order.Owner(), order.Start()).
		Where(
			llsql.Or(
				order.Id().Equal(1).Or(order.Owner().Name().Like("Pavlo%")),
				order.TaxFree().Equal(dbmodel.Confirmation("yes")),
				order.Owner().Version().NotIn(1, 4, 9)),
			order.Start().Between(time.Now(), time.Now())).
		//OrderBy(sql.Desc(order.Id())).
		//GroupBy(order.Id(), order.Owner().Name()).
		//Having(order.Id().Equal(1)).
		Get()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orders)
}

// для забора данных нужно знать:
// 1. какие данные нужно забирать
// 1.1 как их забирать
// 2. по каким условиям их забирать
// 3. какие дополнительные запросы нужно выполнить для полноты данных
// 4. вложенность дополнительных запросов по данным

/*
допустим что у нас есть таблица Users и структура User
таблица Users имеет составной ключ id, version
в которую входит также множество структур Order хранящиеся в таблице Orders
таблица Orders имеет составной ключ id, version и отношение внешнего ключа на таблицу Users user_id, user_version

первый запрос изымает необходимые сущности из таблицы Users
делается аггрегация внешних ключей для выполнения выборки из таблицы Orders
выборка из таблицы Orders и выполнение маппинга сущностей Order к соответсвующим сущностям User

*/
func Test2(t *testing.T) {
	db, err := sql.Open("", "")
	if err != nil {
		t.Fatal(err)
	}

	db.Query("", args...)

}
