package sql

type SortOrder interface{}

func Desc(Field) SortOrder {
	return nil
}

func Asc(Field) SortOrder {
	return nil
}
