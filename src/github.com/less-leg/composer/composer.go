package composer

import "os/user"

type Order struct {
	Id       int
	Version  int
	Customer *string
	Owner    *user.User
}

func Start() {
	Order(Id, Owner).
		Where(Disj(Conj(Id.Equals(1), Customer.Name.Like("Pavlo%")), Customer.Version.NotIn(1, 4, 9))).
		OrderBy(Id.Desc).
		Having().
		Get()
}

type OrderDirection string

type fieldColumn struct {
	Desc OrderDirection
}

type relatedOrderEntity struct {
	Name    fieldColumn
	Version fieldColumn
}

func (f fieldColumn) Equals(int) cond {
	return cond{}
}

func (f fieldColumn) NotIn(int, ...int) cond {
	return cond{}
}

func (f fieldColumn) Like(string) cond {
	return cond{}
}

type builder struct{}
type cond struct{}

var (
	Id, Owner fieldColumn
	Customer  relatedOrderEntity
)

func Disj(c ...cond) cond {
	return cond{}
}

func Conj(c ...cond) cond {
	return cond{}
}

func Order(fields ...fieldColumn) builder {
	return builder{}
}

func (b builder) Where(...cond) builder {
	return b
}

func (b builder) OrderBy(...OrderDirection) builder {
	return b
}

func (b builder) Having() builder {
	return b
}

func (b builder) Get() builder {
	return b
}
