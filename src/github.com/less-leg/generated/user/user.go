package user

type OrderDirection string

type UserField interface {
}

type userBuilder struct {
	fields []UserField
}

func Order() userBuilder {
	return userBuilder{}
}

func User(fields ...UserField) userBuilder {
	return userBuilder{fields: fields}
}

type UserName interface {
	UserField
}

func Name() UserName {
	return nil
}
