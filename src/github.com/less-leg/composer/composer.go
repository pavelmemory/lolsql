package composer

import (
	"github.com/less-leg/test_model"
)

type Order struct {
	Id       int
	Version  int
	Customer *string
	Owner    *test_model.User
}

func Start() {
	Order(Id, Owner()).
		Where(Disj(Id.Is(1).Or(Customer.Name.Like("Pavlo%")), Customer.Version.NotIn(1, 4, 9))).
		OrderBy(Id.Desc).
		GroupBy(Id, Customer.Name).
		Having(Id.Is(1)).
		Get()
}

/*
 - for each entity there is a same-name function to specify fields to retrieve
 - each field represent a global read-only variable and has a bunch of methods to use in
   - WHERE
   - GROUP BY
   - HAVING
   - ORDER BY
  such as
   - Is - same as `=`
   - IsNot - same as `<>` or `!=`
   - Like - same as `LIKE '_va%' AND LIKE '123va%'`
   - LikeOr - same as `LIKE '_va%' OR LIKE '11va%'`
   - NotLike - same as `NOT LIKE '_va%' AND NOT LIKE '123va%'`
   - NotLikeOr - same as `NOT LIKE '_va%' OR NOT LIKE '11va%'`
   - Greater - same as `>`
   - Lower - same as `<`
   - IsOrGreater - same as `>=`
   - IsOrLower - same as `<=`
   - Between - same as `BETWEEN 1 AND 4`
   - IsNull - same as `IS NULL`
   - IsNotNull - same as `IS NOT NULL`
   - In - same as `IN (...)`
   - NotIn - same as `NOT IN (...)`
*/

func (b orderBuilder) GroupBy(column ...orderField) orderBuilder {
	return b
}

type OrderDirection string

type orderField struct {
	Desc OrderDirection
}

type userField struct {
	Desc OrderDirection
}

type relatedOrderEntity struct {
	Name    orderField
	Version orderField
}

func (f orderField) Is(int) cond {
	return cond{}
}

func (f orderField) NotIn(int, ...int) cond {
	return cond{}
}

func (f orderField) Like(string) cond {
	return cond{}
}

type orderBuilder struct {
	fields []orderField
}

type userBuilder struct {
	fields []userField
}

type cond struct{}

var (
	Id       orderField
	Customer relatedOrderEntity
)

func Disj(c ...cond) cond {
	return cond{}
}

func Conj(c ...cond) cond {
	return cond{}
}

func (c cond) And(cond) cond {
	return c
}

func (c cond) Or(cond) cond {
	return c
}

var allFields = []orderField{Id, Owner()}

func Owner(fields ...userField) userBuilder {
	return User(fields...)
}

func User(fields ...userField) userBuilder {
	return userBuilder{fields: fields}
}

func Order(fields ...orderField) orderBuilder {
	if len(fields) == 0 {
		return orderBuilder{fields: allFields}
	} else {
		return orderBuilder{fields: fields}
	}
	return orderBuilder{}
}

func (b orderBuilder) Where(...cond) orderBuilder {
	return b
}

func (b orderBuilder) OrderBy(...OrderDirection) orderBuilder {
	return b
}

func (b orderBuilder) Having(...cond) orderBuilder {
	return b
}

func (b orderBuilder) Get() []Order {
	return b
}

func (b orderBuilder) GetPtr() []*Order {
	return b
}
